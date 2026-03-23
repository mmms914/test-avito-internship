package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/converter"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/service/booking"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
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
		handleServiceError(w, err, bh.logger)
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
		handleServiceError(w, err, bh.logger)
		return
	}

	writeJSON(w, http.StatusOK, models.ListBookingResponseWithPagination{
		Bookings: converter.BookingArrayToResponse(bookings),
		Pagination: &models.Pagination{
			Page:     pageInt,
			PageSize: pageSizeInt,
			Total:    len(bookings),
		},
	})
}

func (bh *BookingHandler) ListMy(w http.ResponseWriter, r *http.Request) {
	if err := checkUserRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	bookings, err := bh.service.ListForUser(r.Context())
	if err != nil {
		handleServiceError(w, err, bh.logger)
		return
	}

	writeJSON(w, http.StatusOK, converter.BookingArrayToResponse(bookings))
}

func (bh *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	if err := checkUserRole(r); err != nil {
		writeError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	bookingIDStr := r.URL.Query().Get("bookingID")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		writeError(w, models.InvalidRequestErrorCode, "bookingID must be uuid", http.StatusBadRequest)
		return
	}

	cancelledBooking, err := bh.service.Cancel(r.Context(), bookingID)
	if err != nil {
		handleServiceError(w, err, bh.logger)
		return
	}

	writeJSON(w, http.StatusOK, converter.BookingToResponse(cancelledBooking))
}
