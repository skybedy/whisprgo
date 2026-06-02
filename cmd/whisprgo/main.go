package main

import (
	"fmt"
	"os"

	"whisprgo/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "whisprgo: %s\n", app.SanitizeText(err.Error()))
		os.Exit(1)
	}

	if err := application.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "whisprgo: %s\n", app.SanitizeText(err.Error()))
		os.Exit(1)
	}
}
