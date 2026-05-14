package core

// OutputMode represents the desired sorting behavior for the output.
type OutputMode string

const (
	// ModeDesc sorts names in descending order by frequency.
	ModeDesc OutputMode = "desc"
	// ModeAsc sorts names in ascending order by frequency.
	ModeAsc OutputMode = "asc"
	// ModeNone leaves the names in an arbitrary order.
	ModeNone OutputMode = "none"
)

// NameFrequency holds a name and its occurrence count.
type NameFrequency struct {
	Name  string
	Count int
}
