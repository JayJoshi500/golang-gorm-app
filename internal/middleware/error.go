package middleware

import (
	"log/slog"
	"net/http"

	"github.com/JayJoshi500/golang-gorm-app/pkg/response"
)

// Recoverer catches panics anywhere downstream, logs them with a stack
// trace, and returns a clean 500 JSON envelope instead of letting the
// connection die or leaking a Go stack trace to the client.
func Recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"error", rec,
						"path", r.URL.Path,
						"method", r.Method,
					)
					response.Error(w, http.StatusInternalServerError, "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
