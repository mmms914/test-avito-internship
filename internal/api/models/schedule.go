package models

import (
	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	DaysOfWeek *[]int  `json:"daysOfWeek" validate:"required"`
	StartTime  *string `json:"startTime" validate:"required"`
	EndTime    *string `json:"endTime" validate:"required"`
}

type ScheduleResponse struct {
	Schedule *ScheduleObject `json:"schedule"`
}

type ScheduleObject struct {
	ID         uuid.UUID `json:"id"`
	RoomID     uuid.UUID `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  string    `json:"startTime"`
	EndTime    string    `json:"endTime"`
}
