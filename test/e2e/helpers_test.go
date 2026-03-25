//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Room struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Capacity    int       `json:"capacity,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Schedule struct {
	ID         string `json:"id"`
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

type Slot struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}

type Booking struct {
	ID             string    `json:"id"`
	SlotID         string    `json:"slotId"`
	UserID         string    `json:"userId"`
	Status         string    `json:"status"`
	ConferenceLink *string   `json:"conferenceLink,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (s *E2ETestSuite) createRoom() *Room {
	reqBody := map[string]interface{}{
		"name":        fmt.Sprintf("E2E Test Room %s", uuid.New().String()[:8]),
		"description": "Room for E2E testing",
		"capacity":    10,
	}

	resp := s.doRequest("POST", "/rooms/create", s.adminToken, reqBody)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusCreated)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var result Room
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	return &result
}

func (s *E2ETestSuite) createSchedule(roomID string) *Schedule {
	reqBody := map[string]interface{}{
		"daysOfWeek": []int{1, 2, 3, 4, 5},
		"startTime":  "09:00",
		"endTime":    "18:00",
	}

	resp := s.doRequest("POST", fmt.Sprintf("/rooms/%s/schedule/create", roomID), s.adminToken, reqBody)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusCreated)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var result Schedule
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	return &result
}

func (s *E2ETestSuite) getAvailableSlots(roomID, date string) []Slot {
	resp := s.doRequest("GET", fmt.Sprintf("/rooms/%s/slots/list?date=%s", roomID, date), s.userToken, nil)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusOK)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result struct {
		Slots []Slot `json:"slots"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	return result.Slots
}

func (s *E2ETestSuite) createBooking(slotID string, createConferenceLink bool) *Booking {
	reqBody := map[string]interface{}{
		"slotId":               slotID,
		"createConferenceLink": createConferenceLink,
	}

	resp := s.doRequest("POST", "/bookings/create", s.userToken, reqBody)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusCreated)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	var result struct {
		Booking Booking `json:"booking"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	return &result.Booking
}

func (s *E2ETestSuite) cancelBooking(bookingID string) *Booking {
	resp := s.doRequest("POST", fmt.Sprintf("/bookings/%s/cancel", bookingID), s.userToken, nil)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusOK)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result struct {
		Booking Booking `json:"booking"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	return &result.Booking
}

func (s *E2ETestSuite) getMyBookings() []Booking {
	resp := s.doRequest("GET", "/bookings/my", s.userToken, nil)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusOK)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result struct {
		Bookings []Booking `json:"bookings"`
	}
	err := json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)

	return result.Bookings
}
