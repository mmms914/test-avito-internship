package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/auth"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/user"
)

type AuthHandler struct {
	logger          *slog.Logger
	service         *user.Service
	validate        *validator.Validate
	jwtSecret       string
	expirationHours int
}

type AuthHandlerConfig struct {
	Logger          *slog.Logger
	Service         *user.Service
	Validate        *validator.Validate
	JwtSecret       string
	ExpirationHours int
}

func NewAuthHandler(c *AuthHandlerConfig) *AuthHandler {
	return &AuthHandler{
		logger:          c.Logger,
		service:         c.Service,
		validate:        c.Validate,
		jwtSecret:       c.JwtSecret,
		expirationHours: c.ExpirationHours,
	}
}

func (ah *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req models.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := ah.validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	var userID string
	switch *req.Role {
	case "admin":
		userID = auth.DefaultAdminUID
	case "user":
		userID = auth.DefaultUserUID
	default:
		writeError(w, models.InvalidRequestErrorCode, "invalid role", http.StatusBadRequest)
	}

	token, err := auth.GenerateToken(userID, *req.Role, ah.jwtSecret, ah.expirationHours)
	if err != nil {
		writeError(w, models.InternalErrorCode, "cannot generate jwt token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, models.LoginResponse{Token: token})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := ah.validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	existingUser, err := ah.service.Login(r.Context(), &dto.UserCredentials{
		Email:    *req.Email,
		Password: *req.Password,
	})
	if err != nil {
		handleServiceError(w, err, ah.logger)
		return
	}

	token, err := auth.GenerateToken(existingUser.ID().String(),
		existingUser.Role().String(),
		ah.jwtSecret,
		ah.expirationHours)
	if err != nil {
		writeError(w, models.InternalErrorCode, "cannot generate jwt token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, models.LoginResponse{Token: token})
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := ah.validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	role, err := converter.UserRoleFromString(*req.Role)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid role", http.StatusBadRequest)
		return
	}

	existingUser, err := ah.service.Register(r.Context(), &dto.UserCreateModel{
		Email:    *req.Email,
		Password: *req.Password,
		Role:     role,
	})
	if err != nil {
		handleServiceError(w, err, ah.logger)
		return
	}

	writeJSON(w, http.StatusCreated, converter.UserToResponse(existingUser))
}
