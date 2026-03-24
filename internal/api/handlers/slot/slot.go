package slot

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Service interface {
	GetAvailableSlots(ctx context.Context, roomID uuid.UUID, date time.Time) ([]*domain.Slot, error)
}

type Logger interface {
	Error(msg string, args ...any)
}

type Handler struct {
	service Service
	logger  Logger
}

type HandlerConfig struct {
	Service Service
	Logger  Logger
}

func NewHandler(c *HandlerConfig) *Handler {
	return &Handler{
		service: c.Service,
		logger:  c.Logger,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomID")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "roomID must be UUID", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.DateOnly, r.URL.Query().Get("date"))
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "date must be valid", http.StatusBadRequest)
		return
	}

	slots, err := h.service.GetAvailableSlots(r.Context(), roomID, date)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, converter.SlotArrayToResponse(slots))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrRoomNotExists):
		handlers.WriteError(w, models.RoomNotFoundErrorCode, "room not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrForbidden):
		handlers.WriteError(w, models.ForbiddenErrorCode, "access denied", http.StatusForbidden)
	case errors.Is(err, errs.ErrUnauthorized):
		handlers.WriteError(w, models.UnauthorizedErrorCode, "unauthorized", http.StatusUnauthorized)
	default:
		h.logger.Error("Unexpected error", "error", err)
		handlers.WriteError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}
