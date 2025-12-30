package model

import (
	"time"

	"github.com/lib/pq"
)

type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Nickname    string         `gorm:"size:50" json:"nickname"`
	Avatar      string         `gorm:"size:255" json:"avatar"`
	Role        string         `gorm:"size:20;default:'user'" json:"role"`
	Permissions pq.StringArray `gorm:"type:text[]" json:"permissions" swaggertype:"array,string"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// User represents a user in the system.
// @Description User represents a user in the system with role-based access control.
type UserResponse struct {
	ID          uint     `json:"id" example:"1"`
	Username    string   `json:"username" example:"admin"`
	Nickname    string   `json:"nickname" example:"管理员"`
	Avatar      string   `json:"avatar" example:"https://example.com/avatar.jpg"`
	Role        string   `json:"role" example:"admin"`
	Permissions []string `json:"permissions" swaggertype:"array,string" example:"[\"read\", \"write\"]"`
	CreatedAt   string   `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   string   `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}

func (User) TableName() string {
	return "users"
}