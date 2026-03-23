package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

var validate = validator.New()

func handleServiceError(w http.ResponseWriter, err error, logger *slog.Logger) {
	switch {
	case errors.Is(err, errs.ErrRoomNotExists):
		writeError(w, models.RoomNotFoundErrorCode, "room not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrForbidden):
		writeError(w, models.ForbiddenErrorCode, "access denied", http.StatusForbidden)
	case errors.Is(err, errs.ErrUnauthorized):
		writeError(w, models.UnauthorizedErrorCode, "unauthorized", http.StatusUnauthorized)
	case errors.Is(err, errs.ErrScheduleExists):
		writeError(w, models.ScheduleExistsErrorCode, "schedule already exists", http.StatusConflict)
	case errors.Is(err, errs.ErrUserNotFound):
		writeError(w, models.NotFoundErrorCode, "user not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrRoomNotExists):
		writeError(w, models.RoomNotFoundErrorCode, "room not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrBookingNotFound):
		writeError(w, models.BookingNotFoundErrorCode, "booking not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrSlotNotFound):
		writeError(w, models.SlotNotFoundErrorCode, "slot not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrSlotAlreadyBooked):
		writeError(w, models.SlotAlreadyBookedErrorCode, "slot already exists", http.StatusConflict)
	default:
		logger.Error("Unexpected error", "error", err)
		writeError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: &models.Error{
			Code:    code,
			Message: message,
		},
	})
}

func writeValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: &models.Error{
			Code:    models.InvalidRequestErrorCode,
			Message: fmt.Sprintf("validation not passed: %v", err),
		},
	})
}

func checkUserRole(r *http.Request) error {
	auth, err := domain.GetCredentialsFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("getting credentials: %w", err)
	}

	if auth.Role != domain.UserRole {
		return fmt.Errorf("user is not an user role")
	}
	return nil
}

func checkAdminRole(r *http.Request) error {
	auth, err := domain.GetCredentialsFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("getting credentials: %w", err)
	}

	if auth.Role != domain.AdminRole {
		return fmt.Errorf("user is not an admin role")
	}
	return nil
}
