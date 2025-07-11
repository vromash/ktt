package fastbank

import (
	"bytes"
	"context"
	"encoding/json"
	"financing-aggregator/internal/dto"
	"fmt"
	"io"
	"net/http"

	"financing-aggregator/internal/banks"
)

type FastBank struct {
	BaseURL string
	Client  *http.Client
}

func NewFastBank(baseURL string) banks.Bank {
	return &FastBank{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (b *FastBank) Name() string {
	return "fastbank"
}

func (b *FastBank) SubmitApplication(ctx context.Context, data dto.ApplicationDTO) (dto.OfferDTO, error) {
	reqData := ApplicationRequest{
		PhoneNumber:              data.Phone,
		Email:                    data.Email,
		MonthlyIncomeAmount:      data.MonthlyIncome,
		MonthlyCreditLiabilities: data.MonthlyCreditLiabilities,
		Dependents:               data.Dependents,
		AgreeToDataSharing:       data.AgreeToDataSharing,
		Amount:                   data.Amount,
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
		return dto.OfferDTO{}, fmt.Errorf("fastbank: unexpected status %d: %s", resp.StatusCode, string(body))
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

func (b *FastBank) GetApplication(ctx context.Context, id string) (dto.OfferDTO, error) {
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
		return dto.OfferDTO{}, fmt.Errorf("fastbank: unexpected status %d: %s", resp.StatusCode, string(body))
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
