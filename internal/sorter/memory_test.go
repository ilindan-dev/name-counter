package sorter_test

import (
	"reflect"
	"testing"

	"github.com/ilindan-dev/name-counter/internal/core"
	"github.com/ilindan-dev/name-counter/internal/sorter"
)

// TestMemorySorter_Sort verifies MemorySorter ordering across modes:
// ModeDesc and ModeAsc: items are ordered by count with deterministic tie-breaking for equal counts.
// ModeNone: returns results without ordering while preserving the number of items.
func TestMemorySorter_Sort(t *testing.T) {
	inputCounts := map[string]int{
		"Миша":  1,
		"Алёна": 3,
		"Дима":  2,
		"Антон": 2,
	}

	tests := []struct {
		name     string
		mode     core.OutputMode
		expected []core.NameFrequency
		wantErr  bool
	}{
		{
			name: "Sort desc mode",
			mode: core.ModeDesc,
			expected: []core.NameFrequency{
				{Name: "Алёна", Count: 3},
				{Name: "Антон", Count: 2},
				{Name: "Дима", Count: 2},
				{Name: "Миша", Count: 1},
			},
			wantErr: false,
		},
		{
			name: "Sort asc mode",
			mode: core.ModeAsc,
			expected: []core.NameFrequency{
				{Name: "Миша", Count: 1},
				{Name: "Антон", Count: 2},
				{Name: "Дима", Count: 2},
				{Name: "Алёна", Count: 3},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := sorter.NewMemorySorter(tt.mode)
			got, err := s.Sort(inputCounts)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Sort() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("\nExpected: %v\nGot: %v", tt.expected, got)
			}
		})
	}

	t.Run("Sort in none mode)", func(t *testing.T) {
		s := sorter.NewMemorySorter(core.ModeNone)
		got, err := s.Sort(inputCounts)
		if err != nil {
			t.Fatalf("Sort() error = %v", err)
		}

		if len(got) != len(inputCounts) {
			t.Errorf("Sort() length = %v, expected %v", len(got), len(inputCounts))
		}
	})
}
