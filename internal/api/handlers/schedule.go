package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/schedule"
)

type ScheduleHandler struct {
	service  *schedule.Service
	validate *validator.Validate
	logger   *slog.Logger
}

type ScheduleHandlerConfig struct {
	Service  *schedule.Service
	Validate *validator.Validate
	Logger   *slog.Logger
}

func NewScheduleHandler(c *ScheduleHandlerConfig) *ScheduleHandler {
	return &ScheduleHandler{
		service:  c.Service,
		validate: c.Validate,
		logger:   c.Logger,
	}
}

func (sh *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := checkAdminRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "admin role required", http.StatusForbidden)
		return
	}

	roomIDStr := chi.URLParam(r, "roomID")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "roomID must be UUID", http.StatusBadRequest)
		return
	}

	var req models.CreateScheduleRequest
	if decodeErr := json.NewDecoder(r.Body).Decode(&req); decodeErr != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := sh.validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	days, err := converter.IntArrayToWeekdays(req.DaysOfWeek)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid days of week", http.StatusBadRequest)
		return
	}

	startTime, err := converter.StringHourMinuteToTime(req.StartTime)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid start time", http.StatusBadRequest)
		return
	}

	endTime, err := converter.StringHourMinuteToTime(req.EndTime)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid start time", http.StatusBadRequest)
		return
	}

	createdSchedule, err := sh.service.Create(r.Context(), &domain.ScheduleInitSpecs{
		RoomID:     roomID,
		DaysOfWeek: days,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		handleServiceError(w, err, sh.logger)
		return
	}

	writeJSON(w, http.StatusCreated, converter.ScheduleToResponse(createdSchedule))
}
