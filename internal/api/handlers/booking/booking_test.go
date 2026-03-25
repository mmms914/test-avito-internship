package booking_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/booking"
	mocks "github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/booking/mocks"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

var errInternal = errors.New("internal error")

func TestHandler_Create(t *testing.T) {
	tests := map[string]struct {
		requestBody    interface{}
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - create booking": {
			requestBody: models.CreateBookingRequest{
				SlotID:               ptr.To(uuid.New().String()),
				CreateConferenceLink: ptr.To(false),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*dto.BookingCreateModel")).Return(&domain.Booking{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
		},
		"failure - invalid request body": {
			requestBody:    "invalid json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - missing slotId": {
			requestBody: models.CreateBookingRequest{
				CreateConferenceLink: ptr.To(false),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - slot not found": {
			requestBody: models.CreateBookingRequest{
				SlotID:               ptr.To(uuid.New().String()),
				CreateConferenceLink: ptr.To(false),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*dto.BookingCreateModel")).Return(nil, errs.ErrSlotNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   models.SlotNotFoundErrorCode,
		},
		"failure - slot already booked": {
			requestBody: models.CreateBookingRequest{
				SlotID:               ptr.To(uuid.New().String()),
				CreateConferenceLink: ptr.To(false),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*dto.BookingCreateModel")).Return(nil, errs.ErrSlotAlreadyBooked)
			},
			expectedStatus: http.StatusConflict,
			expectedCode:   models.SlotAlreadyBookedErrorCode,
		},
		"failure - slot time in past": {
			requestBody: models.CreateBookingRequest{
				SlotID:               ptr.To(uuid.New().String()),
				CreateConferenceLink: ptr.To(false),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*dto.BookingCreateModel")).Return(nil, errs.ErrSlotTimeInPast)
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - internal error": {
			requestBody: models.CreateBookingRequest{
				SlotID:               ptr.To(uuid.New().String()),
				CreateConferenceLink: ptr.To(false),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*dto.BookingCreateModel")).Return(nil, errInternal)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   models.InternalErrorCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			validate := validator.New()
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			handler := booking.NewHandler(&booking.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(setContext(req.Context(), domain.UserRole))

			rr := httptest.NewRecorder()
			handler.Create(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_List(t *testing.T) {
	tests := map[string]struct {
		queryParams    string
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - list bookings": {
			queryParams: "page=1&pageSize=10",
			setupMock: func(m *mocks.Service) {
				m.On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).Return([]*domain.Booking{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		"success - default pagination": {
			queryParams: "",
			setupMock: func(m *mocks.Service) {
				m.On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).Return([]*domain.Booking{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		"failure - service error": {
			queryParams: "",
			setupMock: func(m *mocks.Service) {
				m.On("List", mock.Anything, mock.AnythingOfType("*dto.BookingFilter")).Return(nil, errInternal)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   models.InternalErrorCode,
		},
		"failure - invalid page": {
			queryParams:    "page=0&pageSize=10",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid pageSize": {
			queryParams:    "page=1&pageSize=200",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - negative page": {
			queryParams:    "page=-1&pageSize=10",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid page format": {
			queryParams:    "page=abc&pageSize=10",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			validate := validator.New()
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			handler := booking.NewHandler(&booking.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			req := httptest.NewRequest(http.MethodGet, "/bookings/list?"+tt.queryParams, nil)
			req = req.WithContext(setContext(req.Context(), domain.AdminRole))

			rr := httptest.NewRecorder()
			handler.List(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_ListMy(t *testing.T) {
	tests := map[string]struct {
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - list my bookings": {
			setupMock: func(m *mocks.Service) {
				m.On("ListForUser", mock.Anything).Return([]*domain.Booking{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		"failure - internal error": {
			setupMock: func(m *mocks.Service) {
				m.On("ListForUser", mock.Anything).Return(nil, errInternal)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   models.InternalErrorCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			validate := validator.New()
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			handler := booking.NewHandler(&booking.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
			req = req.WithContext(setContext(req.Context(), domain.UserRole))

			rr := httptest.NewRecorder()
			handler.ListMy(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_Cancel(t *testing.T) {
	tests := map[string]struct {
		bookingID      string
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - cancel booking": {
			bookingID: uuid.New().String(),
			setupMock: func(m *mocks.Service) {
				m.On("Cancel", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&domain.Booking{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		"failure - invalid booking id": {
			bookingID:      "invalid-uuid",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - booking not found": {
			bookingID: uuid.New().String(),
			setupMock: func(m *mocks.Service) {
				m.On("Cancel", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errs.ErrBookingNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   models.BookingNotFoundErrorCode,
		},
		"failure - forbidden": {
			bookingID: uuid.New().String(),
			setupMock: func(m *mocks.Service) {
				m.On("Cancel", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errs.ErrForbidden)
			},
			expectedStatus: http.StatusForbidden,
			expectedCode:   models.ForbiddenErrorCode,
		},
		"failure - internal error": {
			bookingID: uuid.New().String(),
			setupMock: func(m *mocks.Service) {
				m.On("Cancel", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errInternal)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   models.InternalErrorCode,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			validate := validator.New()
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			handler := booking.NewHandler(&booking.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			req := httptest.NewRequest(http.MethodPost, "/bookings/"+tt.bookingID+"/cancel", nil)
			req = req.WithContext(setContext(req.Context(), domain.UserRole))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("bookingId", tt.bookingID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.Cancel(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func setContext(ctx context.Context, role domain.Role) context.Context {
	return context.WithValue(context.WithValue(ctx, domain.UserRoleKey, role), domain.UserIDKey, uuid.New())
}
