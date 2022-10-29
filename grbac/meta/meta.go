package meta

import (
	"errors"
)

// define a set of errors
var (
	ErrFieldIncomplete = errors.New("incomplete fields")
	ErrEmptyStructure  = errors.New("empty structure")
)
