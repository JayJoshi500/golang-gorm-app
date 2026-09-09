package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/JayJoshi500/golang-gorm-app/internal/config"
	"github.com/JayJoshi500/golang-gorm-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Connect opens a GORM/Postgres connection using values from cfg, tunes
// the underlying connection pool, and returns the *gorm.DB handle.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	slog.Info("connected to database", "host", cfg.DBHost, "name", cfg.DBName)
	return db, nil
}

// AutoMigrate creates/updates tables for all known models. Fine for
// development; consider a real migration tool (goose, atlas, etc.) for
// production schema changes.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}
