// Package sorter provides implementations of the core.Sorter interface.
package sorter

import (
	"cmp"
	"slices"

	"github.com/ilindan-dev/name-counter/internal/core"
)

// Compile check that memorySorter implements core.Sorter.
var _ core.Sorter = (*memorySorter)(nil)

// memorySorter implements core.Sorter using in-memory slice sorting.
type memorySorter struct {
	mode core.OutputMode
}

// NewMemorySorter creates a new Sorter that performs sorting in RAM.
func NewMemorySorter(mode core.OutputMode) core.Sorter {
	return &memorySorter{
		mode: mode,
	}
}

// Sort converts the frequency map into a slice and sorts it based on the configured mode.
func (s *memorySorter) Sort(counts map[string]int) ([]core.NameFrequency, error) {
	result := make([]core.NameFrequency, 0, len(counts))

	for name, count := range counts {
		result = append(result, core.NameFrequency{
			Name:  name,
			Count: count,
		})
	}

	if s.mode == core.ModeNone {
		return result, nil
	}

	slices.SortFunc(result, func(a, b core.NameFrequency) int {
		if a.Count == b.Count {
			return cmp.Compare(a.Name, b.Name)
		}

		if s.mode == core.ModeAsc {
			return cmp.Compare(a.Count, b.Count)
		}

		return cmp.Compare(b.Count, a.Count)
	})

	return result, nil
}
