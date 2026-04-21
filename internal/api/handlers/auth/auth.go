package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mmms914/test-avito-internship/internal/api/auth"
	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/api/handlers"
	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
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

// DummyLogin godoc
// @Summary      Тестовый логин
// @Description  Получить JWT токен по роли (admin/user)
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.DummyLoginRequest true "Роль пользователя"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /dummyLogin [post]
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
		return
	}

	token, err := auth.GenerateToken(userID, *req.Role, h.jwtSecret, h.expirationHours)
	if err != nil {
		handlers.WriteError(w, models.InternalErrorCode, "cannot generate jwt token", http.StatusInternalServerError)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, models.LoginResponse{Token: token})
}

// Login godoc
// @Summary      Авторизация
// @Description  Авторизация по email и паролю
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Учетные данные"
// @Success      200 {object} models.LoginResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /login [post]
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

// Register godoc
// @Summary      Регистрация
// @Description  Регистрация нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Данные для регистрации"
// @Success      201 {object} models.UserResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /register [post]
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
		handlers.WriteError(w, models.UnauthorizedErrorCode, "user not found", http.StatusUnauthorized)
	case errors.Is(err, errs.ErrUserAlreadyExists):
		handlers.WriteError(w, models.InvalidRequestErrorCode, "user already exists", http.StatusBadRequest)
	default:
		h.logger.Error("Unexpected error", "error", err)
		handlers.WriteError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}
