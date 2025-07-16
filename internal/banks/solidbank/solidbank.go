package solidbank

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"financing-aggregator/internal/dto"
	"fmt"
	"io"
	"net/http"

	"financing-aggregator/internal/banks"
)

type SolidBank struct {
	BaseURL string
	Client  *http.Client
}

func NewSolidBank(baseURL string) banks.Bank {
	return &SolidBank{
		BaseURL: baseURL,
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (b *SolidBank) Name() string {
	return "solidbank"
}

func (b *SolidBank) SubmitApplication(ctx context.Context, data dto.ApplicationDTO) (dto.OfferDTO, error) {
	reqData := ApplicationRequest{
		Phone:           data.Phone,
		Email:           data.Email,
		MonthlyIncome:   data.MonthlyIncome,
		MonthlyExpenses: data.MonthlyExpenses,
		MaritalStatus:   data.MaritalStatus,
		AgreeToBeScored: data.AgreeToBeScored,
		Amount:          data.Amount,
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return dto.OfferDTO{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.BaseURL+"/applications", bytes.NewBuffer(reqBody))
	if err != nil {
		return dto.OfferDTO{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.Client.Do(req)
	if err != nil {
		return dto.OfferDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return dto.OfferDTO{}, fmt.Errorf("solidbank: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var response ApplicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.OfferDTO{}, err
	}
	return dto.OfferDTO{
		ExternalID:           response.ID,
		Status:               response.Status,
		Bank:                 b.Name(),
		MonthlyPaymentAmount: response.Offer.MonthlyPaymentAmount,
		TotalRepaymentAmount: response.Offer.TotalRepaymentAmount,
		NumberOfPayments:     response.Offer.NumberOfPayments,
		AnnualPercentageRate: response.Offer.AnnualPercentageRate,
		FirstRepaymentDate:   response.Offer.FirstRepaymentDate,
	}, nil
}

func (b *SolidBank) GetApplication(ctx context.Context, id string) (dto.OfferDTO, error) {
	url := fmt.Sprintf("%s/applications/%s", b.BaseURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dto.OfferDTO{}, err
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return dto.OfferDTO{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return dto.OfferDTO{}, fmt.Errorf("solidbank: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var response ApplicationResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.OfferDTO{}, err
	}
	return dto.OfferDTO{
		ExternalID:           response.ID,
		Status:               response.Status,
		Bank:                 b.Name(),
		MonthlyPaymentAmount: response.Offer.MonthlyPaymentAmount,
		TotalRepaymentAmount: response.Offer.TotalRepaymentAmount,
		NumberOfPayments:     response.Offer.NumberOfPayments,
		AnnualPercentageRate: response.Offer.AnnualPercentageRate,
		FirstRepaymentDate:   response.Offer.FirstRepaymentDate,
	}, nil
}
