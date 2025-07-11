package banks

import (
	"context"
	"financing-aggregator/internal/dto"
)

type (
	Bank interface {
		Name() string
		SubmitApplication(ctx context.Context, data dto.ApplicationDTO) (dto.OfferDTO, error)
		GetApplication(ctx context.Context, id string) (dto.OfferDTO, error)
	}
)
