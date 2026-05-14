// Package counter provides implementations of the core.Counter interface.
package counter

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ilindan-dev/name-counter/internal/core"
)

// defaultBatchSize contains the default batch size value
// to override the input value if it is invalid
const defaultBatchSize = 1000

// Compile check that batchCounter implements core.Counter.
var _ core.Counter = (*batchCounter)(nil)

// batchCounter implements core.Counter using a batch-reading approach.
type batchCounter struct {
	batchSize int
}

// NewBatchCounter creates a new instance of a batch-based Counter.
// batchSize determines how many lines are read into memory before processing.
func NewBatchCounter(batchSize int) core.Counter {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	return &batchCounter{
		batchSize: batchSize,
	}
}

// Count reads from the io.Reader line by line, grouping them into batches
// to aggregate the frequencies into a map.
func (c *batchCounter) Count(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(r)

	batch := make([]string, 0, c.batchSize)

	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" {
			continue
		}

		batch = append(batch, name)

		if len(batch) >= c.batchSize {
			c.processBatch(batch, counts)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		c.processBatch(batch, counts)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", core.ErrReadInput, err)
	}

	return counts, nil
}

// processBatch iterates over a batch of names and updates the frequency map.
func (c *batchCounter) processBatch(batch []string, counts map[string]int) {
	for _, name := range batch {
		counts[name]++
	}
}
