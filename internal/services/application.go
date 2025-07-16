package services

import (
	"context"
	"financing-aggregator/internal/banks"
	"financing-aggregator/internal/controllers/ws"
	"financing-aggregator/internal/dto"
	"financing-aggregator/internal/mapper"
	"financing-aggregator/internal/models"
	"financing-aggregator/internal/repositories"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ApplicationService interface {
	SubmitApplication(ctx context.Context, app dto.ApplicationDTO) (dto.ApplicationDTO, error)
	GetApplication(ctx context.Context, id string) (dto.ApplicationDTO, error)
	UpdateApplicationStatuses(ctx context.Context)
}

type applicationService struct {
	logger          *zap.Logger
	banks           map[string]banks.Bank
	applicationRepo repositories.ApplicationRepository
	offerRepo       repositories.OfferRepository
	wsHandler       ws.WebSocketHandler
}

func NewApplicationService(
	logger *zap.Logger,
	allBanks []banks.Bank,
	applicationRepo repositories.ApplicationRepository,
	offerRepo repositories.OfferRepository,
	wsHandler ws.WebSocketHandler,
) ApplicationService {
	bankMap := lo.SliceToMap(allBanks, func(b banks.Bank) (string, banks.Bank) {
		return b.Name(), b
	})

	return &applicationService{
		logger:          logger,
		banks:           bankMap,
		applicationRepo: applicationRepo,
		offerRepo:       offerRepo,
		wsHandler:       wsHandler,
	}
}

func (s *applicationService) SubmitApplication(ctx context.Context, app dto.ApplicationDTO) (dto.ApplicationDTO, error) {
	appModel := mapper.MapApplicationDTOToModel(app)
	if err := s.applicationRepo.Create(ctx, &appModel); err != nil {
		s.logger.Error("failed to create application", zap.Error(err))
		return dto.ApplicationDTO{}, fmt.Errorf("failed to create application: %v", err)
	}

	for _, bank := range s.banks {
		go func(b banks.Bank, a dto.ApplicationDTO) {
			ctx := context.Background()

			offer, err := b.SubmitApplication(ctx, a)
			if err != nil {
				s.logger.Error("failed to submit application", zap.Error(err), zap.String("bank", b.Name()), zap.String("id", a.ID))
				return
			}

			offerModel := mapper.MapOfferDTOToModel(offer, appModel.ID)
			if err := s.offerRepo.Create(ctx, &offerModel); err != nil {
				s.logger.Error("failed to create offer", zap.Error(err), zap.String("bank", b.Name()), zap.String("id", a.ID))
				return
			}
		}(bank, app)
	}

	app.ID = appModel.ID.String()
	return app, nil
}

func (s *applicationService) GetApplication(ctx context.Context, id string) (dto.ApplicationDTO, error) {
	application, err := s.applicationRepo.GetWithProcessedOffers(ctx, id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ApplicationDTO{}, err
		}
		return dto.ApplicationDTO{}, fmt.Errorf("failed to get application: %v", err)
	}

	return mapper.MapApplicationModelToDTO(application), nil
}

func (s *applicationService) UpdateApplicationStatuses(ctx context.Context) {
	offers, err := s.offerRepo.List(ctx, repositories.OfferListFilter{
		Status: models.OfferStatusDraft,
	})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("failed to list draft offers", zap.Error(err))
		}
		return
	}

	for _, offer := range offers {
		bank, ok := s.banks[offer.Bank]
		if !ok {
			s.logger.Error("offer belongs to unknown bank", zap.String("bank", offer.Bank))
			continue
		}

		bankOffer, err := bank.GetApplication(ctx, offer.ExternalID)
		if err != nil {
			s.logger.Error("failed to get application from bank", zap.Error(err), zap.String("bank", offer.Bank), zap.String("id", offer.ExternalID))
			continue
		}

		if bankOffer.Status == offer.Status {
			continue
		}

		model := mapper.MapOfferDTOToModel(bankOffer, offer.ApplicationID)
		if bankOffer.NumberOfPayments == 0 {
			model.Status = models.OfferStatusDeclined
		}

		err = s.offerRepo.Update(ctx, offer.ID.String(), model)
		if err != nil {
			s.logger.Error("failed to get update offer", zap.Error(err), zap.String("bank", offer.Bank), zap.String("id", offer.ID.String()))
			continue
		}

		if bankOffer.Status == models.OfferStatusDeclined {
			continue
		}

		s.wsHandler.BroadcastNewOffer(offer.ApplicationID.String(), mapper.MapOfferDTOToResponse(bankOffer))
	}
}
