package converter

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func ScheduleToResponse(s *domain.Schedule) *models.ScheduleResponse {
	return &models.ScheduleResponse{
		Schedule: ScheduleToObject(s),
	}
}

func ScheduleToObject(s *domain.Schedule) *models.ScheduleObject {
	return &models.ScheduleObject{
		ID:         s.ID(),
		RoomID:     s.RoomID(),
		DaysOfWeek: WeekdaysToIntArray(s.DaysOfWeek()),
		StartTime:  TimeHourMinuteToString(s.StartTime()),
		EndTime:    TimeHourMinuteToString(s.EndTime()),
	}
}
