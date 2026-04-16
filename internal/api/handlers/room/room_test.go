package room_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mmms914/test-avito-internship/internal/api/handlers/room"
	mocks "github.com/mmms914/test-avito-internship/internal/api/handlers/room/mocks"
	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
	"github.com/mmms914/test-avito-internship/pkg/ptr"
)

func TestHandler_Create(t *testing.T) {
	tests := map[string]struct {
		requestBody    any
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - create room": {
			requestBody: models.CreateRoomRequest{
				Name:        ptr.To("Conference Room A"),
				Description: ptr.To("Spacious room with projector"),
				Capacity:    ptr.To(10),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.RoomInitSpecs")).Return(&domain.Room{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
		},
		"success - create room without description and capacity": {
			requestBody: models.CreateRoomRequest{
				Name: ptr.To("Small Meeting Room"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.RoomInitSpecs")).Return(&domain.Room{}, nil)
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
		"failure - missing name": {
			requestBody: models.CreateRoomRequest{
				Description: ptr.To("test"),
				Capacity:    ptr.To(10),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - internal error": {
			requestBody: models.CreateRoomRequest{
				Name: ptr.To("Conference Room"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.RoomInitSpecs")).Return(nil, assert.AnError)
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
			handler := room.NewHandler(&room.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(setContext(req.Context(), domain.AdminRole))

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
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - list rooms": {
			setupMock: func(m *mocks.Service) {
				m.On("GetAll", mock.Anything).Return([]*domain.Room{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
		},
		"failure - internal error": {
			setupMock: func(m *mocks.Service) {
				m.On("GetAll", mock.Anything).Return(nil, assert.AnError)
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
			handler := room.NewHandler(&room.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)

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

func setContext(ctx context.Context, role domain.Role) context.Context {
	return context.WithValue(context.WithValue(ctx, domain.UserRoleKey, role), domain.UserIDKey, uuid.New())
}
