package service

import (
	"errors"
	"fmt"
	"time"

	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

type ContractService struct {
	contractRepo *repository.ContractRepository
	tenantRepo   *repository.TenantRepository
}

func NewContractService(contractRepo *repository.ContractRepository, tenantRepo *repository.TenantRepository) *ContractService {
	return &ContractService{
		contractRepo: contractRepo,
		tenantRepo:   tenantRepo,
	}
}

func (s *ContractService) Create(contract *model.Contract) (*model.Contract, error) {
	if contract.ContractNo == "" {
		contract.ContractNo = s.generateContractNo()
	}
	if err := s.contractRepo.Create(contract); err != nil {
		return nil, err
	}
	// 设置 TenantName
	if contract.TenantID > 0 {
		tenant, err := s.tenantRepo.FindByID(contract.TenantID)
		if err == nil && tenant != nil {
			contract.TenantName = tenant.Name
		}
	}
	return contract, nil
}

func (s *ContractService) GetByID(id uint) (*model.Contract, error) {
	return s.contractRepo.FindByID(id)
}

func (s *ContractService) Update(contract *model.Contract) error {
	return s.contractRepo.Update(contract)
}

func (s *ContractService) Delete(id uint) error {
	return s.contractRepo.Delete(id)
}

func (s *ContractService) List(page, pageSize int, keyword, status string, startDateFrom, startDateTo *time.Time, tenantID uint) ([]model.Contract, int64, error) {
	return s.contractRepo.List(page, pageSize, keyword, status, startDateFrom, startDateTo, tenantID)
}

func (s *ContractService) Activate(id uint) error {
	contract, err := s.contractRepo.FindByID(id)
	if err != nil {
		return errors.New("合同不存在")
	}
	if contract.Status != "draft" {
		return errors.New("只有草稿状态的合同可以激活")
	}
	contract.Status = "active"
	return s.contractRepo.Update(contract)
}

func (s *ContractService) Terminate(id uint) error {
	contract, err := s.contractRepo.FindByID(id)
	if err != nil {
		return errors.New("合同不存在")
	}
	if contract.Status == "terminated" {
		return errors.New("合同已终止")
	}
	contract.Status = "terminated"
	return s.contractRepo.Update(contract)
}

func (s *ContractService) GetExpiring(days int) ([]model.Contract, error) {
	return s.contractRepo.GetExpiring(days)
}

func (s *ContractService) generateContractNo() string {
	return fmt.Sprintf("HT%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
}
