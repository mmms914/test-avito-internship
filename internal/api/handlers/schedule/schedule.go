package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Service interface {
	Create(ctx context.Context, specs *domain.ScheduleInitSpecs) (*domain.Schedule, error)
}

type Validator interface {
	Struct(s any) error
}

type Logger interface {
	Error(msg string, args ...any)
}

type Handler struct {
	service  Service
	validate Validator
	logger   Logger
}

type HandlerConfig struct {
	Service  Service
	Validate Validator
	Logger   Logger
}

func NewHandler(c *HandlerConfig) *Handler {
	return &Handler{
		service:  c.Service,
		validate: c.Validate,
		logger:   c.Logger,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := handlers.CheckAdminRole(r); err != nil {
		handlers.WriteError(w, models.ForbiddenErrorCode, "admin role required", http.StatusForbidden)
		return
	}

	roomIDStr := chi.URLParam(r, "roomID")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "roomID must be UUID", http.StatusBadRequest)
		return
	}

	var req models.CreateScheduleRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validate.Struct(req); validationErrors != nil {
		handlers.WriteValidationError(w, errors.Join(validationErrors))
		return
	}

	days, err := converter.IntArrayToWeekdays(req.DaysOfWeek)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid days of week", http.StatusBadRequest)
		return
	}

	startTime, err := converter.StringHourMinuteToTime(req.StartTime)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid start time", http.StatusBadRequest)
		return
	}

	endTime, err := converter.StringHourMinuteToTime(req.EndTime)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid start time", http.StatusBadRequest)
		return
	}

	createdSchedule, err := h.service.Create(r.Context(), &domain.ScheduleInitSpecs{
		RoomID:     roomID,
		DaysOfWeek: days,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, converter.ScheduleToResponse(createdSchedule))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrRoomNotExists):
		handlers.WriteError(w, models.RoomNotFoundErrorCode, "room not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrForbidden):
		handlers.WriteError(w, models.ForbiddenErrorCode, "access denied", http.StatusForbidden)
	case errors.Is(err, errs.ErrUnauthorized):
		handlers.WriteError(w, models.UnauthorizedErrorCode, "unauthorized", http.StatusUnauthorized)
	case errors.Is(err, errs.ErrScheduleExists):
		handlers.WriteError(w, models.ScheduleExistsErrorCode, "schedule already exists", http.StatusConflict)
	default:
		h.logger.Error("Unexpected error", "error", err)
		handlers.WriteError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}
