package handlers

import (
	"encoding/json"
	"errors"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/booking"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
	"log/slog"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	service *booking.Service
	logger  *slog.Logger
}

func NewBookingHandler(service *booking.Service, logger *slog.Logger) *BookingHandler {
	return &BookingHandler{
		service: service,
		logger:  logger,
	}
}

func (bh *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := checkUserRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	var req models.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := validate.Struct(req); validationErrors != nil {
		writeValidationError(w, errors.Join(validationErrors))
		return
	}

	createdBooking, err := bh.service.Create(r.Context(), &dto.BookingCreateModel{
		SlotID:               *req.SlotID,
		CreateConferenceLink: req.CreateConferenceLink,
	})
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, converter.BookingToResponse(createdBooking))
}

func (bh *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	if err := checkAdminRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "admin role required", http.StatusForbidden)
		return
	}

	pageInt := domain.DefaultPage
	page := r.URL.Query().Get("page")
	if page != "" {
		pInt, err := strconv.Atoi(page)
		if err != nil || pInt < domain.MinPage {
			writeError(w, models.InvalidRequestErrorCode, "page must be an integer and greater than 0", http.StatusBadRequest)
			return
		}

		pageInt = pInt
	}

	pageSizeInt := domain.DefaultPageSize
	pageSize := r.URL.Query().Get("pageSize")
	if pageSize != "" {
		pSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pSizeInt < domain.MinPageSize || pSizeInt > domain.MaxPageSize {
			writeError(w, models.InvalidRequestErrorCode, "page size must be an integer and between 1 and 100", http.StatusBadRequest)
			return
		}

		pageSizeInt = pSizeInt
	}

	bookings, err := bh.service.List(r.Context(), &dto.BookingFilter{
		Page:     ptr.To(pageInt),
		PageSize: ptr.To(pageSizeInt),
	})
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, converter.BookingArrayToResponse(bookings))
}
