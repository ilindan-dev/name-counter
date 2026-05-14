// Package core defines the domain models, interfaces, and common errors for the application.
package core

import "errors"

// ErrReadInput indicates a failure occurred while reading the data source.
var ErrReadInput = errors.New("failed to read input data")
