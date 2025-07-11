package repositories

import (
	"context"
	"financing-aggregator/internal/models"
	"gorm.io/gorm"
)

type OfferRepository interface {
	Create(ctx context.Context, offer *models.Offer) error
	List(ctx context.Context, filter OfferListFilter) ([]models.Offer, error)
	Update(ctx context.Context, id string, offer models.Offer) error
}

type OfferListFilter struct {
	Status string
}

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{db: db}
}

func (r *offerRepository) Create(ctx context.Context, offer *models.Offer) error {
	return r.db.WithContext(ctx).Create(offer).Error
}

func (r *offerRepository) List(ctx context.Context, filter OfferListFilter) ([]models.Offer, error) {
	var offers []models.Offer
	err := r.db.WithContext(ctx).Where("status = ?", filter.Status).Find(&offers).Error
	return offers, err
}

func (r *offerRepository) Update(ctx context.Context, id string, offer models.Offer) error {
	return r.db.WithContext(ctx).Model(&models.Offer{}).Where("id = ?", id).Updates(offer).Error
}
