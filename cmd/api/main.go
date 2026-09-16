package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JayJoshi500/golang-gorm-app/internal/config"
	"github.com/JayJoshi500/golang-gorm-app/internal/database"
	"github.com/JayJoshi500/golang-gorm-app/internal/handlers"
	"github.com/JayJoshi500/golang-gorm-app/internal/router"
	"github.com/JayJoshi500/golang-gorm-app/internal/services"
	applogger "github.com/JayJoshi500/golang-gorm-app/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := applogger.New(cfg.LogLevel)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	authService := services.NewAuthService(db, cfg.JWTSecret, cfg.JWTExpiryTime)
	authHandler := handlers.NewAuthHandler(authService, log)

	restoService := services.NewRestoService(db)
	restroHandler := handlers.NewRestoHandler(restoService, log)

	r := router.New(router.Handlers{Auth: authHandler, Resto: restroHandler}, log)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run the server in a goroutine so we can listen for shutdown signals
	// on the main goroutine below.
	go func() {
		log.Info("server starting", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	log.Info("server stopped cleanly")
}
