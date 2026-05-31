package main

import (
	"fmt"
	"os"

	"whispergo/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "whispergo: %v\n", err)
		os.Exit(1)
	}

	if err := application.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "whispergo: %v\n", err)
		os.Exit(1)
	}
}
