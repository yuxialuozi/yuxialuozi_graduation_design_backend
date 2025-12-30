package model

import (
	"time"
)

type Fee struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TenantID   uint       `gorm:"not null;index" json:"tenantId"`
	Tenant     Tenant     `gorm:"foreignKey:TenantID" json:"-"`
	TenantName string     `gorm:"-" json:"tenantName"`
	RoomNo     string     `gorm:"size:20" json:"roomNo"`
	FeeType    string     `gorm:"size:20;not null" json:"feeType"`
	Amount     float64    `gorm:"type:decimal(10,2)" json:"amount"`
	Period     string     `gorm:"size:20" json:"period"`
	DueDate    time.Time  `json:"dueDate"`
	PaidDate   *time.Time `json:"paidDate"`
	Status     string     `gorm:"size:20;default:'unpaid'" json:"status"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

func (Fee) TableName() string {
	return "fees"
}
