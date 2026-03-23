package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/slot"
)

type SlotHandler struct {
	service *slot.Service
	logger  *slog.Logger
}

type SlotHandlerConfig struct {
	Service *slot.Service
	Logger  *slog.Logger
}

func NewSlotHandler(c *SlotHandlerConfig) *SlotHandler {
	return &SlotHandler{
		service: c.Service,
		logger:  c.Logger,
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
		handleServiceError(w, err, sh.logger)
		return
	}

	writeJSON(w, http.StatusCreated, converter.SlotArrayToResponse(slots))
}
