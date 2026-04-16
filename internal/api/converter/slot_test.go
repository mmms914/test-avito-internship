package converter_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/domain"
)

func TestSlotToObject(t *testing.T) {
	slotID := uuid.New()
	roomID := uuid.New()
	startTime := time.Now().UTC()
	endTime := startTime.Add(30 * time.Minute)

	slot := domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
		ID:        slotID,
		RoomID:    roomID,
		StartTime: startTime,
		EndTime:   endTime,
	}))

	result := converter.SlotToObject(slot)

	assert.NotNil(t, result)
	assert.Equal(t, slotID, result.ID)
	assert.Equal(t, roomID, result.RoomID)
	assert.Equal(t, startTime, result.StartTime)
	assert.Equal(t, endTime, result.EndTime)
}

func TestSlotArrayToResponse(t *testing.T) {
	tests := map[string]struct {
		slots    []*domain.Slot
		expected int
	}{
		"multiple slots": {
			slots: []*domain.Slot{
				domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    uuid.New(),
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC().Add(30 * time.Minute),
				})),
				domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    uuid.New(),
					StartTime: time.Now().UTC().Add(1 * time.Hour),
					EndTime:   time.Now().UTC().Add(1*time.Hour + 30*time.Minute),
				})),
			},
			expected: 2,
		},
		"single slot": {
			slots: []*domain.Slot{
				domain.NewSlot(domain.WithSlotRestoreSpecs(&domain.SlotRestoreSpecs{
					ID:        uuid.New(),
					RoomID:    uuid.New(),
					StartTime: time.Now().UTC(),
					EndTime:   time.Now().UTC().Add(30 * time.Minute),
				})),
			},
			expected: 1,
		},
		"empty array": {
			slots:    []*domain.Slot{},
			expected: 0,
		},
		"nil array": {
			slots:    nil,
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := converter.SlotArrayToResponse(tt.slots)

			assert.NotNil(t, result)
			assert.Len(t, result.Slots, tt.expected)

			for i, slot := range tt.slots {
				if i < len(result.Slots) {
					assert.Equal(t, slot.ID(), result.Slots[i].ID)
					assert.Equal(t, slot.RoomID(), result.Slots[i].RoomID)
					assert.Equal(t, slot.StartTime(), result.Slots[i].StartTime)
					assert.Equal(t, slot.EndTime(), result.Slots[i].EndTime)
				}
			}
		})
	}
}
