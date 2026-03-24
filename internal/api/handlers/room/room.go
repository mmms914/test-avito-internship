package room

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

type Service interface {
	Create(ctx context.Context, specs *domain.RoomInitSpecs) (*domain.Room, error)
	GetAll(ctx context.Context) ([]*domain.Room, error)
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
		handlers.WriteError(w, models.ForbiddenErrorCode,
			fmt.Sprintf("admin role required: %v", err), http.StatusForbidden)
		return
	}

	var req models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validate.Struct(req); validationErrors != nil {
		handlers.WriteValidationError(w, errors.Join(validationErrors))
		return
	}

	createdRoom, err := h.service.Create(r.Context(), &domain.RoomInitSpecs{
		Name:        *req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, converter.RoomToResponse(createdRoom))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.GetAll(r.Context())
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, converter.RoomArrayToResponse(rooms))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrForbidden):
		handlers.WriteError(w, models.ForbiddenErrorCode, "access denied", http.StatusForbidden)
	case errors.Is(err, errs.ErrUnauthorized):
		handlers.WriteError(w, models.UnauthorizedErrorCode, "unauthorized", http.StatusUnauthorized)
	default:
		h.logger.Error("Unexpected error", "error", err)
		handlers.WriteError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}
