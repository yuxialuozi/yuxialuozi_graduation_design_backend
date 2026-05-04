package service

import (
	"errors"
	"time"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

// TenantPortalService 租户端服务
type TenantPortalService struct {
	tenantRepo      *repository.TenantRepository
	contractRepo    *repository.ContractRepository
	roomRepo        *repository.RoomRepository
	feeRepo         *repository.FeeRepository
	maintenanceRepo *repository.MaintenanceRepository
	maintenanceService *MaintenanceService
	userRepo        *repository.UserRepository
}

func NewTenantPortalService(
	tenantRepo *repository.TenantRepository,
	contractRepo *repository.ContractRepository,
	roomRepo *repository.RoomRepository,
	feeRepo *repository.FeeRepository,
	maintenanceRepo *repository.MaintenanceRepository,
	maintenanceService *MaintenanceService,
	userRepo *repository.UserRepository,
) *TenantPortalService {
	return &TenantPortalService{
		tenantRepo:         tenantRepo,
		contractRepo:       contractRepo,
		roomRepo:           roomRepo,
		feeRepo:            feeRepo,
		maintenanceRepo:    maintenanceRepo,
		maintenanceService: maintenanceService,
		userRepo:           userRepo,
	}
}

// TenantProfileResponse 租户个人信息响应
type TenantProfileResponse struct {
	UserID       uint      `json:"userId"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	TenantID     uint      `json:"tenantId"`
	TenantName   string    `json:"tenantName"`
	Contact      string    `json:"contact"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	ActiveRoom   int       `json:"activeRoom"`    // 活跃房间数
	ActiveContract int     `json:"activeContract"` // 有效合同数
	UnpaidFee    float64   `json:"unpaidFee"`      // 未缴费用
	PendingMaintenance int `json:"pendingMaintenance"` // 待处理维修
}

// TenantDashboardResponse 租户仪表盘响应
type TenantDashboardResponse struct {
	TotalContract   int                        `json:"totalContract"`    // 合同总数
	ActiveContract  int                        `json:"activeContract"`   // 有效合同数
	TotalRoom       int                        `json:"totalRoom"`        // 房间总数
	TotalFee        float64                    `json:"totalFee"`         // 累计费用
	UnpaidFee       float64                    `json:"unpaidFee"`        // 未缴费用
	PaidFee         float64                    `json:"paidFee"`          // 已缴费用
	TotalMaintenance int                       `json:"totalMaintenance"` // 维修工单总数
	PendingMaintenance int                     `json:"pendingMaintenance"` // 待处理维修
	FeeTrend        []FeeTrendItem             `json:"feeTrend"`         // 费用趋势
	RecentFees      []model.Fee                `json:"recentFees"`       // 最近账单
	ContractList    []model.Contract           `json:"contractList"`     // 合同列表
}

// FeeTrendItem 费用趋势项
type FeeTrendItem struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
}

// GetProfile 获取租户个人信息
func (s *TenantPortalService) GetProfile(userID, tenantID uint) (*TenantProfileResponse, error) {
	// 获取用户信息
	var username, nickname string
	user, err := s.userRepo.FindByID(userID)
	if err == nil {
		username = user.Username
		nickname = user.Nickname
	}

	if tenantID == 0 {
		return &TenantProfileResponse{
			UserID:   userID,
			Username: username,
			Nickname: nickname,
		}, nil
	}

	tenant, err := s.tenantRepo.FindByID(tenantID)
	if err != nil {
		return nil, err
	}

	// 统计活跃房间数
	rooms, err := s.roomRepo.FindByTenantID(tenantID)
	if err != nil {
		rooms = []model.Room{}
	}

	// 统计有效合同数
	contracts, err := s.contractRepo.FindByTenantID(tenantID)
	if err != nil {
		contracts = []model.Contract{}
	}
	activeContract := 0
	for _, c := range contracts {
		if c.Status == "active" {
			activeContract++
		}
	}

	// 统计未缴费用
	var unpaidAmount float64
	fees, _, _ := s.feeRepo.List(1, 1000, tenantID, "", "", "unpaid", "")
	for _, f := range fees {
		unpaidAmount += f.Amount
	}

	// 统计待处理维修（只统计当前租户的）
	pendingMaintenance := int64(0)
	maintenances, _, _ := s.maintenanceRepo.List(1, 1000, "", "", "pending", "", tenantID)
	pendingMaintenance = int64(len(maintenances))

	return &TenantProfileResponse{
		UserID:            userID,
		Username:          username,
		Nickname:          nickname,
		TenantID:          tenant.ID,
		TenantName:        tenant.Name,
		Contact:           tenant.ContactPerson,
		Phone:             tenant.Phone,
		Email:             tenant.Email,
		ActiveRoom:        len(rooms),
		ActiveContract:    activeContract,
		UnpaidFee:         unpaidAmount,
		PendingMaintenance: int(pendingMaintenance),
	}, nil
}

