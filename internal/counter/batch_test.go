package counter

import (
	"strings"
	"testing"
)

// TestBatchCounter_Count verifies that the batch counter correctly aggregates name counts from
// multiple input scenarios: standard input, inputs with extra whitespace and empty lines, and
// empty input. These cases exercise trimming, line-splitting and batch processing to cover
// prior technical issues related to batch splitting and line parsing.
func TestBatchCounter_Count(t *testing.T) {
	tests := []struct {
		name      string
		batchSize int
		input     string
		expected  map[string]int
		wantErr   bool
	}{
		{
			name:      "Successful count with standard input",
			batchSize: 2,
			input:     "Алёна\nМиша\nАлёна\nДима\n",
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
				"Дима":  1,
			},
			wantErr: false,
		},
		{
			name:      "Batch size is less than 0",
			batchSize: -10,
			input:     "Алёна\nМиша\nАлёна\nДима\n",
			expected: map[string]int{
				"Алёна": 2,
				"Миша":  1,
				"Дима":  1,
			},
			wantErr: false,
		},
		{
			name:      "Input with spaces, tabs and empty lines",
			batchSize: 10,
			input:     "  Алёна  \n\t\t\nМиша\n\t\n",
			expected: map[string]int{
				"Алёна": 1,
				"Миша":  1,
			},
			wantErr: false,
		},
		{
			name:      "Empty input",
			batchSize: 10,
			input:     "",
			expected:  map[string]int{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewBatchCounter(tt.batchSize)
			r := strings.NewReader(tt.input)

			got, err := c.Count(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Count() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("Count() map length = %v, expected length %v", len(got), len(tt.expected))
			}

			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("Count()[%q] = %v, expected %v", k, got[k], v)
				}
			}
		})
	}
}
