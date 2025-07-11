package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Offer struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	ApplicationID        uuid.UUID `json:"applicationId"`
	ExternalID           string    `json:"externalId"`
	Bank                 string    `json:"bank"`
	Status               string    `gorm:"type:application_status_enum" json:"status"`
	MonthlyPaymentAmount float64   `json:"monthlyPaymentAmount"`
	TotalRepaymentAmount float64   `json:"totalRepaymentAmount"`
	NumberOfPayments     int       `json:"numberOfPayments"`
	AnnualPercentageRate float64   `json:"annualPercentageRate"`
	FirstRepaymentDate   time.Time `json:"firstRepaymentDate"`
}

func (o *Offer) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
