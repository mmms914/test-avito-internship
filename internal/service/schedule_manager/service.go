package schedule_manager

import (
	"fmt"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/google/uuid"
)

type RoomRepository interface {
	Exists(roomID uuid.UUID) (bool, error)
}

type ScheduleRepository interface {
	Create(schedule *domain.Schedule) (*domain.Schedule, error)
}

type ScheduleManager struct {
	roomRepo     RoomRepository
	scheduleRepo ScheduleRepository
}

type Config struct {
	roomRepo     RoomRepository
	scheduleRepo ScheduleRepository
}

func NewScheduleManager(c *Config) *ScheduleManager {
	return &ScheduleManager{
		roomRepo:     c.roomRepo,
		scheduleRepo: c.scheduleRepo,
	}
}
func (sm *ScheduleManager) Create(schedule *domain.Schedule) (*domain.Schedule, error) {
	roomExists, err := sm.roomRepo.Exists(schedule.RoomID())
	if err != nil {
		return nil, fmt.Errorf("checking if room exists: %v", err)
	}

	if !roomExists {
		return nil, fmt.Errorf("room %s does not exist", schedule.RoomID())
	}

	createdSchedule, err := sm.scheduleRepo.Create(schedule)
	if err != nil {
		return nil, fmt.Errorf("creating schedule: %w", err)
	}

	return createdSchedule, nil
}
