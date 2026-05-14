// Package main is the entry point for the name-counter executable.
package main

import (
	"log/slog"
	"os"

	"github.com/ilindan-dev/name-counter/internal/cli"
)

// main executes the root command and exits with a non-zero status code on failure.
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	err := cli.Execute(logger)
	if err != nil {
		logger.Error("application execution failed", "error", err)
		os.Exit(1)
	}
}
