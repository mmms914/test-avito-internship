package schedule_test

import (
	"bytes"
	"context"
	"encoding/json"
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

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/schedule"
	mocks "github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/schedule/mocks"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func TestHandler_Create(t *testing.T) {
	tests := map[string]struct {
		roomID         string
		requestBody    interface{}
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
	}{
		"success - create schedule": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1, 2, 3, 4, 5}),
				StartTime:  ptr.To("09:00"),
				EndTime:    ptr.To("18:00"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.ScheduleInitSpecs")).Return(&domain.Schedule{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
		},
		"success - create schedule with single day": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("10:00"),
				EndTime:    ptr.To("14:00"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.ScheduleInitSpecs")).Return(&domain.Schedule{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
		},
		"failure - invalid roomId": {
			roomID:         "invalid-uuid",
			requestBody:    models.CreateScheduleRequest{},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid request body": {
			roomID:         uuid.New().String(),
			requestBody:    "invalid json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - missing daysOfWeek": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				StartTime: ptr.To("09:00"),
				EndTime:   ptr.To("18:00"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid day of week": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{8}),
				StartTime:  ptr.To("09:00"),
				EndTime:    ptr.To("18:00"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - missing startTime": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				EndTime:    ptr.To("18:00"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - missing endTime": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("09:00"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid startTime format": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("119:00"),
				EndTime:    ptr.To("18:00"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - invalid endTime format": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("9:00"),
				EndTime:    ptr.To("116:00"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
		},
		"failure - room not found": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("09:00"),
				EndTime:    ptr.To("18:00"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.ScheduleInitSpecs")).Return(nil, errs.ErrRoomNotExists)
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   models.RoomNotFoundErrorCode,
		},
		"failure - schedule already exists": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("09:00"),
				EndTime:    ptr.To("18:00"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.ScheduleInitSpecs")).Return(nil, errs.ErrScheduleExists)
			},
			expectedStatus: http.StatusConflict,
			expectedCode:   models.ScheduleExistsErrorCode,
		},
		"failure - internal error": {
			roomID: uuid.New().String(),
			requestBody: models.CreateScheduleRequest{
				DaysOfWeek: ptr.To([]int{1}),
				StartTime:  ptr.To("09:00"),
				EndTime:    ptr.To("18:00"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*domain.ScheduleInitSpecs")).Return(nil, assert.AnError)
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
			handler := schedule.NewHandler(&schedule.HandlerConfig{
				Service:  mockService,
				Validate: validate,
				Logger:   logger,
			})

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/rooms/"+tt.roomID+"/schedule/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(setContext(req.Context(), domain.AdminRole))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("roomId", tt.roomID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

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

func setContext(ctx context.Context, role domain.Role) context.Context {
	return context.WithValue(context.WithValue(ctx, domain.UserRoleKey, role), domain.UserIDKey, uuid.New())
}
