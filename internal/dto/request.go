package dto

import "time"

// Auth
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Tenant
type CreateTenantRequest struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contactPerson"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Status        string `json:"status"`
}

type UpdateTenantRequest struct {
	Name          string `json:"name"`
	ContactPerson string `json:"contactPerson"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Status        string `json:"status"`
}

type TenantListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
}

// Contract
type CreateContractRequest struct {
	TenantID   uint      `json:"tenantId" binding:"required"`
	ContractNo string    `json:"contractNo"`
	StartDate  time.Time `json:"startDate" binding:"required"`
	EndDate    time.Time `json:"endDate" binding:"required"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
}

type UpdateContractRequest struct {
	TenantID   uint      `json:"tenantId"`
	ContractNo string    `json:"contractNo"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"`
}

type ContractListRequest struct {
	Page          int    `form:"page,default=1"`
	PageSize      int    `form:"pageSize,default=10"`
	Keyword       string `form:"keyword"`
	Status        string `form:"status"`
	StartDateFrom string `form:"startDateFrom"`
	StartDateTo   string `form:"startDateTo"`
}

// Room
type CreateRoomRequest struct {
	RoomNo      string  `json:"roomNo" binding:"required"`
	Building    string  `json:"building"`
	Floor       int     `json:"floor"`
	Area        float64 `json:"area"`
	MonthlyRent float64 `json:"monthlyRent"`
	Status      string  `json:"status"`
}

type UpdateRoomRequest struct {
	RoomNo      string  `json:"roomNo"`
	Building    string  `json:"building"`
	Floor       int     `json:"floor"`
	Area        float64 `json:"area"`
	MonthlyRent float64 `json:"monthlyRent"`
	Status      string  `json:"status"`
}

type RoomListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Building string `form:"building"`
	Status   string `form:"status"`
}

type AssignTenantRequest struct {
	TenantID uint `json:"tenantId" binding:"required"`
}

// Fee
type CreateFeeRequest struct {
	TenantID uint      `json:"tenantId" binding:"required"`
	RoomNo   string    `json:"roomNo"`
	FeeType  string    `json:"feeType" binding:"required"`
	Amount   float64   `json:"amount" binding:"required"`
	Period   string    `json:"period"`
	DueDate  time.Time `json:"dueDate" binding:"required"`
	Status   string    `json:"status"`
}

type UpdateFeeRequest struct {
	TenantID uint      `json:"tenantId"`
	RoomNo   string    `json:"roomNo"`
	FeeType  string    `json:"feeType"`
	Amount   float64   `json:"amount"`
	Period   string    `json:"period"`
	DueDate  time.Time `json:"dueDate"`
	Status   string    `json:"status"`
}

type FeeListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	TenantID uint   `form:"tenantId"`
	RoomNo   string `form:"roomNo"`
	FeeType  string `form:"feeType"`
	Status   string `form:"status"`
	Period   string `form:"period"`
}

type PayFeeRequest struct {
	PaidDate *time.Time `json:"paidDate"`
}

// Maintenance
type CreateMaintenanceRequest struct {
	TenantID    uint   `json:"tenantId" binding:"required"`
	RoomNo      string `json:"roomNo"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type UpdateMaintenanceRequest struct {
	TenantID    uint   `json:"tenantId"`
	RoomNo      string `json:"roomNo"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Assignee    string `json:"assignee"`
}

type MaintenanceListRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"pageSize,default=10"`
	Keyword  string `form:"keyword"`
	Type     string `form:"type"`
	Status   string `form:"status"`
	Priority string `form:"priority"`
}

type AssignMaintenanceRequest struct {
	Assignee string `json:"assignee" binding:"required"`
}

type CompleteMaintenanceRequest struct {
	CompletedAt *time.Time `json:"completedAt"`
}

// Report
type ReportQueryRequest struct {
	Start   string `form:"start"`
	End     string `form:"end"`
	GroupBy string `form:"groupBy"`
	Limit   int    `form:"limit,default=10"`
}
