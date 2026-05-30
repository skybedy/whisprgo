package main

import (
	"fmt"
	"os"

	"whisprgo/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "whisprgo: %v\n", err)
		os.Exit(1)
	}

	if err := application.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "whisprgo: %v\n", err)
		os.Exit(1)
	}
}
