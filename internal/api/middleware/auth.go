package appmiddleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/avito-internships/test-backend-1-mmms914/internal/api/auth"
	"github.com/avito-internships/test-backend-1-mmms914/internal/api/models"
	"github.com/avito-internships/test-backend-1-mmms914/internal/domain"
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
				writeError(w, models.UnauthorizedErrorCode, "authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != bearerPrefix {
				writeError(w, models.UnauthorizedErrorCode, "invalid authorization header format, expected Bearer token", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			claims, err := auth.ValidateToken(token, jwtSecret)
			if err != nil {
				writeError(w, models.UnauthorizedErrorCode, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), domain.UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, domain.UserRoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeError(w http.ResponseWriter, code, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}
