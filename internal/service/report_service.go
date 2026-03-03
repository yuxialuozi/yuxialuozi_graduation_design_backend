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
	Total  float64                     `json:"total"`
	ByDay  []repository.IncomeByDay    `json:"byDay"`
	ByType []repository.FeeComposition `json:"byType"`
}

func (s *ReportService) GetIncomeReport(start, end time.Time, groupBy string) (*IncomeReport, error) {
	// 获取实际支付的费用
	actualTotal, err := s.feeRepo.SumByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	// 获取合同预期收入
	contractTotal, err := s.contractRepo.GetTotalIncomeByPeriod(start, end)
	if err != nil {
		return nil, err
	}

	// 总收入 = 实际支付 + 合同预期收入
	total := actualTotal + contractTotal

	// 获取按天统计的实际支付收入
	actualByDay, err := s.feeRepo.GetIncomeByDay(start, end)
	if err != nil {
		return nil, err
	}

	// 获取按天统计的合同预期收入
	contractByDay, err := s.contractRepo.GetIncomeByDay(start, end)
	if err != nil {
		return nil, err
	}

	// 合并按天收入：实际支付 + 合同预期
	combinedByDay := make([]repository.IncomeByDay, 0)
	dayMap := make(map[string]float64)

	// 添加实际支付收入
	for _, item := range actualByDay {
		dayMap[item.Day] += item.Amount
	}

	// 添加合同预期收入
	for _, item := range contractByDay {
		dayMap[item.Day] += item.Amount
	}

	// 转换为切片
	for day, amount := range dayMap {
		combinedByDay = append(combinedByDay, repository.IncomeByDay{
			Day:    day,
			Amount: amount,
		})
	}

	// 按日期排序
	for i := 0; i < len(combinedByDay); i++ {
		for j := i + 1; j < len(combinedByDay); j++ {
			if combinedByDay[i].Day > combinedByDay[j].Day {
				combinedByDay[i], combinedByDay[j] = combinedByDay[j], combinedByDay[i]
			}
		}
	}

	// 获取费用类型构成（仅实际支付）
	byType, err := s.feeRepo.GetComposition(start, end)
	if err != nil {
		return nil, err
	}

	return &IncomeReport{
		Total:  total,
		ByDay:  combinedByDay,
		ByType: byType,
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
	// 费用构成基于实际支付的费用
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
	// 租户排名基于实际支付的费用
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

	// 生成收入图表数据（2025年全年，按月）
	incomeChart := s.UpdateDashboardIncomeChart()

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
