package service

import (
	"time"

	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

type FeeService struct {
	feeRepo    *repository.FeeRepository
	tenantRepo *repository.TenantRepository
}

func NewFeeService(feeRepo *repository.FeeRepository, tenantRepo *repository.TenantRepository) *FeeService {
	return &FeeService{
		feeRepo:    feeRepo,
		tenantRepo: tenantRepo,
	}
}

func (s *FeeService) Create(fee *model.Fee) error {
	return s.feeRepo.Create(fee)
}

func (s *FeeService) GetByID(id uint) (*model.Fee, error) {
	return s.feeRepo.FindByID(id)
}

func (s *FeeService) Update(fee *model.Fee) error {
	return s.feeRepo.Update(fee)
}

func (s *FeeService) Delete(id uint) error {
	return s.feeRepo.Delete(id)
}

func (s *FeeService) List(page, pageSize int, tenantID uint, roomNo, feeType, status, period string) ([]model.Fee, int64, error) {
	return s.feeRepo.List(page, pageSize, tenantID, roomNo, feeType, status, period)
}

func (s *FeeService) Pay(id uint, paidDate *time.Time) error {
	fee, err := s.feeRepo.FindByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	if paidDate == nil {
		paidDate = &now
	}

	fee.PaidDate = paidDate
	fee.Status = "paid"
	return s.feeRepo.Update(fee)
}
