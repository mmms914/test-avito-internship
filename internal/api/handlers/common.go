package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: &models.Error{
			Code:    code,
			Message: message,
		},
	})
}

func WriteValidationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(models.ErrorResponse{
		Error: &models.Error{
			Code:    models.InvalidRequestErrorCode,
			Message: fmt.Sprintf("validation not passed: %v", err),
		},
	})
}

func CheckUserRole(r *http.Request) error {
	auth, err := domain.GetCredentialsFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("getting credentials: %w", err)
	}

	if auth.Role != domain.UserRole {
		return errors.New("user is not an user role")
	}
	return nil
}

func CheckAdminRole(r *http.Request) error {
	auth, err := domain.GetCredentialsFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("getting credentials: %w", err)
	}

	if auth.Role != domain.AdminRole {
		return errors.New("user is not an admin role")
	}
	return nil
}
