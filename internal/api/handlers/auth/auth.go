package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/auth"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Service interface {
	Login(ctx context.Context, credentials *dto.UserCredentials) (*domain.User, error)
	Register(ctx context.Context, credentials *dto.UserCreateModel) (*domain.User, error)
}

type Validator interface {
	Struct(s any) error
}

type Logger interface {
	Error(msg string, args ...any)
}

type Handler struct {
	logger          Logger
	service         Service
	validate        Validator
	jwtSecret       string
	expirationHours int
}

type HandlerConfig struct {
	Logger          Logger
	Service         Service
	Validate        Validator
	JwtSecret       string
	ExpirationHours int
}

func NewHandler(c *HandlerConfig) *Handler {
	return &Handler{
		logger:          c.Logger,
		service:         c.Service,
		validate:        c.Validate,
		jwtSecret:       c.JwtSecret,
		expirationHours: c.ExpirationHours,
	}
}

func (h *Handler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req models.DummyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validate.Struct(req); validationErrors != nil {
		handlers.WriteValidationError(w, errors.Join(validationErrors))
		return
	}

	var userID string
	switch *req.Role {
	case "admin":
		userID = auth.DefaultAdminUID
	case "user":
		userID = auth.DefaultUserUID
	default:
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid role", http.StatusBadRequest)
	}

	token, err := auth.GenerateToken(userID, *req.Role, h.jwtSecret, h.expirationHours)
	if err != nil {
		handlers.WriteError(w, models.InternalErrorCode, "cannot generate jwt token", http.StatusInternalServerError)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, models.LoginResponse{Token: token})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validate.Struct(req); validationErrors != nil {
		handlers.WriteValidationError(w, errors.Join(validationErrors))
		return
	}

	existingUser, err := h.service.Login(r.Context(), &dto.UserCredentials{
		Email:    *req.Email,
		Password: *req.Password,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	token, err := auth.GenerateToken(existingUser.ID().String(),
		existingUser.Role().String(),
		h.jwtSecret,
		h.expirationHours)
	if err != nil {
		handlers.WriteError(w, models.InternalErrorCode, "cannot generate jwt token", http.StatusInternalServerError)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, models.LoginResponse{Token: token})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validate.Struct(req); validationErrors != nil {
		handlers.WriteValidationError(w, errors.Join(validationErrors))
		return
	}

	role, err := converter.UserRoleFromString(*req.Role)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid role", http.StatusBadRequest)
		return
	}

	existingUser, err := h.service.Register(r.Context(), &dto.UserCreateModel{
		Email:    *req.Email,
		Password: *req.Password,
		Role:     role,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, converter.UserToResponse(existingUser))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrForbidden):
		handlers.WriteError(w, models.ForbiddenErrorCode, "access denied", http.StatusForbidden)
	case errors.Is(err, errs.ErrUnauthorized):
		handlers.WriteError(w, models.UnauthorizedErrorCode, "unauthorized", http.StatusUnauthorized)
	case errors.Is(err, errs.ErrUserNotFound):
		handlers.WriteError(w, models.NotFoundErrorCode, "user not found", http.StatusNotFound)
	default:
		h.logger.Error("Unexpected error", "error", err)
		handlers.WriteError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}
