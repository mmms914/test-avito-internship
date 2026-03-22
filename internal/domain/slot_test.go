package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
)

func TestGenerateSlots(t *testing.T) {
	// Фиксированные времена для тестов
	fixedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)   // Понедельник
	tuesdayDate := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // Вторник
	sundayDate := time.Date(2024, 1, 21, 0, 0, 0, 0, time.UTC)  // Воскресенье

	roomID := uuid.New()

	tests := []struct {
		name     string
		schedule *domain.Schedule
		date     time.Time
		expected func(t *testing.T, slots []*domain.Slot)
	}{
		{
			name: "success - generate slots for working day (9:00-18:00)",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    roomID,
					StartTime: 9 * time.Hour,
					EndTime:   18 * time.Hour,
					DaysOfWeek: []time.Weekday{
						time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
					},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// 9:00 до 18:00 = 9 часов = 18 слотов по 30 минут
				assert.Len(t, slots, 18)

				// Проверяем первый слот
				if len(slots) > 0 {
					expectedStart := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
					assert.Equal(t, expectedStart, slots[0].StartTime())
					assert.Equal(t, expectedStart.Add(domain.SlotDuration), slots[0].EndTime())
				}

				// Проверяем последний слот
				if len(slots) > 0 {
					expectedStart := time.Date(2024, 1, 15, 17, 30, 0, 0, time.UTC)
					assert.Equal(t, expectedStart, slots[17].StartTime())
					assert.Equal(t, expectedStart.Add(domain.SlotDuration), slots[17].EndTime())
				}

				// Проверяем, что все слоты имеют корректные ID и roomID
				for _, slot := range slots {
					assert.NotEqual(t, uuid.Nil, slot.ID())
					assert.Equal(t, roomID, slot.RoomID())
				}

				// Проверяем, что интервалы между слотами равны SlotDuration
				for i := 1; i < len(slots); i++ {
					expectedInterval := domain.SlotDuration
					actualInterval := slots[i].StartTime().Sub(slots[i-1].StartTime())
					assert.Equal(t, expectedInterval, actualInterval,
						"Slot %d has wrong interval", i)
				}
			},
		},
		{
			name: "success - generate slots for half day (10:00-14:00)",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  10 * time.Hour,
					EndTime:    14 * time.Hour,
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// 10:00 до 14:00 = 4 часа = 8 слотов
				assert.Len(t, slots, 8)

				// Проверяем первый слот
				if len(slots) > 0 {
					expectedStart := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
					assert.Equal(t, expectedStart, slots[0].StartTime())
				}

				// Проверяем последний слот
				if len(slots) > 0 {
					expectedStart := time.Date(2024, 1, 15, 13, 30, 0, 0, time.UTC)
					assert.Equal(t, expectedStart, slots[7].StartTime())
				}
			},
		},
		{
			name: "success - generate slots with non-zero minutes start (9:30-17:30)",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  9*time.Hour + 30*time.Minute,
					EndTime:    17*time.Hour + 30*time.Minute,
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// 9:30 до 17:30 = 8 часов = 16 слотов
				assert.Len(t, slots, 16)

				// Проверяем первый слот
				if len(slots) > 0 {
					expectedStart := time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)
					assert.Equal(t, expectedStart, slots[0].StartTime())
					assert.Equal(t, expectedStart.Add(domain.SlotDuration), slots[0].EndTime())
				}
			},
		},
		{
			name: "empty - day not in schedule",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    roomID,
					StartTime: 9 * time.Hour,
					EndTime:   18 * time.Hour,
					DaysOfWeek: []time.Weekday{
						time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday,
					},
				}),
			),
			date: sundayDate, // Воскресенье
			expected: func(t *testing.T, slots []*domain.Slot) {
				assert.Empty(t, slots)
			},
		},
		{
			name: "empty - schedule has no days",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  9 * time.Hour,
					EndTime:    18 * time.Hour,
					DaysOfWeek: []time.Weekday{},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				assert.Empty(t, slots)
			},
		},
		{
			name: "success - exact slot boundary (9:00-12:00)",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  9 * time.Hour,
					EndTime:    12 * time.Hour,
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// 9:00 до 12:00 = 3 часа = 6 слотов
				assert.Len(t, slots, 6)

				// Проверяем точные времена
				expectedTimes := []struct {
					start time.Time
					end   time.Time
				}{
					{time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC), time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)},
					{time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC), time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)},
					{time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)},
					{time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)},
					{time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC), time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC)},
					{time.Date(2024, 1, 15, 11, 30, 0, 0, time.UTC), time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
				}

				for i, expected := range expectedTimes {
					assert.Equal(t, expected.start, slots[i].StartTime(), "Slot %d start time mismatch", i)
					assert.Equal(t, expected.end, slots[i].EndTime(), "Slot %d end time mismatch", i)
				}
			},
		},
		{
			name: "success - slot duration exactly matches (9:00-9:30)",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  9 * time.Hour,
					EndTime:    9*time.Hour + domain.SlotDuration,
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// Всего 1 слот
				assert.Len(t, slots, 1)

				if len(slots) > 0 {
					expectedStart := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
					expectedEnd := time.Date(2024, 1, 15, 9, 30, 0, 0, time.UTC)
					assert.Equal(t, expectedStart, slots[0].StartTime())
					assert.Equal(t, expectedEnd, slots[0].EndTime())
				}
			},
		},
		{
			name: "success - multiple days in schedule",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  9 * time.Hour,
					EndTime:    12 * time.Hour,
					DaysOfWeek: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday},
				}),
			),
			date: tuesdayDate, // Вторник
			expected: func(t *testing.T, slots []*domain.Slot) {
				assert.Len(t, slots, 6) // 3 часа = 6 слотов

				// Проверяем, что слоты созданы для правильного дня
				for _, slot := range slots {
					assert.Equal(t, 2024, slot.StartTime().Year())
					assert.Equal(t, time.January, slot.StartTime().Month())
					assert.Equal(t, 16, slot.StartTime().Day()) // 16 января - вторник
				}
			},
		},
		{
			name: "success - generate slots in UTC",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  9 * time.Hour,
					EndTime:    17 * time.Hour,
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// Проверяем, что все слоты в UTC
				for _, slot := range slots {
					assert.Equal(t, time.UTC, slot.StartTime().Location())
					assert.Equal(t, time.UTC, slot.EndTime().Location())
				}

				// 9:00 до 17:00 = 8 часов = 16 слотов
				assert.Len(t, slots, 16)
			},
		},
		{
			name: "empty - end time before start time",
			schedule: domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  18 * time.Hour,
					EndTime:    9 * time.Hour, // Некорректное расписание
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			),
			date: fixedDate,
			expected: func(t *testing.T, slots []*domain.Slot) {
				// При некорректных данных слотов не будет
				assert.Empty(t, slots)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slots := domain.GenerateSlots(tt.schedule, tt.date)
			tt.expected(t, slots)
		})
	}
}

