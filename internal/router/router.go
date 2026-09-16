package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/JayJoshi500/golang-gorm-app/internal/handlers"
	appmiddleware "github.com/JayJoshi500/golang-gorm-app/internal/middleware"
	"github.com/JayJoshi500/golang-gorm-app/pkg/response"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Handlers struct {
	Auth  *handlers.AuthHandler
	Resto *handlers.RestoHandler
}

// New builds the full chi.Mux: global middleware, health check, and the
// versioned API route tree.
func New(h Handlers, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware, outermost first.
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(appmiddleware.RequestLogger(logger))
	r.Use(appmiddleware.Recoverer(logger))
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(api chi.Router) {
		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/register", h.Auth.Register)
			auth.Post("/login", h.Auth.Login)
			auth.Post("/logout", h.Auth.Logout)
		})

		api.Route("/resto", func(auth chi.Router) {
			auth.Get("/availableSlots", h.Resto.AvailableSlots)
			auth.Get("/timingSlots", h.Resto.TimingSlots)
		})

	})

	return r
}
