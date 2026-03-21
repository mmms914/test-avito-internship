package errs

import "errors"

var (
	ErrDayIsInvalid = errors.New("day should be between 1 and 7")
)
