package services

import (
	"context"

	"github.com/JayJoshi500/golang-gorm-app/internal/models"
	"gorm.io/gorm"
)

type RestoService interface {
	GetAvailableSlots(ctx context.Context) ([]*models.RestoTables, error)
	GetTimingSlots(ctx context.Context) ([]*models.RestoTimingSlots, error)
}

type RestoSlots struct {
	db *gorm.DB
}

// NewAuthService wires a GORM handle plus JWT settings into an AuthService.
func NewRestoService(db *gorm.DB) RestoService {
	return &RestoSlots{
		db: db,
	}
}

func (r *RestoSlots) GetAvailableSlots(ctx context.Context) ([]*models.RestoTables, error) {
	var restoSlots []*models.RestoTables

	if err := r.db.
		WithContext(ctx).
		Where("booked_val < capacity").
		Find(&restoSlots).Error; err != nil {
		return nil, err
	}

	return restoSlots, nil
}

func (r *RestoSlots) GetTimingSlots(ctx context.Context) ([]*models.RestoTimingSlots, error) {
	var restoTimingSlots []*models.RestoTimingSlots

	if err := r.db.WithContext(ctx).Find(&restoTimingSlots).Error; err != nil {
		return nil, err
	}

	return restoTimingSlots, nil
}
