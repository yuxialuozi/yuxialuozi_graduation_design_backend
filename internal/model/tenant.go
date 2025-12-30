package model

import (
	"time"
)

type Tenant struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	ContactPerson string    `gorm:"size:50" json:"contactPerson"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Email         string    `gorm:"size:100" json:"email"`
	Status        string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (Tenant) TableName() string {
	return "tenants"
}
