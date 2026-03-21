package errs

type CustomError struct {
	code    string
	message string
}

func NewCustomError(code string, message string) *CustomError {
	return &CustomError{code: code, message: message}
}

func (e *CustomError) Code() string {
	return e.code
}

func (e *CustomError) Error() string {
	return e.message
}
