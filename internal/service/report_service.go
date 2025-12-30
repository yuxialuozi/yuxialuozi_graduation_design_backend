package service

import (
	"time"

	"yuxialuozi_graduation_design_backend/internal/repository"
)

type ReportService struct {
	feeRepo         *repository.FeeRepository
	roomRepo        *repository.RoomRepository
	maintenanceRepo *repository.MaintenanceRepository
	tenantRepo      *repository.TenantRepository
	contractRepo    *repository.ContractRepository
}

func NewReportService(
	feeRepo *repository.FeeRepository,
	roomRepo *repository.RoomRepository,
	maintenanceRepo *repository.MaintenanceRepository,
	tenantRepo *repository.TenantRepository,
	contractRepo *repository.ContractRepository,
) *ReportService {
	return &ReportService{
		feeRepo:         feeRepo,
		roomRepo:        roomRepo,
		maintenanceRepo: maintenanceRepo,
		tenantRepo:      tenantRepo,
		contractRepo:    contractRepo,
	}
}

type IncomeReport struct {
	Total     float64                       `json:"total"`
	ByMonth   []repository.IncomeByMonth    `json:"byMonth"`
	ByType    []repository.FeeComposition   `json:"byType"`
}

func (s *ReportService) GetIncomeReport(start, end time.Time, groupBy string) (*IncomeReport, error) {
	total, err := s.feeRepo.SumByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	byMonth, err := s.feeRepo.GetIncomeByMonth(start, end)
	if err != nil {
		return nil, err
	}

	byType, err := s.feeRepo.GetComposition(start, end)
	if err != nil {
		return nil, err
	}

	return &IncomeReport{
		Total:   total,
		ByMonth: byMonth,
		ByType:  byType,
	}, nil
}

type OccupancyReport struct {
	TotalRooms    int64   `json:"totalRooms"`
	OccupiedRooms int64   `json:"occupiedRooms"`
	VacantRooms   int64   `json:"vacantRooms"`
	OccupancyRate float64 `json:"occupancyRate"`
}

func (s *ReportService) GetOccupancyReport(start, end time.Time) (*OccupancyReport, error) {
	totalRooms, err := s.roomRepo.CountTotal()
	if err != nil {
		return nil, err
	}

	occupiedRooms, err := s.roomRepo.CountByStatus("occupied")
	if err != nil {
		return nil, err
	}

	vacantRooms, err := s.roomRepo.CountByStatus("vacant")
	if err != nil {
		return nil, err
	}

	var occupancyRate float64
	if totalRooms > 0 {
		occupancyRate = float64(occupiedRooms) / float64(totalRooms) * 100
	}

	return &OccupancyReport{
		TotalRooms:    totalRooms,
		OccupiedRooms: occupiedRooms,
		VacantRooms:   vacantRooms,
		OccupancyRate: occupancyRate,
	}, nil
}

func (s *ReportService) GetFeeComposition(start, end time.Time) ([]repository.FeeComposition, error) {
	return s.feeRepo.GetComposition(start, end)
}

type MaintenanceReport struct {
	ByType   []repository.MaintenanceStats       `json:"byType"`
	ByStatus []repository.MaintenanceStatusStats `json:"byStatus"`
}

func (s *ReportService) GetMaintenanceStats(start, end time.Time) (*MaintenanceReport, error) {
	byType, err := s.maintenanceRepo.GetStatsByType(start, end)
	if err != nil {
		return nil, err
	}

	byStatus, err := s.maintenanceRepo.GetStatsByStatus(start, end)
	if err != nil {
		return nil, err
	}

	return &MaintenanceReport{
		ByType:   byType,
		ByStatus: byStatus,
	}, nil
}

func (s *ReportService) GetTenantRanking(limit int, start, end time.Time) ([]repository.TenantFeeRanking, error) {
	return s.feeRepo.GetTenantRanking(limit, start, end)
}

type DashboardData struct {
	TotalTenants      int64   `json:"totalTenants"`
	TotalRooms        int64   `json:"totalRooms"`
	OccupiedRooms     int64   `json:"occupiedRooms"`
	OccupancyRate     float64 `json:"occupancyRate"`
	ActiveContracts   int64   `json:"activeContracts"`
	PendingFees       int64   `json:"pendingFees"`
	UnpaidAmount      float64 `json:"unpaidAmount"`
	PendingMaintenance int64  `json:"pendingMaintenance"`
}

func (s *ReportService) GetDashboardData() (*DashboardData, error) {
	tenants, err := s.tenantRepo.FindAll()
	if err != nil {
		return nil, err
	}

	totalRooms, err := s.roomRepo.CountTotal()
	if err != nil {
		return nil, err
	}

	occupiedRooms, err := s.roomRepo.CountByStatus("occupied")
	if err != nil {
		return nil, err
	}

	activeContracts, err := s.contractRepo.CountByStatus("active")
	if err != nil {
		return nil, err
	}

	pendingFees, err := s.feeRepo.CountByStatus("unpaid")
	if err != nil {
		return nil, err
	}

	unpaidAmount, err := s.feeRepo.SumUnpaidAmount()
	if err != nil {
		return nil, err
	}

	pendingMaintenance, err := s.maintenanceRepo.CountByStatus("pending")
	if err != nil {
		return nil, err
	}

	var occupancyRate float64
	if totalRooms > 0 {
		occupancyRate = float64(occupiedRooms) / float64(totalRooms) * 100
	}

	return &DashboardData{
		TotalTenants:       int64(len(tenants)),
		TotalRooms:         totalRooms,
		OccupiedRooms:      occupiedRooms,
		OccupancyRate:      occupancyRate,
		ActiveContracts:    activeContracts,
		PendingFees:        pendingFees,
		UnpaidAmount:       unpaidAmount,
		PendingMaintenance: pendingMaintenance,
	}, nil
}
