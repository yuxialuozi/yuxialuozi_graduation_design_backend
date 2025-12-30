package service

import (
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

type TenantService struct {
	tenantRepo *repository.TenantRepository
}

func NewTenantService(tenantRepo *repository.TenantRepository) *TenantService {
	return &TenantService{tenantRepo: tenantRepo}
}

func (s *TenantService) Create(tenant *model.Tenant) error {
	return s.tenantRepo.Create(tenant)
}

func (s *TenantService) GetByID(id uint) (*model.Tenant, error) {
	return s.tenantRepo.FindByID(id)
}

func (s *TenantService) Update(tenant *model.Tenant) error {
	return s.tenantRepo.Update(tenant)
}

func (s *TenantService) Delete(id uint) error {
	return s.tenantRepo.Delete(id)
}

func (s *TenantService) List(page, pageSize int, keyword, status string) ([]model.Tenant, int64, error) {
	return s.tenantRepo.List(page, pageSize, keyword, status)
}

func (s *TenantService) GetAll() ([]model.Tenant, error) {
	return s.tenantRepo.FindAll()
}
