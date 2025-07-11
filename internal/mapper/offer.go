package mapper

import (
	"financing-aggregator/internal/dto"
	"financing-aggregator/internal/models"
	"github.com/google/uuid"
	"time"
)

const dateFormat = "2006-01-02"

func MapOfferDTOToModel(in dto.OfferDTO, applicationID uuid.UUID) models.Offer {
	firstRepaymentDate, _ := time.Parse(dateFormat, in.FirstRepaymentDate)
	return models.Offer{
		ApplicationID:        applicationID,
		ExternalID:           in.ExternalID,
		Bank:                 in.Bank,
		Status:               in.Status,
		MonthlyPaymentAmount: in.MonthlyPaymentAmount,
		TotalRepaymentAmount: in.TotalRepaymentAmount,
		NumberOfPayments:     in.NumberOfPayments,
		AnnualPercentageRate: in.AnnualPercentageRate,
		FirstRepaymentDate:   firstRepaymentDate,
	}
}

func MapOfferModelToDTO(in models.Offer) dto.OfferDTO {
	return dto.OfferDTO{
		ExternalID:           in.ExternalID,
		Bank:                 in.Bank,
		Status:               in.Status,
		MonthlyPaymentAmount: in.MonthlyPaymentAmount,
		TotalRepaymentAmount: in.TotalRepaymentAmount,
		NumberOfPayments:     in.NumberOfPayments,
		AnnualPercentageRate: in.AnnualPercentageRate,
		FirstRepaymentDate:   in.FirstRepaymentDate.Format(dateFormat),
	}
}
