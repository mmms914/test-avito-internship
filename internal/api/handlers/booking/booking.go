package booking

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/api/handlers"
	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/internal/dto"
	"github.com/mmms914/test-avito-internship/internal/errs"
	"github.com/mmms914/test-avito-internship/pkg/ptr"
)

type Service interface {
	Create(ctx context.Context, model *dto.BookingCreateModel) (*domain.Booking, error)
	List(ctx context.Context, filter *dto.BookingFilter) ([]*domain.Booking, error)
	ListForUser(ctx context.Context) ([]*domain.Booking, error)
	Cancel(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error)
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

// Create godoc
// @Summary      Создать бронь
// @Description  Создание брони на слот (только user)
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateBookingRequest true "Данные брони"
// @Success      201 {object} models.BookingResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      409 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /bookings/create [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if err := handlers.CheckUserRole(r); err != nil {
		handlers.WriteError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	var req models.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "invalid request body", http.StatusBadRequest)
		return
	}

	if validationErrors := h.validate.Struct(req); validationErrors != nil {
		handlers.WriteValidationError(w, errors.Join(validationErrors))
		return
	}

	createdBooking, err := h.service.Create(r.Context(), &dto.BookingCreateModel{
		SlotID:               uuid.MustParse(*req.SlotID),
		CreateConferenceLink: req.CreateConferenceLink,
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, converter.BookingToResponse(createdBooking))
}

// List godoc
// @Summary      Все брони
// @Description  Получить список всех броней с пагинацией (только admin)
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Номер страницы" default(1)
// @Param        pageSize query int false "Размер страницы" default(20) maximum(100)
// @Success      200 {object} models.ListBookingResponseWithPag
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /bookings/list [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if err := handlers.CheckAdminRole(r); err != nil {
		handlers.WriteError(w, models.ForbiddenErrorCode, "admin role required", http.StatusForbidden)
		return
	}

	pageInt := domain.DefaultPage
	page := r.URL.Query().Get("page")
	if page != "" {
		pInt, err := strconv.Atoi(page)
		if err != nil || pInt < domain.MinPage {
			handlers.WriteError(w, models.InvalidRequestErrorCode,
				"page must be an integer and greater than 0", http.StatusBadRequest)
			return
		}

		pageInt = pInt
	}

	pageSizeInt := domain.DefaultPageSize
	pageSize := r.URL.Query().Get("pageSize")
	if pageSize != "" {
		pSizeInt, err := strconv.Atoi(pageSize)
		if err != nil || pSizeInt < domain.MinPageSize || pSizeInt > domain.MaxPageSize {
			handlers.WriteError(w, models.InvalidRequestErrorCode,
				"page size must be an integer and between 1 and 100", http.StatusBadRequest)
			return
		}

		pageSizeInt = pSizeInt
	}

	bookings, err := h.service.List(r.Context(), &dto.BookingFilter{
		Page:     ptr.To(pageInt),
		PageSize: ptr.To(pageSizeInt),
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, converter.BookingArrayToListResponseWithPag(
		bookings,
		&models.Pagination{
			Page:     pageInt,
			PageSize: pageSizeInt,
			Total:    len(bookings),
		}))
}

// ListMy godoc
// @Summary      Мои брони
// @Description  Получить список броней текущего пользователя (только user)
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.ListBookingResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /bookings/my [get]
func (h *Handler) ListMy(w http.ResponseWriter, r *http.Request) {
	if err := handlers.CheckUserRole(r); err != nil {
		handlers.WriteError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	bookings, err := h.service.ListForUser(r.Context())
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, converter.BookingArrayToListResponse(bookings))
}

// Cancel godoc
// @Summary      Отменить бронь
// @Description  Отмена брони (только user, идемпотентно)
// @Tags         bookings
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        bookingId path string true "ID брони"
// @Success      200 {object} models.BookingResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      401 {object} models.ErrorResponse
// @Failure      403 {object} models.ErrorResponse
// @Failure      404 {object} models.ErrorResponse
// @Failure      500 {object} models.ErrorResponse
// @Router       /bookings/{bookingId}/cancel [post]
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	if err := handlers.CheckUserRole(r); err != nil {
		handlers.WriteError(w, models.ForbiddenErrorCode, "user role required", http.StatusForbidden)
		return
	}

	bookingIDStr := chi.URLParam(r, "bookingId")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		handlers.WriteError(w, models.InvalidRequestErrorCode, "bookingId must be uuid", http.StatusBadRequest)
		return
	}

	cancelledBooking, err := h.service.Cancel(r.Context(), bookingID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	handlers.WriteJSON(w, http.StatusOK, converter.BookingToResponse(cancelledBooking))
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errs.ErrRoomNotExists):
		handlers.WriteError(w, models.RoomNotFoundErrorCode, "room not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrForbidden):
		handlers.WriteError(w, models.ForbiddenErrorCode, "access denied", http.StatusForbidden)
	case errors.Is(err, errs.ErrUnauthorized):
		handlers.WriteError(w, models.UnauthorizedErrorCode, "unauthorized", http.StatusUnauthorized)
	case errors.Is(err, errs.ErrUserNotFound):
		handlers.WriteError(w, models.NotFoundErrorCode, "user not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrBookingNotFound):
		handlers.WriteError(w, models.BookingNotFoundErrorCode, "booking not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrSlotNotFound):
		handlers.WriteError(w, models.SlotNotFoundErrorCode, "slot not found", http.StatusNotFound)
	case errors.Is(err, errs.ErrSlotAlreadyBooked):
		handlers.WriteError(w, models.SlotAlreadyBookedErrorCode, "slot already booked", http.StatusConflict)
	case errors.Is(err, errs.ErrSlotTimeInPast):
		handlers.WriteError(w, models.InvalidRequestErrorCode, "slot time in past", http.StatusBadRequest)
	default:
		h.logger.Error("Unexpected error", "error", err)
		handlers.WriteError(w, models.InternalErrorCode, "internal server error", http.StatusInternalServerError)
	}
}
