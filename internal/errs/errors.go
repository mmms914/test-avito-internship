package errs

var (
	ErrDayIsInvalid      = NewCustomError(InvalidRequestErrorCode, "day should be between 1 and 7")
	ErrRoomNotExists     = NewCustomError(RoomNotFoundErrorCode, "room does not exist")
	ErrUnauthorized      = NewCustomError(UnauthorizedErrorCode, "unauthorized")
	ErrNotFound          = NewCustomError(NotFoundErrorCode, "not found")
	ErrUserNotFound      = NewCustomError(NotFoundErrorCode, "user not found")
	ErrSlotNotFound      = NewCustomError(SlotNotFoundErrorCode, "slot not found")
	ErrSlotAlreadyBooked = NewCustomError(SlotAlreadyBookedErrorCode, "slot already booked")
	ErrBookingNotFound   = NewCustomError(BookingNotFoundErrorCode, "booking not found")
	ErrForbidden         = NewCustomError(ForbiddenErrorCode, "forbidden")
	ErrScheduleExists    = NewCustomError(ScheduleExistsErrorCode, "schedule already exists")
	ErrInternalError     = NewCustomError(InternalErrorCode, "internal error")
)
