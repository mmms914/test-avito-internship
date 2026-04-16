package slot

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/api/handlers"
	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
)

type Service interface {
	GetAvailableSlots(ctx context.Context, f *dto.SlotFilter) ([]*domain.Slot, error)
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

// List godoc
// @Summary      Доступные слоты
// @Description  Получить доступные для бронирования слоты по переговорке и дате
// @Tags         slots
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        roomId path string true "ID переговорки"
// @Param        date query string true "Дата в формате YYYY-MM-DD"
// @Success      200 {object} models.SlotsResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /rooms/{roomId}/slots/list [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	roomIDStr := chi.URLParam(r, "roomId")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "roomId must be UUID", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.DateOnly, r.URL.Query().Get("date"))
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "date must be valid", http.StatusBadRequest)
		return
	}

	f := &dto.SlotFilter{
		RoomID: roomID,
		Date:   date,
	}

	slots, err := h.service.GetAvailableSlots(r.Context(), f)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, converter.SlotArrayToResponse(slots))
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
