package slot_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/slot"
	mocks "github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/slot/mocks"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
)

func TestHandler_List(t *testing.T) {
	tests := map[string]struct {
		roomID         string
		queryDate      string
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - list available slots": {
			roomID:    uuid.New().String(),
			queryDate: time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02"),
			setupMock: func(m *mocks.Service) {
				m.On("GetAvailableSlots", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).Return([]*domain.Slot{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		"failure - invalid roomId": {
			roomID:         "invalid-uuid",
			queryDate:      time.Now().UTC().Format("2006-01-02"),
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - missing date": {
			roomID:         uuid.New().String(),
			queryDate:      "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid date format": {
			roomID:         uuid.New().String(),
			queryDate:      "2024-15-01",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - room not found": {
			roomID:    uuid.New().String(),
			queryDate: time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02"),
			setupMock: func(m *mocks.Service) {
				m.On("GetAvailableSlots", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).Return(nil, errs.ErrRoomNotExists)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   models.RoomNotFoundErrorCode,
		},
		"failure - internal error": {
			roomID:    uuid.New().String(),
			queryDate: time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02"),
			setupMock: func(m *mocks.Service) {
				m.On("GetAvailableSlots", mock.Anything, mock.AnythingOfType("*dto.SlotFilter")).Return(nil, assert.AnError)
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

			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			handler := slot.NewHandler(&slot.HandlerConfig{
				Service: mockService,
				Logger:  logger,
			})

			url := "/rooms/" + tt.roomID + "/slots/list"
			if tt.queryDate != "" {
				url = url + "?date=" + tt.queryDate
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("roomId", tt.roomID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

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
