package errs

var (
	ErrRoomNotExists     = NewCustomError(RoomNotFoundErrorCode, "room does not exist")
	ErrUnauthorized      = NewCustomError(UnauthorizedErrorCode, "unauthorized")
	ErrSlotTimeInPast    = NewCustomError(InternalErrorCode, "slot time is in past")
	ErrUserNotFound      = NewCustomError(NotFoundErrorCode, "user not found")
	ErrUserAlreadyExists = NewCustomError(UserAlreadyExistsErrorCode, "user with that email already exists")
	ErrSlotNotFound      = NewCustomError(SlotNotFoundErrorCode, "slot not found")
	ErrSlotAlreadyBooked = NewCustomError(SlotAlreadyBookedErrorCode, "slot already booked")
	ErrBookingNotFound   = NewCustomError(BookingNotFoundErrorCode, "booking not found")
	ErrForbidden         = NewCustomError(ForbiddenErrorCode, "forbidden")
	ErrScheduleExists    = NewCustomError(ScheduleExistsErrorCode, "schedule already exists")
)
