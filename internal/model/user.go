package model

import (
	"github.com/lib/pq"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Nickname    string         `gorm:"size:50" json:"nickname"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Email       string         `gorm:"size:100" json:"email"`
	Role        string         `gorm:"size:20;default:'user'" json:"role"`
	Status      string         `gorm:"size:20;default:'active'" json:"status"`
	Permissions pq.StringArray `gorm:"type:text[]" json:"permissions" swaggertype:"array,string"`
	TenantID    uint           `gorm:"index" json:"tenantId"`
	Tenant      *Tenant        `gorm:"foreignKey:TenantID" json:"-"`
	CreatedAt   CustomTime     `json:"createdAt"`
	UpdatedAt   CustomTime     `json:"updatedAt"`
}

// User represents a user in the system.
// @Description User represents a user in the system with role-based access control.
type UserResponse struct {
	ID          uint     `json:"id" example:"1"`
	Username    string   `json:"username" example:"admin"`
	Nickname    string   `json:"nickname" example:"管理员"`
	Avatar      string   `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone       string   `json:"phone" example:"13800138000"`
	Email       string   `json:"email" example:"admin@example.com"`
	Role        string   `json:"role" example:"admin"`
	Status      string   `json:"status" example:"active"`
	Permissions []string `json:"permissions" swaggertype:"array,string" example:"[\"read\", \"write\"]"`
	TenantID    uint     `json:"tenantId" example:"1"`
	CreatedAt   string   `json:"createdAt" example:"2024-01-01 00:00:00"`
	UpdatedAt   string   `json:"updatedAt" example:"2024-01-01 00:00:00"`
}

func (User) TableName() string {
	return "users"
}
