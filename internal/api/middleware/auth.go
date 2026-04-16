package appmiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/mmms914/test-avito-internship/internal/api/auth"
	"github.com/mmms914/test-avito-internship/internal/api/converter"
	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer"
)

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(authorizationHeader)
			if authHeader == "" {
				writeUnauthorizedError(w, "authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != bearerPrefix {
				writeUnauthorizedError(w, "invalid authorization header format, expected Bearer token")
				return
			}

			token := parts[1]

			claims, err := auth.ValidateToken(token, jwtSecret)
			if err != nil {
				writeUnauthorizedError(w, "invalid or expired token")
				return
			}

			userRole, err := converter.UserRoleFromString(claims.Role)
			if err != nil {
				writeUnauthorizedError(w, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), domain.UserIDKey, uuid.MustParse(claims.UserID))
			ctx = context.WithValue(ctx, domain.UserRoleKey, userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorizedError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    models.UnauthorizedErrorCode,
			"message": message,
		},
	})
}
