// internal/api/handlers/room.go
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/room"
)

type RoomHandler struct {
	service *room.Service
	logger  *slog.Logger
}

func NewRoomHandler(service *room.Service, logger *slog.Logger) *RoomHandler {
	return &RoomHandler{
		service: service,
		logger:  logger,
	}
}

func (rh *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := checkAdminRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "admin role required", http.StatusForbidden)
		return
	}

	var req models.CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	createdRoom, err := rh.service.Create(r.Context(), &domain.RoomInitSpecs{
		Name:        *req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	})
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, converter.RoomToResponse(createdRoom))
}

func (rh *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	rooms, err := rh.service.GetAll(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.RoomArrayToResponse(rooms))
}
