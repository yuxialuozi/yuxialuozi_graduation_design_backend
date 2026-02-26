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
	Total   float64                     `json:"total"`
	ByMonth []repository.IncomeByMonth  `json:"byMonth"`
	ByType  []repository.FeeComposition `json:"byType"`
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
	TotalTenants       int64   `json:"totalTenants"`
	TotalRooms         int64   `json:"totalRooms"`
	OccupiedRooms      int64   `json:"occupiedRooms"`
	OccupancyRate      float64 `json:"occupancyRate"`
	ActiveContracts    int64   `json:"activeContracts"`
	PendingFees        int64   `json:"pendingFees"`
	UnpaidAmount       float64 `json:"unpaidAmount"`
	PendingMaintenance int64   `json:"pendingMaintenance"`
	// 图表数据，匹配前端期望
	IncomeChart            []IncomeChartData            `json:"incomeChart"`
	MaintenanceStatusChart []MaintenanceStatusChartData `json:"maintenanceStatusChart"`
	FeeTypeChart           []FeeTypeChartData           `json:"feeTypeChart"`
}

// 添加图表数据结构类型
type IncomeChartData struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type MaintenanceStatusChartData struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type FeeTypeChartData struct {
	FeeType string  `json:"feeType"`
	Amount  float64 `json:"amount"`
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

	// 生成收入图表数据（最近6个月）
	incomeChart := make([]IncomeChartData, 0)
	now := time.Now()
	for i := 5; i >= 0; i-- {
		// 使用 AddDate 方法减去月份
		startDate := now.AddDate(0, -i, 0)
		// 将日期设置为当月第一天
		startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())

		startOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		monthlyIncome, _ := s.feeRepo.SumByPeriod(startOfMonth, endOfMonth)
		incomeChart = append(incomeChart, IncomeChartData{
			Date:   startOfMonth.Format("2006-01"),
			Amount: monthlyIncome,
		})
	}

	// 生成维修状态图表数据
	maintenanceStatusChart := make([]MaintenanceStatusChartData, 0)
	statuses := []string{"pending", "processing", "completed", "cancelled"}
	for _, status := range statuses {
		count, err := s.maintenanceRepo.CountByStatus(status)
		if err != nil {
			count = 0
		}
		maintenanceStatusChart = append(maintenanceStatusChart, MaintenanceStatusChartData{
			Status: status,
			Count:  count,
		})
	}

	// 生成费用类型图表数据
	feeTypeChart := make([]FeeTypeChartData, 0)
	feeTypes := []string{"rent", "water", "electricity", "property", "other"}
	for _, feeType := range feeTypes {
		amount, err := s.feeRepo.SumByType(feeType)
		if err != nil {
			amount = 0
		}
		feeTypeChart = append(feeTypeChart, FeeTypeChartData{
			FeeType: feeType,
			Amount:  amount,
		})
	}

	return &DashboardData{
		TotalTenants:           int64(len(tenants)),
		TotalRooms:             totalRooms,
		OccupiedRooms:          occupiedRooms,
		OccupancyRate:          occupancyRate,
		ActiveContracts:        activeContracts,
		PendingFees:            pendingFees,
		UnpaidAmount:           unpaidAmount,
		PendingMaintenance:     pendingMaintenance,
		IncomeChart:            incomeChart,
		MaintenanceStatusChart: maintenanceStatusChart,
		FeeTypeChart:           feeTypeChart,
	}, nil
}
