package repository

import (
	"time"

	"gorm.io/gorm"

	"yuxialuozi_graduation_design_backend/internal/model"
)

type FeeRepository struct {
	db *gorm.DB
}

func NewFeeRepository(db *gorm.DB) *FeeRepository {
	return &FeeRepository{db: db}
}

func (r *FeeRepository) Create(fee *model.Fee) error {
	return r.db.Create(fee).Error
}

func (r *FeeRepository) FindByID(id uint) (*model.Fee, error) {
	var fee model.Fee
	if err := r.db.Preload("Tenant").First(&fee, id).Error; err != nil {
		return nil, err
	}
	fee.TenantName = fee.Tenant.Name
	return &fee, nil
}

func (r *FeeRepository) Update(fee *model.Fee) error {
	return r.db.Save(fee).Error
}

func (r *FeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Fee{}, id).Error
}

func (r *FeeRepository) List(page, pageSize int, tenantID uint, roomNo, feeType, status, period string) ([]model.Fee, int64, error) {
	var fees []model.Fee
	var total int64

	query := r.db.Model(&model.Fee{}).Preload("Tenant")

	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if roomNo != "" {
		query = query.Where("room_no = ?", roomNo)
	}
	if feeType != "" {
		query = query.Where("fee_type = ?", feeType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if period != "" {
		query = query.Where("period = ?", period)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("due_date DESC").Find(&fees).Error; err != nil {
		return nil, 0, err
	}

	for i := range fees {
		fees[i].TenantName = fees[i].Tenant.Name
	}

	return fees, total, nil
}

func (r *FeeRepository) SumByTypeAndPeriod(feeType string, start, end time.Time) (float64, error) {
	var sum float64
	err := r.db.Model(&model.Fee{}).
		Where("fee_type = ? AND status = 'paid' AND COALESCE(paid_date, due_date, created_at) >= ? AND COALESCE(paid_date, due_date, created_at) <= ?", feeType, start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *FeeRepository) SumByPeriod(start, end time.Time) (float64, error) {
	var sum float64
	err := r.db.Model(&model.Fee{}).
		Where("status = 'paid' AND COALESCE(paid_date, due_date, created_at) >= ? AND COALESCE(paid_date, due_date, created_at) <= ?", start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	return sum, err
}

func (r *FeeRepository) CountByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Fee{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FeeRepository) SumUnpaidAmount() (float64, error) {
	var sum float64
	err := r.db.Model(&model.Fee{}).
		Where("status IN ('unpaid', 'overdue')").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	return sum, err
}

type FeeComposition struct {
	FeeType string  `json:"feeType"`
	Amount  float64 `json:"amount"`
}

func (r *FeeRepository) GetComposition(start, end time.Time) ([]FeeComposition, error) {
	var compositions []FeeComposition
	err := r.db.Model(&model.Fee{}).
		Select("fee_type, COALESCE(SUM(amount), 0) as amount").
		Where("status = 'paid' AND COALESCE(paid_date, due_date, created_at) >= ? AND COALESCE(paid_date, due_date, created_at) <= ?", start, end).
		Group("fee_type").
		Scan(&compositions).Error
	return compositions, err
}

type IncomeByMonth struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

func (r *FeeRepository) GetIncomeByMonth(start, end time.Time) ([]IncomeByMonth, error) {
	var incomes []IncomeByMonth
	err := r.db.Model(&model.Fee{}).
		Select("TO_CHAR(COALESCE(paid_date, due_date, created_at), 'YYYY-MM') as month, COALESCE(SUM(amount), 0) as amount").
		Where("status = 'paid' AND COALESCE(paid_date, due_date, created_at) >= ? AND COALESCE(paid_date, due_date, created_at) <= ?", start, end).
		Group("TO_CHAR(COALESCE(paid_date, due_date, created_at), 'YYYY-MM')").
		Order("month ASC").
		Scan(&incomes).Error
	return incomes, err
}

type IncomeByDay struct {
	Day    string  `json:"day"`
	Amount float64 `json:"amount"`
}

func (r *FeeRepository) GetIncomeByDay(start, end time.Time) ([]IncomeByDay, error) {
	var incomes []IncomeByDay
	err := r.db.Model(&model.Fee{}).
		Select("DATE(COALESCE(paid_date, due_date, created_at)) as day, COALESCE(SUM(amount), 0) as amount").
		Where("status = 'paid' AND COALESCE(paid_date, due_date, created_at) >= ? AND COALESCE(paid_date, due_date, created_at) <= ?", start, end).
		Group("DATE(COALESCE(paid_date, due_date, created_at))").
		Order("day ASC").
		Scan(&incomes).Error
	return incomes, err
}

type TenantFeeRanking struct {
	TenantID   uint    `json:"tenantId"`
	TenantName string  `json:"tenantName"`
	Amount     float64 `json:"amount"`
}

func (r *FeeRepository) GetTenantRanking(limit int, start, end time.Time) ([]TenantFeeRanking, error) {
	var rankings []TenantFeeRanking
	err := r.db.Model(&model.Fee{}).
		Select("fees.tenant_id, tenants.name as tenant_name, COALESCE(SUM(fees.amount), 0) as amount").
		Joins("LEFT JOIN tenants ON fees.tenant_id = tenants.id").
		Where("fees.status = 'paid' AND COALESCE(fees.paid_date, fees.due_date, fees.created_at) >= ? AND COALESCE(fees.paid_date, fees.due_date, fees.created_at) <= ?", start, end).
		Group("fees.tenant_id, tenants.name").
		Order("amount DESC").
		Limit(limit).
		Scan(&rankings).Error
	return rankings, err
}

// SumByType 按费用类型求和
func (r *FeeRepository) SumByType(feeType string) (float64, error) {
	var sum float64
	err := r.db.Model(&model.Fee{}).
		Where("fee_type = ? AND status = 'paid'", feeType).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum).Error
	return sum, err
}
