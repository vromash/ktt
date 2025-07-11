package fastbank

type (
	ApplicationRequest struct {
		PhoneNumber              string  `json:"phoneNumber"`
		Email                    string  `json:"email"`
		MonthlyIncomeAmount      float64 `json:"monthlyIncomeAmount"`
		MonthlyCreditLiabilities float64 `json:"monthlyCreditLiabilities"`
		Dependents               int     `json:"dependents"`
		AgreeToDataSharing       bool    `json:"agreeToDataSharing"`
		Amount                   float64 `json:"amount"`
	}

	ApplicationResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Offer  Offer  `json:"offer"`
	}

	Offer struct {
		MonthlyPaymentAmount float64 `json:"monthlyPaymentAmount"`
		TotalRepaymentAmount float64 `json:"totalRepaymentAmount"`
		NumberOfPayments     int     `json:"numberOfPayments"`
		AnnualPercentageRate float64 `json:"annualPercentageRate"`
		FirstRepaymentDate   string  `json:"firstRepaymentDate"`
	}
)
