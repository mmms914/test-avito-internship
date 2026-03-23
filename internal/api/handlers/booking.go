package handlers

import (
	"encoding/json"
	"errors"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/booking"
	"log/slog"
	"net/http"
)

type BookingHandler struct {
	service *booking.Service
	logger  *slog.Logger
}

func NewBookingHandler(service *booking.Service, logger *slog.Logger) *BookingHandler {
	return &BookingHandler{
		service: service,
		logger:  logger,
	}
}

func (bh *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := checkUserRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	var req models.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	createdBooking, err := bh.service.Create(r.Context(), &dto.BookingCreateModel{
		SlotID:               *req.SlotID,
		CreateConferenceLink: req.CreateConferenceLink,
	})
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, converter.BookingToResponse(createdBooking))
}

func (rh *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	rooms, err := rh.service.GetAll(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.RoomArrayToResponse(rooms))
}
