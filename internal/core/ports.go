package core

import "io"

// Counter defines the contract for reading data and aggregating name frequencies.
type Counter interface {
	// Count reads names from the provided reader and returns an aggregated frequency map.
	// The implementation details (e.g., batching, buffering) are hidden.
	Count(r io.Reader) (map[string]int, error)
}

// Sorter defines the contract for sorting the aggregated name frequencies.
// This allows future implementations, such as external sorting via temp files,
// to be seamlessly integrated.
type Sorter interface {
	// Sort takes a map of name frequencies and returns a slice formatted
	// according to the implementation's rules (e.g., by OutputMode).
	Sort(counts map[string]int) ([]NameFrequency, error)
}