func TestSlot_IsInPast(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name      string
		startTime time.Time
		expected  bool
	}{
		{
			name:      "slot in past",
			startTime: now.Add(-1 * time.Hour),
			expected:  true,
		},
		{
			name:      "slot in future",
			startTime: now.Add(1 * time.Hour),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot := domain.NewSlot(
				domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    uuid.New(),
					StartTime: tt.startTime,
					EndTime:   tt.startTime.Add(domain.SlotDuration),
				}),
			)

			assert.Equal(t, tt.expected, slot.IsInPast())
		})
	}
}

func TestGenerateSlots_NonDeterministic(t *testing.T) {
	roomID := uuid.New()
	schedule := domain.NewSchedule(
		domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
			ID:         uuid.New(),
			RoomID:     roomID,
			StartTime:  9 * time.Hour,
			EndTime:    17 * time.Hour,
			DaysOfWeek: []time.Weekday{time.Monday},
		}),
	)

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	slots1 := domain.GenerateSlots(schedule, date)
	slots2 := domain.GenerateSlots(schedule, date)

	assert.Len(t, slots1, len(slots2))

	// Времена должны совпадать
	for i := range slots1 {
		assert.Equal(t, slots1[i].StartTime(), slots2[i].StartTime())
		assert.Equal(t, slots1[i].EndTime(), slots2[i].EndTime())
		assert.Equal(t, slots1[i].RoomID(), slots2[i].RoomID())

		// ID должны отличаться (так как генерируются новые UUID)
		assert.NotEqual(t, slots1[i].ID(), slots2[i].ID())
	}
}

func TestGenerateSlots_EdgeCases(t *testing.T) {
	roomID := uuid.New()

	tests := []struct {
		name      string
		startTime time.Duration
		endTime   time.Duration
		expected  int
	}{
		{
			name:      "exactly 30 minutes",
			startTime: 9 * time.Hour,
			endTime:   9*time.Hour + domain.SlotDuration,
			expected:  1,
		},
		{
			name:      "exactly 1 hour",
			startTime: 9 * time.Hour,
			endTime:   10 * time.Hour,
			expected:  2,
		},
		{
			name:      "exactly 24 hours",
			startTime: 0,
			endTime:   24 * time.Hour,
			expected:  48,
		},
		{
			name:      "non-aligned start (9:15)",
			startTime: 9*time.Hour + 15*time.Minute,
			endTime:   10 * time.Hour,
			expected:  1, // только один слот с 9:15 до 9:45
		},
		{
			name:      "non-aligned end (10:45)",
			startTime: 9 * time.Hour,
			endTime:   10*time.Hour + 45*time.Minute,
			expected:  3, // 9:00-9:30, 9:30-10:00, 10:00-10:30
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule := domain.NewSchedule(
				domain.WithScheduleRestoreSpecs(&domain.ScheduleRestoreSpecs{
					ID:         uuid.New(),
					RoomID:     roomID,
					StartTime:  tt.startTime,
					EndTime:    tt.endTime,
					DaysOfWeek: []time.Weekday{time.Monday},
				}),
			)
			date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			slots := domain.GenerateSlots(schedule, date)
			assert.Len(t, slots, tt.expected)
		})
	}
}

func TestNewSlot_WithInitSpecs(t *testing.T) {
	roomID := uuid.New()
	startTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	endTime := startTime.Add(domain.SlotDuration)

	slot := domain.NewSlot(
		domain.WithSlotInitSpecs(&domain.SlotInitSpecs{
			RoomID:    roomID,
			StartTime: startTime,
			EndTime:   endTime,
		}),
	)

	assert.Equal(t, roomID, slot.RoomID())
	assert.Equal(t, startTime, slot.StartTime())
	assert.Equal(t, endTime, slot.EndTime())
}

func TestNewSlot_WithRestoreSpecs(t *testing.T) {
	id := uuid.New()
	roomID := uuid.New()
	startTime := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	endTime := startTime.Add(domain.SlotDuration)

	slot := domain.NewSlot(
		domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
			ID:        id,
			RoomID:    roomID,
			StartTime: startTime,
			EndTime:   endTime,
		}),
	)

	assert.Equal(t, id, slot.ID())
	assert.Equal(t, roomID, slot.RoomID())
	assert.Equal(t, startTime, slot.StartTime())
	assert.Equal(t, endTime, slot.EndTime())
}
