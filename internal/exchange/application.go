package exchange

type ApplicationRequest struct {
	Phone                    string  `json:"phone" validate:"e164,startswith=+371,len=12"`
	Email                    string  `json:"email" validate:"email"`
	MonthlyIncome            float64 `json:"monthlyIncome" validate:"gte=0"`
	MonthlyExpenses          float64 `json:"monthlyExpenses" validate:"gte=0"`
	MonthlyCreditLiabilities float64 `json:"monthlyCreditLiabilities" validate:"gte=0"`
	MaritalStatus            string  `json:"maritalStatus" validate:"oneof=SINGLE MARRIED DIVORCED COHABITING"`
	Dependents               int     `json:"dependents" validate:"gte=0"`
	AgreeToDataSharing       bool    `json:"agreeToDataSharing"`
	AgreeToBeScored          bool    `json:"agreeToBeScored"`
	Amount                   float64 `json:"amount" validate:"gte=0"`
}

type ApplicationResponse struct {
	ID                       string          `json:"id"`
	Phone                    string          `json:"phone"`
	Email                    string          `json:"email"`
	MonthlyIncome            float64         `json:"monthlyIncome"`
	MonthlyExpenses          float64         `json:"monthlyExpenses"`
	MonthlyCreditLiabilities float64         `json:"monthlyCreditLiabilities"`
	MaritalStatus            string          `json:"maritalStatus"`
	Dependents               int             `json:"dependents"`
	AgreeToDataSharing       bool            `json:"agreeToDataSharing"`
	AgreeToBeScored          bool            `json:"agreeToBeScored"`
	Amount                   float64         `json:"amount"`
	Offers                   []OfferResponse `json:"offers,omitempty"`
}

type OfferResponse struct {
	MonthlyPaymentAmount float64 `json:"monthlyPaymentAmount"`
	TotalRepaymentAmount float64 `json:"totalRepaymentAmount"`
	NumberOfPayments     int     `json:"numberOfPayments"`
	AnnualPercentageRate float64 `json:"annualPercentageRate"`
	FirstRepaymentDate   string  `json:"firstRepaymentDate"`
}
