package service

import (
	"fmt"
	"time"

	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

type MaintenanceService struct {
	maintenanceRepo *repository.MaintenanceRepository
	tenantRepo      *repository.TenantRepository
}

func NewMaintenanceService(maintenanceRepo *repository.MaintenanceRepository, tenantRepo *repository.TenantRepository) *MaintenanceService {
	return &MaintenanceService{
		maintenanceRepo: maintenanceRepo,
		tenantRepo:      tenantRepo,
	}
}

func (s *MaintenanceService) Create(maintenance *model.Maintenance) error {
	if maintenance.TicketNo == "" {
		maintenance.TicketNo = s.generateTicketNo()
	}
	return s.maintenanceRepo.Create(maintenance)
}

func (s *MaintenanceService) GetByID(id uint) (*model.Maintenance, error) {
	return s.maintenanceRepo.FindByID(id)
}

func (s *MaintenanceService) Update(maintenance *model.Maintenance) error {
	return s.maintenanceRepo.Update(maintenance)
}

func (s *MaintenanceService) Delete(id uint) error {
	return s.maintenanceRepo.Delete(id)
}

func (s *MaintenanceService) List(page, pageSize int, keyword, maintenanceType, status, priority string) ([]model.Maintenance, int64, error) {
	return s.maintenanceRepo.List(page, pageSize, keyword, maintenanceType, status, priority)
}

func (s *MaintenanceService) Assign(id uint, assignee string) error {
	maintenance, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return err
	}

	maintenance.Assignee = assignee
	maintenance.Status = "processing"
	return s.maintenanceRepo.Update(maintenance)
}

func (s *MaintenanceService) Complete(id uint, completedAt *time.Time) error {
	maintenance, err := s.maintenanceRepo.FindByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	if completedAt == nil {
		completedAt = &now
	}

	maintenance.CompletedAt = completedAt
	maintenance.Status = "completed"
	return s.maintenanceRepo.Update(maintenance)
}

func (s *MaintenanceService) generateTicketNo() string {
	return fmt.Sprintf("WX%s%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
}
