package services

import (
	"context"
	"financing-aggregator/internal/banks"
	"financing-aggregator/internal/dto"
	"financing-aggregator/internal/exchange"
	mock_banks "financing-aggregator/internal/mocks/banks"
	mock_ws "financing-aggregator/internal/mocks/controllers/ws"
	mock_repositories "financing-aggregator/internal/mocks/repositories"
	"financing-aggregator/internal/models"
	"financing-aggregator/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
	"time"
)

type applicationServiceTestSuite struct {
	suite.Suite
	ctrl   *gomock.Controller
	logger *zap.Logger

	applicationRepository *mock_repositories.MockApplicationRepository
	offerRepository       *mock_repositories.MockOfferRepository
	banks                 []banks.Bank
	bank1                 *mock_banks.MockBank
	bank2                 *mock_banks.MockBank
	wsHandler             *mock_ws.MockWebSocketHandler

	service *applicationService
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(applicationServiceTestSuite))
}

func (s *applicationServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.logger = zap.NewNop()

	s.applicationRepository = mock_repositories.NewMockApplicationRepository(s.ctrl)
	s.offerRepository = mock_repositories.NewMockOfferRepository(s.ctrl)
	s.bank1 = mock_banks.NewMockBank(s.ctrl)
	s.bank2 = mock_banks.NewMockBank(s.ctrl)
	s.banks = []banks.Bank{s.bank1, s.bank2}
	s.wsHandler = mock_ws.NewMockWebSocketHandler(s.ctrl)

	s.bank1.EXPECT().Name().Return("bank1").AnyTimes()
	s.bank2.EXPECT().Name().Return("bank2").AnyTimes()

	s.service = NewApplicationService(s.logger, s.banks, s.applicationRepository, s.offerRepository, s.wsHandler).(*applicationService)
}

func (s *applicationServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *applicationServiceTestSuite) Test_SubmitApplication() {
	applicationDTO := getTestApplicationDTO()
	offerDTO1 := getTestOfferDTO("bank1")
	offerDTO2 := getTestOfferDTO("bank2")

	s.Run("application submitted", func() {
		s.applicationRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		s.bank1.EXPECT().SubmitApplication(gomock.Any(), applicationDTO).Return(offerDTO1, nil)
		s.bank2.EXPECT().SubmitApplication(gomock.Any(), applicationDTO).Return(offerDTO2, nil)
		s.offerRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Times(2).Return(nil)

		_, err := s.service.SubmitApplication(context.Background(), applicationDTO)
		time.Sleep(1 * time.Second)
		s.NoError(err)
	})

	s.Run("error occurs while saving application", func() {
		s.applicationRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

		_, err := s.service.SubmitApplication(context.Background(), applicationDTO)
		time.Sleep(1 * time.Second)
		s.Error(err)
		s.Contains(err.Error(), "db error")
	})

	s.Run("error occurs while submitting application to bank", func() {
		s.applicationRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		s.bank1.EXPECT().SubmitApplication(gomock.Any(), applicationDTO).Return(offerDTO1, errors.New("bank error"))
		s.bank2.EXPECT().SubmitApplication(gomock.Any(), applicationDTO).Return(offerDTO2, nil)
		s.offerRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		_, err := s.service.SubmitApplication(context.Background(), applicationDTO)
		time.Sleep(1 * time.Second)
		s.NoError(err)
	})

	s.Run("error occurs while saving offer", func() {
		s.applicationRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		s.bank1.EXPECT().SubmitApplication(gomock.Any(), applicationDTO).Return(offerDTO1, nil)
		s.bank2.EXPECT().SubmitApplication(gomock.Any(), applicationDTO).Return(offerDTO2, nil)
		s.offerRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("offer save error")).AnyTimes()

		_, err := s.service.SubmitApplication(context.Background(), applicationDTO)
		time.Sleep(1 * time.Second)
		s.NoError(err)
	})
}

func (s *applicationServiceTestSuite) Test_GetApplication() {
	applicationDTO := getTestApplicationDTO()
	applicationModel := getTestApplicationModel()

	s.Run("application found", func() {
		s.applicationRepository.EXPECT().GetWithProcessedOffers(gomock.Any(), applicationDTO.ID).Return(applicationModel, nil)
		actual, err := s.service.GetApplication(context.Background(), applicationDTO.ID)
		s.NoError(err)
		s.NotNil(actual)
		s.Equal(applicationDTO, actual)
	})

	s.Run("error occurs because application not found", func() {
		s.applicationRepository.EXPECT().GetWithProcessedOffers(gomock.Any(), applicationDTO.ID).Return(models.Application{}, gorm.ErrRecordNotFound)
		actual, err := s.service.GetApplication(context.Background(), applicationDTO.ID)
		s.Error(err)
		s.Equal(dto.ApplicationDTO{}, actual)
		s.Contains(err.Error(), "not found")
	})

	s.Run("error occurs while getting application", func() {
		s.applicationRepository.EXPECT().GetWithProcessedOffers(gomock.Any(), applicationDTO.ID).Return(models.Application{}, errors.New("db error"))
		actual, err := s.service.GetApplication(context.Background(), applicationDTO.ID)
		s.Error(err)
		s.Equal(dto.ApplicationDTO{}, actual)
		s.Contains(err.Error(), "db error")
	})
}

