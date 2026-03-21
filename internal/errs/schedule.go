package errs

const (
	ValidationErrorCode string = "VALIDATION_ERROR"
)

var (
	ErrDayIsInvalid = NewCustomError(ValidationErrorCode, "day should be between 1 and 7")
)
