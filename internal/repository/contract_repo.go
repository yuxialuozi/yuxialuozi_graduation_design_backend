package repository

import (
	"time"

	"gorm.io/gorm"

	"yuxialuozi_graduation_design_backend/internal/model"
)

type ContractRepository struct {
	db *gorm.DB
}

func NewContractRepository(db *gorm.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) Create(contract *model.Contract) error {
	return r.db.Create(contract).Error
}

func (r *ContractRepository) FindByID(id uint) (*model.Contract, error) {
	var contract model.Contract
	if err := r.db.Preload("Tenant").First(&contract, id).Error; err != nil {
		return nil, err
	}
	contract.TenantName = contract.Tenant.Name
	return &contract, nil
}

func (r *ContractRepository) FindByContractNo(contractNo string) (*model.Contract, error) {
	var contract model.Contract
	if err := r.db.Where("contract_no = ?", contractNo).First(&contract).Error; err != nil {
		return nil, err
	}
	return &contract, nil
}

func (r *ContractRepository) Update(contract *model.Contract) error {
	return r.db.Save(contract).Error
}

func (r *ContractRepository) Delete(id uint) error {
	return r.db.Delete(&model.Contract{}, id).Error
}

func (r *ContractRepository) List(page, pageSize int, keyword, status string, startDateFrom, startDateTo *time.Time) ([]model.Contract, int64, error) {
	var contracts []model.Contract
	var total int64

	query := r.db.Model(&model.Contract{}).Preload("Tenant")

	if keyword != "" {
		query = query.Where("contract_no ILIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDateFrom != nil {
		query = query.Where("start_date >= ?", startDateFrom)
	}
	if startDateTo != nil {
		query = query.Where("start_date <= ?", startDateTo)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&contracts).Error; err != nil {
		return nil, 0, err
	}

	for i := range contracts {
		contracts[i].TenantName = contracts[i].Tenant.Name
	}

	return contracts, total, nil
}

func (r *ContractRepository) FindByTenantID(tenantID uint) ([]model.Contract, error) {
	var contracts []model.Contract
	if err := r.db.Where("tenant_id = ?", tenantID).Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *ContractRepository) CountByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Contract{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
