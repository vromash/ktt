package dto

type (
	ApplicationDTO struct {
		ID                       string
		Phone                    string
		Email                    string
		Amount                   float64
		MonthlyIncome            float64
		MonthlyExpenses          float64
		MonthlyCreditLiabilities float64
		MaritalStatus            string
		Dependents               int
		AgreeToDataSharing       bool
		AgreeToBeScored          bool
		Offers                   []OfferDTO
	}

	OfferDTO struct {
		ExternalID           string
		Status               string
		Bank                 string
		MonthlyPaymentAmount float64
		TotalRepaymentAmount float64
		NumberOfPayments     int
		AnnualPercentageRate float64
		FirstRepaymentDate   string
	}
)
