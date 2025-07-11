package repositories

import (
	"context"
	"financing-aggregator/internal/models"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(ctx context.Context, app *models.Application) error
	GetWithProcessedOffers(ctx context.Context, id string) (models.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *models.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *applicationRepository) GetWithProcessedOffers(ctx context.Context, id string) (models.Application, error) {
	var app models.Application
	err := r.db.WithContext(ctx).Preload("Offers", "status = ?", "PROCESSED").First(&app, "id = ?", id).Error
	if err != nil {
		return models.Application{}, err
	}
	return app, nil
}
