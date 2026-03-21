package errs

import "errors"

var (
	ErrDayIsInvalid error = errors.New("day should be between 1 and 7")
)
