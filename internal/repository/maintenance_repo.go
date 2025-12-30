package repository

import (
	"time"

	"gorm.io/gorm"

	"yuxialuozi_graduation_design_backend/internal/model"
)

type MaintenanceRepository struct {
	db *gorm.DB
}

func NewMaintenanceRepository(db *gorm.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) Create(maintenance *model.Maintenance) error {
	return r.db.Create(maintenance).Error
}

func (r *MaintenanceRepository) FindByID(id uint) (*model.Maintenance, error) {
	var maintenance model.Maintenance
	if err := r.db.Preload("Tenant").First(&maintenance, id).Error; err != nil {
		return nil, err
	}
	maintenance.TenantName = maintenance.Tenant.Name
	return &maintenance, nil
}

func (r *MaintenanceRepository) FindByTicketNo(ticketNo string) (*model.Maintenance, error) {
	var maintenance model.Maintenance
	if err := r.db.Where("ticket_no = ?", ticketNo).First(&maintenance).Error; err != nil {
		return nil, err
	}
	return &maintenance, nil
}

func (r *MaintenanceRepository) Update(maintenance *model.Maintenance) error {
	return r.db.Save(maintenance).Error
}

func (r *MaintenanceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Maintenance{}, id).Error
}

func (r *MaintenanceRepository) List(page, pageSize int, keyword, maintenanceType, status, priority string) ([]model.Maintenance, int64, error) {
	var maintenances []model.Maintenance
	var total int64

	query := r.db.Model(&model.Maintenance{}).Preload("Tenant")

	if keyword != "" {
		query = query.Where("ticket_no ILIKE ? OR description ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if maintenanceType != "" {
		query = query.Where("type = ?", maintenanceType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&maintenances).Error; err != nil {
		return nil, 0, err
	}

	for i := range maintenances {
		maintenances[i].TenantName = maintenances[i].Tenant.Name
	}

	return maintenances, total, nil
}

func (r *MaintenanceRepository) CountByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Maintenance{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

type MaintenanceStats struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

func (r *MaintenanceRepository) GetStatsByType(start, end time.Time) ([]MaintenanceStats, error) {
	var stats []MaintenanceStats
	err := r.db.Model(&model.Maintenance{}).
		Select("type, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Group("type").
		Scan(&stats).Error
	return stats, err
}

type MaintenanceStatusStats struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

func (r *MaintenanceRepository) GetStatsByStatus(start, end time.Time) ([]MaintenanceStatusStats, error) {
	var stats []MaintenanceStatusStats
	err := r.db.Model(&model.Maintenance{}).
		Select("status, COUNT(*) as count").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Group("status").
		Scan(&stats).Error
	return stats, err
}

func (r *MaintenanceRepository) GetLastTicketNo() (string, error) {
	var maintenance model.Maintenance
	err := r.db.Order("created_at DESC").First(&maintenance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}
	return maintenance.TicketNo, nil
}