// GetContracts 获取租户合同列表
func (s *TenantPortalService) GetContracts(tenantID uint, page, pageSize int) (*dto.PageResult, error) {
	contracts, total, err := s.contractRepo.List(page, pageSize, "", "", nil, nil, tenantID)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     contracts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetRooms 获取租户房间列表
func (s *TenantPortalService) GetRooms(tenantID uint, page, pageSize int) (*dto.PageResult, error) {
	rooms, total, err := s.roomRepo.List(page, pageSize, "", "", "", tenantID)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     rooms,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetFees 获取租户账单列表
func (s *TenantPortalService) GetFees(tenantID uint, page, pageSize int, status, feeType, period string) (*dto.PageResult, error) {
	fees, total, err := s.feeRepo.List(page, pageSize, tenantID, "", feeType, status, period)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     fees,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// PayFee 缴纳账单
func (s *TenantPortalService) PayFee(feeID, tenantID uint, paidDate *time.Time) error {
	fee, err := s.feeRepo.FindByID(feeID)
	if err != nil {
		return err
	}

	// 验证账单属于当前租户
	if fee.TenantID != tenantID {
		return errors.New("无权操作此账单")
	}

	fee.PaidDate = &model.CustomTime{Time: *paidDate}
	fee.Status = "paid"
	return s.feeRepo.Update(fee)
}

// GetMaintenance 获取租户维修工单列表
func (s *TenantPortalService) GetMaintenance(tenantID uint, page, pageSize int, status string) (*dto.PageResult, error) {
	// 参数顺序: page, pageSize, keyword, maintenanceType, status, priority, tenantID
	maintenances, total, err := s.maintenanceRepo.List(page, pageSize, "", "", status, "", tenantID)
	if err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:     maintenances,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateMaintenance 创建维修工单
func (s *TenantPortalService) CreateMaintenance(tenantID uint, username string, req *dto.CreateMaintenanceRequest) error {
	// 获取租户信息
	tenant, err := s.tenantRepo.FindByID(tenantID)
	if err != nil {
		return err
	}

	// 获取租户的房间
	rooms, _ := s.roomRepo.FindByTenantID(tenantID)
	roomNo := ""
	if len(rooms) > 0 {
		roomNo = rooms[0].RoomNo
	}

	maintenance := &model.Maintenance{
		TenantID:    tenantID,
		Type:        req.Type,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      "pending",
		RoomNo:      roomNo,
		TenantName:  tenant.Name,
	}

	// 使用 MaintenanceService 来创建（会自动生成工单号）
	_, err = s.maintenanceService.Create(maintenance)
	return err
}

// GetDashboard 获取租户仪表盘数据
func (s *TenantPortalService) GetDashboard(tenantID uint, username string) (*TenantDashboardResponse, error) {
	// 获取合同列表
	contracts, _ := s.contractRepo.FindByTenantID(tenantID)
	activeContract := 0
	for _, c := range contracts {
		if c.Status == "active" {
			activeContract++
		}
	}

	// 获取房间列表
	rooms, _ := s.roomRepo.FindByTenantID(tenantID)

	// 获取账单列表（当前租户）
	fees, _, _ := s.feeRepo.List(1, 1000, tenantID, "", "", "", "")
	var totalFee, unpaidFee, paidFee float64
	recentFees := []model.Fee{}
	for i, f := range fees {
		totalFee += f.Amount
		if f.Status == "paid" {
			paidFee += f.Amount
		} else {
			unpaidFee += f.Amount
		}
		if i < 5 {
			recentFees = append(recentFees, f)
		}
	}

	// 统计维修工单（只统计当前租户的）
	allMaintenances, _, _ := s.maintenanceRepo.List(1, 10000, "", "", "", "", tenantID)
	tenantPending := int64(0)
	tenantProcessing := int64(0)
	for _, m := range allMaintenances {
		if m.Status == "pending" {
			tenantPending++
		} else if m.Status == "processing" {
			tenantProcessing++
		}
	}

	// 计算费用趋势（最近6个月，只统计当前租户的）
	now := time.Now()
	var feeTrend []FeeTrendItem
	for i := 5; i >= 0; i-- {
		month := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.UTC)
		monthEnd := month.AddDate(0, 1, -1)

		// 计算该月当前租户已缴费用
		monthPaid := float64(0)
		for _, f := range fees {
			if f.Status == "paid" && f.PaidDate != nil {
				paidTime := f.PaidDate.Time
				if !paidTime.Before(month) && !paidTime.After(monthEnd) {
					monthPaid += f.Amount
				}
			}
		}

		feeTrend = append(feeTrend, FeeTrendItem{
			Month:  month.Format("2006-01"),
			Amount: monthPaid,
		})
	}

	return &TenantDashboardResponse{
		TotalContract:     len(contracts),
		ActiveContract:    activeContract,
		TotalRoom:         len(rooms),
		TotalFee:          totalFee,
		UnpaidFee:         unpaidFee,
		PaidFee:           paidFee,
		TotalMaintenance:  int(tenantPending) + int(tenantProcessing),
		PendingMaintenance: int(tenantPending),
		FeeTrend:          feeTrend,
		RecentFees:        recentFees,
		ContractList:      contracts,
	}, nil
}