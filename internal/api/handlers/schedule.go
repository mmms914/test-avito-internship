package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/schedule"
)

type ScheduleHandler struct {
	service *schedule.Service
	logger  *slog.Logger
}

func NewScheduleHandler(service *schedule.Service, logger *slog.Logger) *ScheduleHandler {
	return &ScheduleHandler{
		service: service,
		logger:  logger,
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := validate.Struct(req); validationErrors != nil {
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
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, converter.ScheduleToResponse(createdSchedule))
}