func (s *applicationServiceTestSuite) Test_UpdateApplicationStatuses() {
	offerModel := getTestOfferModel("bank1")
	offerModels := []models.Offer{offerModel}

	s.Run("application statuses updated", func() {
		bankOffer := getTestOfferDTO("bank1")
		bankOffer.Status = "PROCESSED"
		updatedOfferModel := getTestOfferModel("bank1")
		updatedOfferModel.Status = "PROCESSED"
		offerResponse := getTestOfferResponse()

		s.offerRepository.EXPECT().List(gomock.Any(), repositories.OfferListFilter{Status: "DRAFT"}).Return(offerModels, nil)
		s.bank1.EXPECT().GetApplication(gomock.Any(), offerModel.ExternalID).Return(bankOffer, nil)
		s.offerRepository.EXPECT().Update(gomock.Any(), offerModel.ID.String(), updatedOfferModel).Return(nil)
		s.wsHandler.EXPECT().BroadcastNewOffer(offerModel.ApplicationID.String(), offerResponse)

		s.service.UpdateApplicationStatuses(context.Background())
	})

	s.Run("error occurs while listing offers", func() {
		s.offerRepository.EXPECT().List(gomock.Any(), repositories.OfferListFilter{Status: "DRAFT"}).Return(nil, errors.New("db error"))

		s.service.UpdateApplicationStatuses(context.Background())
	})

	s.Run("error occurs while getting application from bank", func() {
		s.offerRepository.EXPECT().List(gomock.Any(), repositories.OfferListFilter{Status: "DRAFT"}).Return(offerModels, nil)
		s.bank1.EXPECT().GetApplication(gomock.Any(), offerModel.ExternalID).Return(dto.OfferDTO{}, errors.New("bank error"))

		s.service.UpdateApplicationStatuses(context.Background())
	})

	s.Run("bank application status is the same as offer status", func() {
		s.offerRepository.EXPECT().List(gomock.Any(), repositories.OfferListFilter{Status: "DRAFT"}).Return(offerModels, nil)
		s.bank1.EXPECT().GetApplication(gomock.Any(), offerModel.ExternalID).Return(getTestOfferDTO("bank1"), nil)

		s.service.UpdateApplicationStatuses(context.Background())
	})

	s.Run("error occurs while updating offer", func() {
		bankOffer := getTestOfferDTO("bank1")
		bankOffer.Status = "PROCESSED"
		updatedOfferModel := getTestOfferModel("bank1")
		updatedOfferModel.Status = "PROCESSED"

		s.offerRepository.EXPECT().List(gomock.Any(), repositories.OfferListFilter{Status: "DRAFT"}).Return(offerModels, nil)
		s.bank1.EXPECT().GetApplication(gomock.Any(), offerModel.ExternalID).Return(bankOffer, nil)
		s.offerRepository.EXPECT().Update(gomock.Any(), offerModel.ID.String(), updatedOfferModel).Return(errors.New("update error"))

		s.service.UpdateApplicationStatuses(context.Background())
	})
}

func getTestApplicationDTO() dto.ApplicationDTO {
	return dto.ApplicationDTO{
		ID:              uuid.UUID{}.String(),
		Phone:           "+37122334455",
		Email:           "anakin@skywalker.com",
		Amount:          100,
		MonthlyIncome:   1000,
		MonthlyExpenses: 100,
		Offers:          []dto.OfferDTO{},
	}
}

func getTestApplicationModel() models.Application {
	return models.Application{
		ID:              uuid.UUID{},
		Phone:           "+37122334455",
		Email:           "anakin@skywalker.com",
		Amount:          100,
		MonthlyIncome:   1000,
		MonthlyExpenses: 100,
		Offers:          []models.Offer{},
	}
}

func getTestOfferDTO(bank string) dto.OfferDTO {
	return dto.OfferDTO{
		ExternalID:           bank + "-offer-1",
		Status:               "DRAFT",
		Bank:                 bank,
		MonthlyPaymentAmount: 50,
		TotalRepaymentAmount: 150,
		NumberOfPayments:     3,
		AnnualPercentageRate: 10.0,
		FirstRepaymentDate:   "2025-01-01",
	}
}

func getTestOfferModel(bank string) models.Offer {
	firstRepaymentDate, _ := time.Parse("2006-01-02", "2025-01-01")
	return models.Offer{
		ID:                   uuid.UUID{},
		ExternalID:           bank + "-offer-1",
		Status:               "DRAFT",
		Bank:                 bank,
		MonthlyPaymentAmount: 50,
		TotalRepaymentAmount: 150,
		NumberOfPayments:     3,
		AnnualPercentageRate: 10.0,
		FirstRepaymentDate:   firstRepaymentDate,
		ApplicationID:        uuid.UUID{},
	}
}

func getTestOfferResponse() exchange.OfferResponse {
	return exchange.OfferResponse{
		MonthlyPaymentAmount: 50,
		TotalRepaymentAmount: 150,
		NumberOfPayments:     3,
		AnnualPercentageRate: 10.0,
		FirstRepaymentDate:   "2025-01-01",
	}
}
