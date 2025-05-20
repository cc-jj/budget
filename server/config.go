package server

import (
	"flag"
	"io"
	"strings"
)


type Config struct {
	Host string
	Port string
}

func newConfig(stdin io.Reader, getenv func(string) string) *Config {
	defaultHost := strings.TrimSpace(getenv("SERVER_HOST"))
	if len(defaultHost) == 0 {
		defaultHost = "localhost"
	}
	defaultPort := strings.TrimSpace(getenv("SERVER_PORT"))
	if len(defaultPort) == 0 {
		defaultPort = "4242"
	}
	host := flag.String("host", defaultHost, "The server host. May also set with env var 'SERVER_HOST'.")
	port := flag.String("port", defaultPort, "The server port. May also set with evn var 'SERVER_PORT'.")

	flag.Parse()
	
	return &Config{Host: *host, Port: *port}
}