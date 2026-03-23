package models

const (
	InvalidRequestErrorCode    string = "INVALID_REQUEST"
	UnauthorizedErrorCode      string = "UNAUTHORIZED"
	NotFoundErrorCode          string = "NOT_FOUND"
	RoomNotFoundErrorCode      string = "ROOM_NOT_FOUND"
	SlotNotFoundErrorCode      string = "SLOT_NOT_FOUND"
	SlotAlreadyBookedErrorCode string = "SLOT_ALREADY_BOOKED"
	BookingNotFoundErrorCode   string = "BOOKING_NOT_FOUND"
	ForbiddenErrorCode         string = "FORBIDDEN"
	ScheduleExistsErrorCode    string = "SCHEDULE_EXISTS"
	InternalErrorCode          string = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Error *Error `json:"error"`
}
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
