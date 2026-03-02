package model

type Maintenance struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	TicketNo    string      `gorm:"uniqueIndex;size:50;not null" json:"ticketNo"`
	TenantID    uint        `gorm:"not null;index" json:"tenantId"`
	Tenant      Tenant      `gorm:"foreignKey:TenantID" json:"-"`
	TenantName  string      `gorm:"-" json:"tenantName"`
	RoomNo      string      `gorm:"size:20" json:"roomNo"`
	Type        string      `gorm:"size:20" json:"type"`
	Description string      `gorm:"type:text" json:"description"`
	Priority    string      `gorm:"size:20;default:'medium'" json:"priority"`
	Status      string      `gorm:"size:20;default:'pending'" json:"status"`
	Assignee    string      `gorm:"size:50" json:"assignee"`
	CreatedAt   CustomTime  `json:"createdAt"`
	CompletedAt *CustomTime `json:"completedAt"`
	UpdatedAt   CustomTime  `json:"updatedAt"`
}

func (Maintenance) TableName() string {
	return "maintenances"
}
