package service

import (
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

func (s *ContractService) Create(contract *model.Contract) error {
	if contract.ContractNo == "" {
		contract.ContractNo = s.generateContractNo()
	}
	return s.contractRepo.Create(contract)
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

func (s *ContractService) List(page, pageSize int, keyword, status string, startDateFrom, startDateTo *time.Time) ([]model.Contract, int64, error) {
	return s.contractRepo.List(page, pageSize, keyword, status, startDateFrom, startDateTo)
}

func (s *ContractService) generateContractNo() string {
	return fmt.Sprintf("HT%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
}
