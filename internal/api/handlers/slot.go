package handlers

import (
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/slot"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

type SlotHandler struct {
	service *slot.Service
	logger  *slog.Logger
}

func NewSlotHandler(service *slot.Service, logger *slog.Logger) *SlotHandler {
	return &SlotHandler{
		service: service,
		logger:  logger,
	}
}

func (sh *SlotHandler) List(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomID")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "roomID must be UUID", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.DateOnly, r.URL.Query().Get("date"))
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "date must be valid", http.StatusBadRequest)
		return
	}

	slots, err := sh.service.GetAvailableSlots(r.Context(), roomID, date)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, converter.SlotArrayToResponse(slots))
}
