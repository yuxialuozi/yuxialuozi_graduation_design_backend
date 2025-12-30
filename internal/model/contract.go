package model

import (
	"time"
)

type Contract struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   uint      `gorm:"not null;index" json:"tenantId"`
	Tenant     Tenant    `gorm:"foreignKey:TenantID" json:"-"`
	TenantName string    `gorm:"-" json:"tenantName"`
	ContractNo string    `gorm:"uniqueIndex;size:50;not null" json:"contractNo"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	Amount     float64   `gorm:"type:decimal(10,2)" json:"amount"`
	Status     string    `gorm:"size:20;default:'draft'" json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (Contract) TableName() string {
	return "contracts"
}
