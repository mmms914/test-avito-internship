package auth_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authApi "github.com/avito-internships/test-backend-1-mmms914/internal/api/auth"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/auth"
	mocks "github.com/avito-internships/test-backend-1-mmms914/internal/api/handlers/auth/mocks"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
	"github.com/avito-internships/test-backend-1-mmms914/internal/dto"
	"github.com/avito-internships/test-backend-1-mmms914/internal/errs"
	"github.com/avito-internships/test-backend-1-mmms914/pkg/ptr"
)

func TestHandler_DummyLogin(t *testing.T) {
	tests := map[string]struct {
		name           string
		requestBody    any
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, body []byte)
	}{
		"admin login": {
			requestBody: models.DummyLoginRequest{
				Role: ptr.To("admin"),
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.LoginResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Token)

				claims, err := authApi.ValidateToken(resp.Token, "test-secret-key")
				require.NoError(t, err)
				assert.Equal(t, authApi.DefaultAdminUID, claims.UserID)
				assert.Equal(t, "admin", claims.Role)
			},
		},
		"user login": {
			requestBody: models.DummyLoginRequest{
				Role: ptr.To("user"),
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.LoginResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Token)

				claims, err := authApi.ValidateToken(resp.Token, "test-secret-key")
				require.NoError(t, err)
				assert.Equal(t, authApi.DefaultUserUID, claims.UserID)
				assert.Equal(t, "user", claims.Role)
			},
		},
		"invalid role": {
			requestBody: models.DummyLoginRequest{
				Role: ptr.To("superuser"),
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			validate := validator.New()
			jwtSecret := "test-secret-key"
			expirationHours := 24
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			handler := auth.NewHandler(&auth.HandlerConfig{
				Logger:          logger,
				Service:         nil,
				Validate:        validate,
				JwtSecret:       jwtSecret,
				ExpirationHours: expirationHours,
			})

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.DummyLogin(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	tests := map[string]struct {
		requestBody    any
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, body []byte)
	}{
		"success - admin login": {
			requestBody: models.LoginRequest{
				Email:    ptr.To("admin@example.com"),
				Password: ptr.To("password123"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Login", mock.Anything, &dto.UserCredentials{
					Email:    "admin@example.com",
					Password: "password123",
				}).Return(&domain.User{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.LoginResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
			},
		},
		"success - user login": {
			requestBody: models.LoginRequest{
				Email:    ptr.To("user@example.com"),
				Password: ptr.To("password123"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Login", mock.Anything, &dto.UserCredentials{
					Email:    "user@example.com",
					Password: "password123",
				}).Return(&domain.User{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCode:   "",
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.LoginResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.NotEmpty(t, resp.Token)
			},
		},
		"failure - invalid credentials": {
			requestBody: models.LoginRequest{
				Email:    ptr.To("user@example.com"),
				Password: ptr.To("wrongpassword"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Login", mock.Anything, &dto.UserCredentials{
					Email:    "user@example.com",
					Password: "wrongpassword",
				}).Return(nil, errs.ErrUnauthorized)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   models.UnauthorizedErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - user not found": {
			requestBody: models.LoginRequest{
				Email:    ptr.To("nonexistent@example.com"),
				Password: ptr.To("password123"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Login", mock.Anything, &dto.UserCredentials{
					Email:    "nonexistent@example.com",
					Password: "password123",
				}).Return(nil, errs.ErrUserNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   models.UnauthorizedErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - missing email": {
			requestBody: models.LoginRequest{
				Password: ptr.To("password123"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - missing password": {
			requestBody: models.LoginRequest{
				Email: ptr.To("user@example.com"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - empty email": {
			requestBody: models.LoginRequest{
				Email:    ptr.To(""),
				Password: ptr.To("password123"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - invalid email format": {
			requestBody: models.LoginRequest{
				Email:    ptr.To("invalid-email"),
				Password: ptr.To("password123"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - invalid request body": {
			requestBody:    "invalid json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			validate := validator.New()
			jwtSecret := "test-secret-key"
			expirationHours := 24
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			handler := auth.NewHandler(&auth.HandlerConfig{
				Logger:          logger,
				Service:         mockService,
				Validate:        validate,
				JwtSecret:       jwtSecret,
				ExpirationHours: expirationHours,
			})

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.Login(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}

			if tt.setupMock != nil {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_Register(t *testing.T) {
	tests := map[string]struct {
		requestBody    any
		setupMock      func(*mocks.Service)
		expectedStatus int
		expectedCode   string
		checkResponse  func(t *testing.T, body []byte)
	}{
		"success - register admin": {
			requestBody: models.RegisterRequest{
				Email:    ptr.To("admin@example.com"),
				Password: ptr.To("password123"),
				Role:     ptr.To("admin"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Register", mock.Anything, &dto.UserCreateModel{
					Email:    "admin@example.com",
					Password: "password123",
					Role:     domain.AdminRole,
				}).Return(&domain.User{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.UserResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"success - register user": {
			requestBody: models.RegisterRequest{
				Email:    ptr.To("user@example.com"),
				Password: ptr.To("password123"),
				Role:     ptr.To("user"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Register", mock.Anything, &dto.UserCreateModel{
					Email:    "user@example.com",
					Password: "password123",
					Role:     domain.UserRole,
				}).Return(&domain.User{}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedCode:   "",
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.UserResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - user already exists": {
			requestBody: models.RegisterRequest{
				Email:    ptr.To("existing@example.com"),
				Password: ptr.To("password123"),
				Role:     ptr.To("user"),
			},
			setupMock: func(m *mocks.Service) {
				m.On("Register", mock.Anything, &dto.UserCreateModel{
					Email:    "existing@example.com",
					Password: "password123",
					Role:     domain.UserRole,
				}).Return(nil, errs.ErrUserAlreadyExists)
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - invalid role": {
			requestBody: models.RegisterRequest{
				Email:    ptr.To("user@example.com"),
				Password: ptr.To("password123"),
				Role:     ptr.To("superuser"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
			},
		},
		"failure - missing email": {
			requestBody: models.RegisterRequest{
				Password: ptr.To("password123"),
				Role:     ptr.To("user"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, models.InvalidRequestErrorCode, resp.Error.Code)
			},
		},
		"failure - missing password": {
			requestBody: models.RegisterRequest{
				Email: ptr.To("user@example.com"),
				Role:  ptr.To("user"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, models.InvalidRequestErrorCode, resp.Error.Code)
			},
		},
		"failure - missing role": {
			requestBody: models.RegisterRequest{
				Email:    ptr.To("user@example.com"),
				Password: ptr.To("password123"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, models.InvalidRequestErrorCode, resp.Error.Code)
			},
		},
		"failure - invalid email format": {
			requestBody: models.RegisterRequest{
				Email:    ptr.To("invalid-email"),
				Password: ptr.To("password123"),
				Role:     ptr.To("user"),
			},
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, models.InvalidRequestErrorCode, resp.Error.Code)
			},
		},
		"failure - invalid request body": {
			requestBody:    "invalid json",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.InvalidRequestErrorCode,
			checkResponse: func(t *testing.T, body []byte) {
				var resp models.ErrorResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, models.InvalidRequestErrorCode, resp.Error.Code)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockService := mocks.NewService(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			validate := validator.New()
			jwtSecret := "test-secret-key"
			expirationHours := 24
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			handler := auth.NewHandler(&auth.HandlerConfig{
				Logger:          logger,
				Service:         mockService,
				Validate:        validate,
				JwtSecret:       jwtSecret,
				ExpirationHours: expirationHours,
			})

			var reqBody []byte
			if tt.requestBody != nil {
				reqBody, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.Register(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			if tt.expectedCode != "" {
				var errResp models.ErrorResponse
				err := json.Unmarshal(rr.Body.Bytes(), &errResp)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCode, errResp.Error.Code)
			}
		})
	}
}
