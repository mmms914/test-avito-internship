package errs

const (
	UnauthorizedErrorCode      string = "UNAUTHORIZED"
	NotFoundErrorCode          string = "NOT_FOUND"
	UserAlreadyExistsErrorCode string = "USER_ALREADY_EXISTS"
	RoomNotFoundErrorCode      string = "ROOM_NOT_FOUND"
	SlotNotFoundErrorCode      string = "SLOT_NOT_FOUND"
	SlotAlreadyBookedErrorCode string = "SLOT_ALREADY_BOOKED"
	BookingNotFoundErrorCode   string = "BOOKING_NOT_FOUND"
	ForbiddenErrorCode         string = "FORBIDDEN"
	ScheduleExistsErrorCode    string = "SCHEDULE_EXISTS"
	InternalErrorCode          string = "INTERNAL_ERROR"
)

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
