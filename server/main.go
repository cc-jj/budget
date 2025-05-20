package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)


func newServer(logger *log.Logger) (http.Handler, error) {
	if logger == nil {
		logger = log.Default()
	}

	mux := http.NewServeMux();

	// TODO (cc-jj) Cache headers
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	if err := addRoutes(mux, logger); err != nil {
		return nil, fmt.Errorf("error adding routes: %w", err)
	}
	
	var handler http.Handler = mux
	handler = LoggingMiddleware(logger, handler)
	return handler, nil
}


func Run(
	signalContext    context.Context,
	getenv func(string) string,
	stdin  io.Reader,
	stdout, stderr io.Writer,
) error {
	signalContext, cancel := signal.NotifyContext(signalContext, os.Interrupt)
	defer cancel()

	logger := log.New(stdout, "", log.LstdFlags)
	config := newConfig(stdin, getenv)

	srv, err := newServer(logger)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	
	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-signalContext.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10 * time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
