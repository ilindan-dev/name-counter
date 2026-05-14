// Package cli provides the command-line interface for the name-counter application.
// It handles argument parsing, flag configuration, and execution of the core domain logic.
package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/ilindan-dev/name-counter/internal/core"
	"github.com/ilindan-dev/name-counter/internal/counter"
	"github.com/ilindan-dev/name-counter/internal/sorter"
)

const (
	defaultBatchSize = 1000
	defaultMode      = string(core.ModeDesc)
)

// Execute initializes the root command, configures its flags, and executes it.
// It returns an error if the command execution or argument parsing fails.
func Execute(logger *slog.Logger) error {
	var (
		mode      string
		batchSize int
	)

	rootCmd := &cobra.Command{
		Use:   "name-counter [file]",
		Short: "A fast utility to count name frequencies in a file",
		Long: "name-counter reads a text file containing names (one per line) and counts the frequency of each " +
			"unique name. It supports batch reading for memory efficiency and multiple sorting modes for the output.",
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return runProcess(args[0], mode, batchSize, logger)
		},
	}

	rootCmd.Flags().StringVarP(&mode, "mode", "m", defaultMode,
		"Sorting mode for output: 'desc', 'asc', or 'none'")

	rootCmd.Flags().IntVarP(&batchSize, "batch-size", "b", defaultBatchSize,
		"Number of lines to read in a single batch")

	return rootCmd.Execute()
}

// runProcess orchestrates the core logic by validating inputs, initializing dependencies,
// and executing the data flow from reading the file to printing the results.
func runProcess(filePath, mode string, batchSize int, logger *slog.Logger) error {
	if mode != string(core.ModeDesc) && mode != string(core.ModeAsc) && mode != string(core.ModeNone) {
		return fmt.Errorf("invalid mode: %s. Valid options are: desc, asc, none", mode)
	}

	if batchSize <= 0 {
		return fmt.Errorf("invalid batch-size: %d. Must be greater than 0", batchSize)
	}

	logger.Debug("starting process", "file", filePath, "mode", mode, "batch_size", batchSize)

	//nolint:gosec // this is a CLI tool, reading files from arbitrary user input is intended
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			logger.Warn("failed to close file", "error", closeErr)
		}
	}()

	c := counter.NewBatchCounter(batchSize)
	s := sorter.NewMemorySorter(core.OutputMode(mode))

	counts, err := c.Count(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	logger.Debug("counting finished", "unique_names", len(counts))

	sortedResults, err := s.Sort(counts)
	if err != nil {
		return fmt.Errorf("failed to sort results: %w", err)
	}

	for _, res := range sortedResults {
		fmt.Printf("%s:%d\n", res.Name, res.Count)
	}

	return nil
}
