package model

type Fee struct {
	ID         uint        `gorm:"primaryKey" json:"id"`
	TenantID   uint        `gorm:"not null;index" json:"tenantId"`
	Tenant     Tenant      `gorm:"foreignKey:TenantID" json:"-"`
	TenantName string      `gorm:"-" json:"tenantName"`
	RoomNo     string      `gorm:"size:20" json:"roomNo"`
	FeeType    string      `gorm:"size:20;not null" json:"feeType"`
	Amount     float64     `gorm:"type:decimal(10,2)" json:"amount"`
	Period     string      `gorm:"size:20" json:"period"`
	DueDate    CustomTime  `json:"dueDate"`
	PaidDate   *CustomTime `json:"paidDate"`
	Status     string      `gorm:"size:20;default:'unpaid'" json:"status"`
	CreatedAt  CustomTime  `json:"createdAt"`
	UpdatedAt  CustomTime  `json:"updatedAt"`
}

func (Fee) TableName() string {
	return "fees"
}
