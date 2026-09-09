package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/JayJoshi500/golang-gorm-app/internal/services"
	"github.com/JayJoshi500/golang-gorm-app/pkg/response"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

// RequireAuth validates the bearer JWT on protected routes and rejects
// tokens that are expired, malformed, or have been logged out. Not wired
// into the auth routes themselves (register/login/logout are public),
// but ready to mount on any future protected route group.
func RequireAuth(authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}
			token := strings.TrimSpace(parts[1])

			if authService.IsBlacklisted(token) {
				response.Error(w, http.StatusUnauthorized, "token has been revoked")
				return
			}

			claims, err := authService.ParseToken(token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
