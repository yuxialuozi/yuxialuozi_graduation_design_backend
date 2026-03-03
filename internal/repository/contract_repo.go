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

// 基于合同的收入结构
type ContractIncomeByDay struct {
	Day    string  `json:"day"`
	Amount float64 `json:"amount"`
}

// GetIncomeByDay 基于合同按天统计收入
func (r *ContractRepository) GetIncomeByDay(start, end time.Time) ([]ContractIncomeByDay, error) {
	// 获取在指定时间范围内有效的活跃合同
	var contracts []model.Contract
	err := r.db.Where("status = ? AND start_date <= ? AND end_date >= ?", "active", end, start).Find(&contracts).Error
	if err != nil {
		return nil, err
	}

	// 创建结果切片
	var incomes []ContractIncomeByDay

	// 计算指定日期范围内的每一天
	days := int(end.Sub(start).Hours()/24) + 1
	for i := 0; i < days; i++ {
		currentDate := start.AddDate(0, 0, i)
		dayIncome := 0.0

		for _, contract := range contracts {
			// 检查合同在当前日期是否有效
			contractStart := contract.StartDate.Time
			contractEnd := contract.EndDate.Time

			if !currentDate.Before(contractStart) && !currentDate.After(contractEnd) {
				// 计算合同的总天数
				totalDays := contractEnd.Sub(contractStart).Hours()/24 + 1
				if totalDays > 0 {
					// 按天分配合同收入
					dailyIncome := contract.Amount / totalDays
					dayIncome += dailyIncome
				}
			}
		}

		if dayIncome > 0 {
			incomes = append(incomes, ContractIncomeByDay{
				Day:    currentDate.Format("2006-01-02"),
				Amount: dayIncome,
			})
		}
	}

	return incomes, nil
}

// GetTotalIncomeByPeriod 基于合同计算指定时间段的总收入
func (r *ContractRepository) GetTotalIncomeByPeriod(start, end time.Time) (float64, error) {
	// 获取在指定时间范围内有效的活跃合同
	var contracts []model.Contract
	err := r.db.Where("status = ? AND start_date <= ? AND end_date >= ?", "active", end, start).Find(&contracts).Error
	if err != nil {
		return 0, err
	}

	totalIncome := 0.0
	for _, contract := range contracts {
		// 计算合同与查询时间的重叠天数
		contractStart := contract.StartDate.Time
		contractEnd := contract.EndDate.Time

		overlapStart := contractStart
		if start.After(contractStart) {
			overlapStart = start
		}

		overlapEnd := contractEnd
		if end.Before(contractEnd) {
			overlapEnd = end
		}

		// 如果有重叠
		if !overlapEnd.Before(overlapStart) {
			overlapDays := overlapEnd.Sub(overlapStart).Hours()/24 + 1
			totalDays := contractEnd.Sub(contractStart).Hours()/24 + 1

			if totalDays > 0 {
				// 按比例计算收入
				income := contract.Amount * (overlapDays / totalDays)
				totalIncome += income
			}
		}
	}

	return totalIncome, nil
}

// GetTenantIncomeByContract 获取租户基于合同的收入
type TenantContractIncome struct {
	TenantID   uint    `json:"tenantId"`
	TenantName string  `json:"tenantName"`
	Amount     float64 `json:"amount"`
}

func (r *ContractRepository) GetTenantIncomeByContract(limit int, start, end time.Time) ([]TenantContractIncome, error) {
	// 获取在指定时间范围内有效的活跃合同及其租户信息
	var results []TenantContractIncome

	var contractIncomes []struct {
		TenantID   uint
		TenantName string
		Amount     float64
	}

	err := r.db.Table("contracts").
		Select("contracts.tenant_id, tenants.name as tenant_name, contracts.amount").
		Joins("LEFT JOIN tenants ON contracts.tenant_id = tenants.id").
		Where("contracts.status = ? AND contracts.start_date <= ? AND contracts.end_date >= ?", "active", end, start).
		Scan(&contractIncomes).Error

	if err != nil {
		return nil, err
	}

	// 按租户汇总收入
	tenantMap := make(map[uint]*TenantContractIncome)
	for _, ci := range contractIncomes {
		if _, exists := tenantMap[ci.TenantID]; !exists {
			tenantMap[ci.TenantID] = &TenantContractIncome{
				TenantID:   ci.TenantID,
				TenantName: ci.TenantName,
				Amount:     0,
			}
		}
		tenantMap[ci.TenantID].Amount += ci.Amount
	}

	// 转换为切片并排序
	for _, value := range tenantMap {
		results = append(results, *value)
	}

	// 按金额降序排序
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Amount > results[i].Amount {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// 限制返回数量
	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// GetContractComposition 基于合同的费用构成
type ContractFeeComposition struct {
	FeeType string  `json:"feeType"`
	Amount  float64 `json:"amount"`
}

func (r *ContractRepository) GetContractComposition(start, end time.Time) ([]ContractFeeComposition, error) {
	// 获取在指定时间范围内合同的收入类型构成
	var compositions []ContractFeeComposition

	// 这里假设合同类型可以映射到费用类型
	// 由于当前合同模型没有类型字段，我们使用合同金额作为综合收入
	var contracts []model.Contract
	err := r.db.Where("status = ? AND start_date <= ? AND end_date >= ?", "active", end, start).Find(&contracts).Error
	if err != nil {
		return nil, err
	}

	// 计算总收入
	totalAmount := 0.0
	for _, contract := range contracts {
		totalAmount += contract.Amount
	}

	// 按不同类型分配（这里简化处理，实际应用中可能需要更复杂的逻辑）
	if totalAmount > 0 {
		// 假设70%是租金，10%水费，15%电费，5%其他
		compositions = []ContractFeeComposition{
			{FeeType: "rent", Amount: totalAmount * 0.70},
			{FeeType: "water", Amount: totalAmount * 0.10},
			{FeeType: "electricity", Amount: totalAmount * 0.15},
			{FeeType: "property", Amount: totalAmount * 0.05},
		}
	}

	return compositions, nil
}
