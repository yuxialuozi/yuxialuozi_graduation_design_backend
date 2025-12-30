package model

import (
	"time"
)

type Room struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomNo      string    `gorm:"uniqueIndex;size:20;not null" json:"roomNo"`
	Building    string    `gorm:"size:50" json:"building"`
	Floor       int       `json:"floor"`
	Area        float64   `gorm:"type:decimal(10,2)" json:"area"`
	MonthlyRent float64   `gorm:"type:decimal(10,2)" json:"monthlyRent"`
	Status      string    `gorm:"size:20;default:'vacant'" json:"status"`
	TenantID    *uint     `gorm:"index" json:"tenantId"`
	Tenant      *Tenant   `gorm:"foreignKey:TenantID" json:"-"`
	TenantName  string    `gorm:"-" json:"tenantName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Room) TableName() string {
	return "rooms"
}
