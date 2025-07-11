package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const (
	MaritalStatusSingle     string = "SINGLE"
	MaritalStatusMarried    string = "MARRIED"
	MaritalStatusDivorced   string = "DIVORCED"
	MaritalStatusCohabiting string = "COHABITING"

	OfferStatusDraft     string = "DRAFT"
	OfferStatusProcessed string = "PROCESSED"
)

type Application struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`

	Phone                    string  `json:"phone"`
	Email                    string  `json:"email"`
	MonthlyIncome            float64 `json:"monthlyIncome"`
	MonthlyExpenses          float64 `json:"monthlyExpenses"`
	MonthlyCreditLiabilities float64 `json:"monthlyCreditLiabilities"`
	MaritalStatus            string  `gorm:"type:marital_status_enum" json:"maritalStatus"`
	Dependents               int     `json:"dependents"`
	AgreeToDataSharing       bool    `json:"agreeToDataSharing"`
	AgreeToBeScored          bool    `json:"agreeToBeScored"`
	Amount                   float64 `json:"amount"`
	Offers                   []Offer `gorm:"foreignKey:ApplicationID" json:"offers"`
}

func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}
