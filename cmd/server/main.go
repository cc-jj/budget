package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cc-jj/budgetapp/server"
)

func main() {
	ctx := context.Background()
	if err := server.Run(ctx, os.Getenv, os.Stdin, os.Stdout, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}