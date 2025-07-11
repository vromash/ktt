package mapper

import (
	"financing-aggregator/internal/dto"
	"financing-aggregator/internal/exchange"
	"financing-aggregator/internal/models"
)

func MapApplicationRequestToDTO(in exchange.ApplicationRequest) dto.ApplicationDTO {
	return dto.ApplicationDTO{
		Phone:                    in.Phone,
		Email:                    in.Email,
		MonthlyIncome:            in.MonthlyIncome,
		MonthlyExpenses:          in.MonthlyExpenses,
		MonthlyCreditLiabilities: in.MonthlyCreditLiabilities,
		MaritalStatus:            in.MaritalStatus,
		Dependents:               in.Dependents,
		AgreeToDataSharing:       in.AgreeToDataSharing,
		AgreeToBeScored:          in.AgreeToBeScored,
		Amount:                   in.Amount,
	}
}

func MapApplicationDTOToResponse(in dto.ApplicationDTO) exchange.ApplicationResponse {
	offers := make([]exchange.OfferResponse, 0, len(in.Offers))
	for _, o := range in.Offers {
		offers = append(offers, MapOfferDTOToResponse(o))
	}

	return exchange.ApplicationResponse{
		ID:                 in.ID,
		Phone:              in.Phone,
		Email:              in.Email,
		MonthlyIncome:      in.MonthlyIncome,
		MonthlyExpenses:    in.MonthlyExpenses,
		MaritalStatus:      in.MaritalStatus,
		Dependents:         in.Dependents,
		AgreeToDataSharing: in.AgreeToDataSharing,
		AgreeToBeScored:    in.AgreeToBeScored,
		Amount:             in.Amount,
		Offers:             offers,
	}
}

func MapApplicationDTOToModel(in dto.ApplicationDTO) models.Application {
	return models.Application{
		Phone:                    in.Phone,
		Email:                    in.Email,
		MonthlyIncome:            in.MonthlyIncome,
		MonthlyExpenses:          in.MonthlyExpenses,
		MonthlyCreditLiabilities: in.MonthlyCreditLiabilities,
		MaritalStatus:            in.MaritalStatus,
		Dependents:               in.Dependents,
		AgreeToDataSharing:       in.AgreeToDataSharing,
		AgreeToBeScored:          in.AgreeToBeScored,
		Amount:                   in.Amount,
	}
}

func MapApplicationModelToDTO(in models.Application) dto.ApplicationDTO {
	offers := make([]dto.OfferDTO, 0, len(in.Offers))
	for _, o := range in.Offers {
		offers = append(offers, MapOfferModelToDTO(o))
	}

	return dto.ApplicationDTO{
		ID:                       in.ID.String(),
		Phone:                    in.Phone,
		Email:                    in.Email,
		MonthlyIncome:            in.MonthlyIncome,
		MonthlyExpenses:          in.MonthlyExpenses,
		MonthlyCreditLiabilities: in.MonthlyCreditLiabilities,
		MaritalStatus:            in.MaritalStatus,
		Dependents:               in.Dependents,
		AgreeToDataSharing:       in.AgreeToDataSharing,
		AgreeToBeScored:          in.AgreeToBeScored,
		Amount:                   in.Amount,
		Offers:                   offers,
	}
}

func MapOfferDTOToResponse(in dto.OfferDTO) exchange.OfferResponse {
	return exchange.OfferResponse{
		MonthlyPaymentAmount: in.MonthlyPaymentAmount,
		TotalRepaymentAmount: in.TotalRepaymentAmount,
		NumberOfPayments:     in.NumberOfPayments,
		AnnualPercentageRate: in.AnnualPercentageRate,
		FirstRepaymentDate:   in.FirstRepaymentDate,
	}
}
