//go:build e2e
// +build e2e

package e2e

import (
	"time"
)

func (s *E2ETestSuite) TestCreateBooking() {
	// 1. Создаем переговорку (admin)
	room := s.createRoom()
	s.Require().NotNil(room)

	// 2. Создаем расписание (admin)
	schedule := s.createSchedule(room.ID)
	s.Require().NotNil(schedule)

	// 3. Получаем доступные слоты (user)
	date := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")
	slots := s.getAvailableSlots(room.ID, date)
	s.Require().NotEmpty(slots, "No slots available")

	// 4. Создаем бронь на первый слот (user)
	booking := s.createBooking(slots[0].ID, false)
	s.Require().NotNil(booking)
	s.Require().Equal("active", booking.Status)

	// 5. Проверяем, что слот больше не доступен
	slotsAfterBooking := s.getAvailableSlots(room.ID, date)
	for _, slot := range slotsAfterBooking {
		s.Require().NotEqual(booking.SlotID, slot.ID, "Slot should not be available after booking")
	}

	// 6. Проверяем список броней пользователя
	userBookings := s.getMyBookings()
	found := false
	for _, b := range userBookings {
		if b.ID == booking.ID {
			found = true
			s.Require().Equal("active", b.Status)
			break
		}
	}
	s.Require().True(found, "Booking not found in user's list")
}

func (s *E2ETestSuite) TestCancelBooking() {
	// 1. Создаем переговорку и расписание
	room := s.createRoom()
	s.Require().NotNil(room)

	schedule := s.createSchedule(room.ID)
	s.Require().NotNil(schedule)

	// 2. Получаем слоты и создаем бронь
	date := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")
	slots := s.getAvailableSlots(room.ID, date)
	s.Require().NotEmpty(slots, "No slots available")

	booking := s.createBooking(slots[0].ID, false)
	s.Require().NotNil(booking)

	// 3. Отменяем бронь
	cancelledBooking := s.cancelBooking(booking.ID)
	s.Require().NotNil(cancelledBooking)
	s.Require().Equal("cancelled", cancelledBooking.Status)

	// 4. Проверяем, что слот снова доступен
	slotsAfterCancel := s.getAvailableSlots(room.ID, date)
	slotAvailable := false
	for _, slot := range slotsAfterCancel {
		if slot.ID == booking.SlotID {
			slotAvailable = true
			break
		}
	}
	s.Require().True(slotAvailable, "Slot should be available after cancellation")

	// 5. Проверяем идемпотентность: повторная отмена не должна быть ошибкой
	cancelledAgain := s.cancelBooking(booking.ID)
	s.Require().NotNil(cancelledAgain)
	s.Require().Equal("cancelled", cancelledAgain.Status)

	// 6. Проверяем, что бронь помечена как cancelled в списке пользователя
	userBookings := s.getMyBookings()
	found := false
	for _, b := range userBookings {
		if b.ID == booking.ID {
			found = true
			break
		}
	}
	s.Require().False(found, "Booking found in user's list after cancellation")
}
