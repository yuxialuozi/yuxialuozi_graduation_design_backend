package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"yuxialuozi_graduation_design_backend/internal/knowledge"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

// ToolExecutor executes MCP tool calls
type ToolExecutor struct {
	db        *gorm.DB
	tenantRepo *repository.TenantRepository
}

// NewToolExecutor creates a new tool executor
func NewToolExecutor(db *gorm.DB) *ToolExecutor {
	return &ToolExecutor{
		db:        db,
		tenantRepo: repository.NewTenantRepository(db),
	}
}

// GetTenantProfile returns the tenant's profile information
func (e *ToolExecutor) GetTenantProfile(args map[string]interface{}) (string, error) {
	tenantID, ok := args["tenantId"].(float64)
	if !ok || tenantID == 0 {
		return `{"error": "无法获取租户信息，请确认已登录"}`, nil
	}

	tenant, err := e.tenantRepo.FindByID(uint(tenantID))
	if err != nil {
		return "", fmt.Errorf("查询租户信息失败: %w", err)
	}

	result := map[string]interface{}{
		"租户ID":    tenant.ID,
		"租户名称":  tenant.Name,
		"联系人":    tenant.ContactPerson,
		"联系电话":  tenant.Phone,
		"邮箱":      tenant.Email,
		"状态":      tenant.Status,
		"创建时间":  tenant.CreatedAt.Format("2006-01-02"),
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetTenantFees returns the tenant's fee/billing information
func (e *ToolExecutor) GetTenantFees(args map[string]interface{}) (string, error) {
	var fees []model.Fee
	query := e.db.Model(&model.Fee{})

	// Filter by status if provided
	if status, ok := args["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by tenant ID (in real implementation, get from JWT)
	if tenantID, ok := args["tenantId"].(float64); ok {
		query = query.Where("tenant_id = ?", uint(tenantID))
	}

	// Limit results
	pageSize := 10
	if ps, ok := args["pageSize"].(float64); ok {
		pageSize = int(ps)
		if pageSize > 50 {
			pageSize = 50
		}
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Find(&fees).Error; err != nil {
		return "", fmt.Errorf("查询费用失败: %w", err)
	}

	// Format fees for display
	var result []map[string]interface{}
	for _, fee := range fees {
		statusText := map[string]string{
			"paid":   "已缴",
			"unpaid": "未缴",
			"overdue": "逾期",
		}
		feeTypeText := map[string]string{
			"rent":      "租金",
			"water":     "水费",
			"electricity": "电费",
			"property":  "物业费",
			"other":     "其他",
		}

		result = append(result, map[string]interface{}{
			"房间号":    fee.RoomNo,
			"费用类型":  feeTypeText[fee.FeeType],
			"金额":     fmt.Sprintf("%.2f元", fee.Amount),
			"周期":     fee.Period,
			"应缴日期":  fee.DueDate.Format("2006-01-02"),
			"状态":     statusText[fee.Status],
		})
	}

	if len(result) == 0 {
		return `{"message": "暂无费用记录"}`, nil
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetTenantContracts returns the tenant's contract information
func (e *ToolExecutor) GetTenantContracts(args map[string]interface{}) (string, error) {
	var contracts []model.Contract
	query := e.db.Model(&model.Contract{})

	if status, ok := args["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if tenantID, ok := args["tenantId"].(float64); ok {
		query = query.Where("tenant_id = ?", uint(tenantID))
	}

	if err := query.Order("created_at DESC").Limit(10).Find(&contracts).Error; err != nil {
		return "", fmt.Errorf("查询合同失败: %w", err)
	}

	var result []map[string]interface{}
	for _, c := range contracts {
		statusText := map[string]string{
			"draft":      "草稿",
			"active":     "生效中",
			"expired":    "已过期",
			"terminated": "已终止",
		}

		result = append(result, map[string]interface{}{
			"合同编号":    c.ContractNo,
			"开始日期":    c.StartDate.Format("2006-01-02"),
			"结束日期":    c.EndDate.Format("2006-01-02"),
			"租金金额":    fmt.Sprintf("%.2f元/月", c.Amount),
			"状态":       statusText[c.Status],
		})
	}

	if len(result) == 0 {
		return `{"message": "暂无合同记录"}`, nil
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetTenantMaintenance returns the tenant's maintenance requests
func (e *ToolExecutor) GetTenantMaintenance(args map[string]interface{}) (string, error) {
	var maintenances []model.Maintenance
	query := e.db.Model(&model.Maintenance{})

	if status, ok := args["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	if tenantID, ok := args["tenantId"].(float64); ok {
		query = query.Where("tenant_id = ?", uint(tenantID))
	}

	if err := query.Order("created_at DESC").Limit(10).Find(&maintenances).Error; err != nil {
		return "", fmt.Errorf("查询维修工单失败: %w", err)
	}

	var result []map[string]interface{}
	for _, m := range maintenances {
		statusText := map[string]string{
			"pending":    "待处理",
			"processing": "处理中",
			"completed":  "已完成",
			"cancelled":  "已取消",
		}
		typeText := map[string]string{
			"electrical": "电气",
			"plumbing":   "水管",
			"appliance":  "家电",
			"furniture":  "家具",
			"other":      "其他",
		}
		priorityText := map[string]string{
			"low":    "低",
			"medium": "中",
			"high":   "高",
			"urgent": "紧急",
		}

		completedAt := ""
		if !m.CompletedAt.IsZero() {
			completedAt = m.CompletedAt.Format("2006-01-02 15:04")
		}

		result = append(result, map[string]interface{}{
			"工单编号":    m.TicketNo,
			"房间号":     m.RoomNo,
			"维修类型":    typeText[m.Type],
			"问题描述":    m.Description,
			"优先级":     priorityText[m.Priority],
			"状态":      statusText[m.Status],
			"创建时间":    m.CreatedAt.Format("2006-01-02 15:04"),
			"完成时间":    completedAt,
			"负责人":     m.Assignee,
		})
	}

	if len(result) == 0 {
		return `{"message": "暂无维修工单记录"}`, nil
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetTenantDashboard returns the tenant's dashboard data
func (e *ToolExecutor) GetTenantDashboard(args map[string]interface{}) (string, error) {
	var result struct {
		统计信息 struct {
			有效合同数        int     `json:"有效合同数"`
			租住房间数        int     `json:"租住房间数"`
			累计费用          float64 `json:"累计费用(元)"`
			已缴费用          float64 `json:"已缴费用(元)"`
			未缴费用          float64 `json:"未缴费用(元)"`
			总维修工单数       int     `json:"总维修工单数"`
			待处理维修工单数    int     `json:"待处理维修工单数"`
		} `json:"统计信息"`
	}

	// Get counts from database
	var activeContracts, totalRooms, totalFees, paidFees, unpaidFees int64
	var totalMaintenance, pendingMaintenance int64

	if tenantID, ok := args["tenantId"].(float64); ok {
		tid := uint(tenantID)
		e.db.Model(&model.Contract{}).Where("tenant_id = ? AND status = ?", tid, "active").Count(&activeContracts)
		e.db.Model(&model.Room{}).Where("tenant_id = ?", tid).Count(&totalRooms)
		e.db.Model(&model.Fee{}).Where("tenant_id = ?", tid).Select("COALESCE(SUM(amount), 0)").Scan(&totalFees)
		e.db.Model(&model.Fee{}).Where("tenant_id = ? AND status = ?", tid, "paid").Select("COALESCE(SUM(amount), 0)").Scan(&paidFees)
		e.db.Model(&model.Fee{}).Where("tenant_id = ? AND status IN (?, ?)", tid, "unpaid", "overdue").Select("COALESCE(SUM(amount), 0)").Scan(&unpaidFees)
		e.db.Model(&model.Maintenance{}).Where("tenant_id = ?", tid).Count(&totalMaintenance)
		e.db.Model(&model.Maintenance{}).Where("tenant_id = ? AND status = ?", tid, "pending").Count(&pendingMaintenance)
	} else {
		// Demo data for testing
		result.统计信息 = struct {
			有效合同数        int     `json:"有效合同数"`
			租住房间数        int     `json:"租住房间数"`
			累计费用          float64 `json:"累计费用(元)"`
			已缴费用          float64 `json:"已缴费用(元)"`
			未缴费用          float64 `json:"未缴费用(元)"`
			总维修工单数       int     `json:"总维修工单数"`
			待处理维修工单数    int     `json:"待处理维修工单数"`
		}{
			有效合同数:        1,
			租住房间数:        1,
			累计费用:          12500.00,
			已缴费用:          10000.00,
			未缴费用:          2500.00,
			总维修工单数:       3,
			待处理维修工单数:    1,
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		return string(data), nil
	}

	result.统计信息 = struct {
		有效合同数        int     `json:"有效合同数"`
		租住房间数        int     `json:"租住房间数"`
		累计费用          float64 `json:"累计费用(元)"`
		已缴费用          float64 `json:"已缴费用(元)"`
		未缴费用          float64 `json:"未缴费用(元)"`
		总维修工单数       int     `json:"总维修工单数"`
		待处理维修工单数    int     `json:"待处理维修工单数"`
	}{
		有效合同数:        int(activeContracts),
		租住房间数:        int(totalRooms),
		累计费用:          float64(totalFees),
		已缴费用:          float64(paidFees),
		未缴费用:          float64(unpaidFees),
		总维修工单数:       int(totalMaintenance),
		待处理维修工单数:    int(pendingMaintenance),
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetTenantRooms returns the tenant's room information
func (e *ToolExecutor) GetTenantRooms(args map[string]interface{}) (string, error) {
	var rooms []model.Room
	query := e.db.Model(&model.Room{})

	if tenantID, ok := args["tenantId"].(float64); ok {
		query = query.Where("tenant_id = ?", uint(tenantID))
	}

	if err := query.Find(&rooms).Error; err != nil {
		return "", fmt.Errorf("查询房间失败: %w", err)
	}

	var result []map[string]interface{}
	for _, r := range rooms {
		statusText := map[string]string{
			"vacant":     "空置",
			"occupied":   "已入住",
			"maintenance": "维护中",
		}

		result = append(result, map[string]interface{}{
			"房间号":    r.RoomNo,
			"楼栋":     r.Building,
			"楼层":     r.Floor,
			"面积":     fmt.Sprintf("%.2f平方米", r.Area),
			"月租金":   fmt.Sprintf("%.2f元", r.MonthlyRent),
			"状态":    statusText[r.Status],
		})
	}

	if len(result) == 0 {
		return `{"message": "暂无房间记录"}`, nil
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// SearchKnowledge searches the knowledge base for relevant information
func (e *ToolExecutor) SearchKnowledge(args map[string]interface{}) (string, error) {
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return `{"error": "请提供搜索关键词"}`, nil
	}

	results := knowledge.Search(query, 5)
	if len(results) == 0 {
		return `{"message": "未找到相关知识，请尝试其他关键词"}`, nil
	}

	var result []map[string]interface{}
	for _, r := range results {
		result = append(result, map[string]interface{}{
			"标题":     r.Item.Title,
			"分类":     r.Item.Category,
			"内容":     r.Item.Content,
			"相关度":   fmt.Sprintf("%.2f", r.Score),
		})
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// ExecuteTool executes a tool by name with arguments
func (e *ToolExecutor) ExecuteTool(name string, args map[string]interface{}) (string, error) {
	tools := GetAllTools(e)
	for _, tool := range tools {
		if tool.Name == name {
			return tool.Handler(args)
		}
	}
	return "", fmt.Errorf("未找到工具: %s", name)
}

// FormatFeesForDisplay formats fee data for AI response
func FormatFeesForDisplay(fees []model.Fee) string {
	if len(fees) == 0 {
		return "您目前没有任何费用记录。"
	}

	statusText := map[string]string{
		"paid":   "已缴",
		"unpaid": "未缴",
		"overdue": "逾期",
	}
	feeTypeText := map[string]string{
		"rent":      "租金",
		"water":     "水费",
		"electricity": "电费",
		"property":  "物业费",
		"other":     "其他",
	}

	var result string
	result = fmt.Sprintf("您共有 %d 条费用记录：\n\n", len(fees))

	totalUnpaid := 0.0
	for i, fee := range fees {
		seq := i + 1
		status := statusText[fee.Status]
		if fee.Status == "unpaid" || fee.Status == "overdue" {
			totalUnpaid += fee.Amount
			status = "⚠️ " + status
		}
		result += fmt.Sprintf("%d. %s - %s - %s - %s (%s)\n",
			seq,
			fee.RoomNo,
			feeTypeText[fee.FeeType],
			fmt.Sprintf("%.2f元", fee.Amount),
			fee.Period,
			status,
		)
	}

	if totalUnpaid > 0 {
		result += fmt.Sprintf("\n⚠️ 您有 %.2f 元费用待缴纳，请及时处理！", totalUnpaid)
	}

	return result
}

// FormatContractsForDisplay formats contract data for AI response
func FormatContractsForDisplay(contracts []model.Contract) string {
	if len(contracts) == 0 {
		return "您目前没有任何合同记录。"
	}

	statusText := map[string]string{
		"draft":      "草稿",
		"active":     "生效中",
		"expired":    "已过期",
		"terminated": "已终止",
	}

	var result string
	result = fmt.Sprintf("您共有 %d 份合同：\n\n", len(contracts))

	for i, c := range contracts {
		seq := i + 1
		status := statusText[c.Status]
		if c.Status == "expired" {
			status = "⚠️ " + status
		}

		now := time.Now()
		daysLeft := int(c.EndDate.Time.Sub(now).Hours() / 24)
		daysText := ""
		if daysLeft > 0 && daysLeft <= 30 {
			daysText = fmt.Sprintf(" (还剩 %d 天到期)", daysLeft)
		}

		result += fmt.Sprintf("%d. 合同号: %s\n", seq, c.ContractNo)
		result += fmt.Sprintf("   租期: %s 至 %s%s\n", c.StartDate.Format("2006-01-02"), c.EndDate.Format("2006-01-02"), daysText)
		result += fmt.Sprintf("   租金: %.2f 元/月\n", c.Amount)
		result += fmt.Sprintf("   状态: %s\n\n", status)
	}

	return result
}

// FormatMaintenanceForDisplay formats maintenance data for AI response
func FormatMaintenanceForDisplay(maintenances []model.Maintenance) string {
	if len(maintenances) == 0 {
		return "您目前没有任何维修工单记录。"
	}

	statusText := map[string]string{
		"pending":    "待处理",
		"processing": "处理中",
		"completed":  "已完成",
		"cancelled":  "已取消",
	}
	typeText := map[string]string{
		"electrical": "电气",
		"plumbing":   "水管",
		"appliance":  "家电",
		"furniture":  "家具",
		"other":      "其他",
	}
	priorityText := map[string]string{
		"low":    "低",
		"medium": "中",
		"high":   "高",
		"urgent": "紧急",
	}

	var result string
	pending := 0
	for _, m := range maintenances {
		if m.Status == "pending" || m.Status == "processing" {
			pending++
		}
	}

	result = fmt.Sprintf("您共有 %d 条维修工单，其中 %d 条待处理/处理中：\n\n", len(maintenances), pending)

	for i, m := range maintenances {
		seq := i + 1
		status := statusText[m.Status]
		if m.Status == "pending" {
			status = "⚠️ " + status
		}

		result += fmt.Sprintf("%d. 工单号: %s\n", seq, m.TicketNo)
		result += fmt.Sprintf("   房间: %s | 类型: %s | 优先级: %s\n", m.RoomNo, typeText[m.Type], priorityText[m.Priority])
		result += fmt.Sprintf("   问题: %s\n", m.Description)
		result += fmt.Sprintf("   状态: %s | 提交时间: %s\n", status, m.CreatedAt.Format("2006-01-02 15:04"))
		if m.Assignee != "" {
			result += fmt.Sprintf("   负责人: %s\n", m.Assignee)
		}
		result += "\n"
	}

	return result
}

// GetAdminDashboard returns system-wide dashboard data for admin
func (e *ToolExecutor) GetAdminDashboard(args map[string]interface{}) (string, error) {
	var totalTenants, totalRooms, occupiedRooms int64
	var activeContracts, pendingFees int64
	var unpaidAmount float64
	var pendingMaintenance int64

	e.db.Model(&model.Tenant{}).Count(&totalTenants)
	e.db.Model(&model.Room{}).Count(&totalRooms)
	e.db.Model(&model.Room{}).Where("status = ?", "occupied").Count(&occupiedRooms)
	e.db.Model(&model.Contract{}).Where("status = ?", "active").Count(&activeContracts)
	e.db.Model(&model.Fee{}).Where("status IN (?, ?)", "unpaid", "overdue").Count(&pendingFees)
	e.db.Model(&model.Fee{}).Where("status IN (?, ?)", "unpaid", "overdue").Select("COALESCE(SUM(amount), 0)").Scan(&unpaidAmount)
	e.db.Model(&model.Maintenance{}).Where("status = ?", "pending").Count(&pendingMaintenance)

	occupancyRate := 0.0
	if totalRooms > 0 {
		occupancyRate = float64(occupiedRooms) / float64(totalRooms) * 100
	}

	type monthlyIncome struct {
		Month  string  `json:"month"`
		Amount float64 `json:"amount"`
	}
	var income []monthlyIncome
	e.db.Raw(`
		SELECT TO_CHAR(due_date, 'YYYY-MM') as month, SUM(amount) as amount
		FROM fees
		WHERE status = 'paid' AND due_date >= NOW() - INTERVAL '6 months'
		GROUP BY TO_CHAR(due_date, 'YYYY-MM')
		ORDER BY month DESC LIMIT 6
	`).Scan(&income)

	type feeByType struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount"`
	}
	var feesByType []feeByType
	e.db.Raw(`
		SELECT fee_type as type, SUM(amount) as amount
		FROM fees GROUP BY fee_type
	`).Scan(&feesByType)

	feeTypeText := map[string]string{"rent": "租金", "water": "水费", "electricity": "电费", "property": "物业费", "other": "其他"}
	var formattedFees []map[string]interface{}
	for _, ft := range feesByType {
		formattedFees = append(formattedFees, map[string]interface{}{
			"类型":   feeTypeText[ft.Type],
			"金额":   fmt.Sprintf("%.2f元", ft.Amount),
		})
	}

	result := map[string]interface{}{
		"系统概览": map[string]interface{}{
			"总租户数":     totalTenants,
			"总房间数":     totalRooms,
			"已入住房间数":  occupiedRooms,
			"入住率":       fmt.Sprintf("%.1f%%", occupancyRate),
			"有效合同数":   activeContracts,
			"待缴费账单数":  pendingFees,
			"未缴总金额":   fmt.Sprintf("%.2f元", unpaidAmount),
			"待处理维修数":  pendingMaintenance,
		},
		"月收入趋势(近6月)": income,
		"费用构成":          formattedFees,
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetAdminFees returns system-wide fee statistics for admin
func (e *ToolExecutor) GetAdminFees(args map[string]interface{}) (string, error) {
	statusFilter := ""
	if status, ok := args["status"].(string); ok && status != "" {
		statusFilter = status
	}

	query := e.db.Model(&model.Fee{})
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var fees []model.Fee
	if err := query.Order("created_at DESC").Limit(20).Find(&fees).Error; err != nil {
		return "", fmt.Errorf("查询费用失败: %w", err)
	}

	var total, paid, unpaid int64
	var unpaidAmt float64
	e.db.Model(&model.Fee{}).Count(&total)
	e.db.Model(&model.Fee{}).Where("status = ?", "paid").Count(&paid)
	e.db.Model(&model.Fee{}).Where("status IN (?, ?)", "unpaid", "overdue").Count(&unpaid)
	e.db.Model(&model.Fee{}).Where("status IN (?, ?)", "unpaid", "overdue").Select("COALESCE(SUM(amount), 0)").Scan(&unpaidAmt)

	statusText := map[string]string{"paid": "已缴", "unpaid": "未缴", "overdue": "逾期"}
	feeTypeText := map[string]string{"rent": "租金", "water": "水费", "electricity": "电费", "property": "物业费", "other": "其他"}

	var records []map[string]interface{}
	records = append(records, map[string]interface{}{
		"统计": map[string]interface{}{
			"总账单数": total,
			"已缴":    paid,
			"未缴":    unpaid,
			"未缴总金额": fmt.Sprintf("%.2f元", unpaidAmt),
		},
	})

	for _, fee := range fees {
		tenantName := ""
		if fee.TenantID > 0 {
			var t model.Tenant
			if err := e.db.First(&t, fee.TenantID).Error; err == nil {
				tenantName = t.Name
			}
		}
		records = append(records, map[string]interface{}{
			"租户":    tenantName,
			"房间号":  fee.RoomNo,
			"类型":    feeTypeText[fee.FeeType],
			"金额":    fmt.Sprintf("%.2f元", fee.Amount),
			"周期":    fee.Period,
			"应缴日期": fee.DueDate.Format("2006-01-02"),
			"状态":    statusText[fee.Status],
		})
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	return string(data), nil
}

// GetAdminTenants returns system-wide tenant statistics for admin
func (e *ToolExecutor) GetAdminTenants(args map[string]interface{}) (string, error) {
	var totalTenants, activeTenants, inactiveTenants int64
	e.db.Model(&model.Tenant{}).Count(&totalTenants)
	e.db.Model(&model.Tenant{}).Where("status = ?", "active").Count(&activeTenants)
	e.db.Model(&model.Tenant{}).Where("status != ?", "active").Count(&inactiveTenants)

	type tenantStat struct {
		Name   string `json:"name"`
		Rooms  int    `json:"rooms"`
		Active int    `json:"activeContracts"`
	}
	var stats []tenantStat
	e.db.Raw(`
		SELECT t.name, COALESCE(COUNT(DISTINCT r.id), 0) as rooms, COALESCE(COUNT(DISTINCT c.id), 0) as active_contracts
		FROM tenants t
		LEFT JOIN rooms r ON r.tenant_id = t.id
		LEFT JOIN contracts c ON c.tenant_id = t.id AND c.status = 'active'
		GROUP BY t.id, t.name
		ORDER BY rooms DESC LIMIT 10
	`).Scan(&stats)

	result := map[string]interface{}{
		"统计": map[string]interface{}{
			"总租户数":   totalTenants,
			"有效租户数": activeTenants,
			"非有效租户数": inactiveTenants,
		},
		"租户排名(按房间数)": stats,
	}

	data, _ := json.MarshalIndent(result, "", "  ")
	return string(data), nil
}

// GetAdminContracts returns system-wide contract statistics for admin
func (e *ToolExecutor) GetAdminContracts(args map[string]interface{}) (string, error) {
	statusFilter := ""
	if status, ok := args["status"].(string); ok && status != "" {
		statusFilter = status
	}

	query := e.db.Model(&model.Contract{})
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var contracts []model.Contract
	if err := query.Order("created_at DESC").Limit(20).Find(&contracts).Error; err != nil {
		return "", fmt.Errorf("查询合同失败: %w", err)
	}

	var total, active, expired, terminated int64
	e.db.Model(&model.Contract{}).Count(&total)
	e.db.Model(&model.Contract{}).Where("status = ?", "active").Count(&active)
	e.db.Model(&model.Contract{}).Where("status = ?", "expired").Count(&expired)
	e.db.Model(&model.Contract{}).Where("status = ?", "terminated").Count(&terminated)

	statusText := map[string]string{"draft": "草稿", "active": "生效中", "expired": "已过期", "terminated": "已终止"}

	var records []map[string]interface{}
	records = append(records, map[string]interface{}{
		"统计": map[string]interface{}{
			"总合同数":  total,
			"生效中":   active,
			"已过期":   expired,
			"已终止":   terminated,
		},
	})

	for _, c := range contracts {
		tenantName := ""
		if c.TenantID > 0 {
			var t model.Tenant
			if err := e.db.First(&t, c.TenantID).Error; err == nil {
				tenantName = t.Name
			}
		}
		records = append(records, map[string]interface{}{
			"租户":    tenantName,
			"合同编号": c.ContractNo,
			"租金":    fmt.Sprintf("%.2f元/月", c.Amount),
			"开始日期": c.StartDate.Format("2006-01-02"),
			"结束日期": c.EndDate.Format("2006-01-02"),
			"状态":    statusText[c.Status],
		})
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	return string(data), nil
}

// GetAdminMaintenance returns system-wide maintenance statistics for admin
func (e *ToolExecutor) GetAdminMaintenance(args map[string]interface{}) (string, error) {
	statusFilter := ""
	if status, ok := args["status"].(string); ok && status != "" {
		statusFilter = status
	}

	query := e.db.Model(&model.Maintenance{})
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	var maintenances []model.Maintenance
	if err := query.Order("created_at DESC").Limit(20).Find(&maintenances).Error; err != nil {
		return "", fmt.Errorf("查询维修工单失败: %w", err)
	}

	var total, pending, processing, completed int64
	e.db.Model(&model.Maintenance{}).Count(&total)
	e.db.Model(&model.Maintenance{}).Where("status = ?", "pending").Count(&pending)
	e.db.Model(&model.Maintenance{}).Where("status = ?", "processing").Count(&processing)
	e.db.Model(&model.Maintenance{}).Where("status = ?", "completed").Count(&completed)

	statusText := map[string]string{"pending": "待处理", "processing": "处理中", "completed": "已完成", "cancelled": "已取消"}
	typeText := map[string]string{"electrical": "电气", "plumbing": "水管", "appliance": "家电", "furniture": "家具", "other": "其他"}
	priorityText := map[string]string{"low": "低", "medium": "中", "high": "高", "urgent": "紧急"}

	var records []map[string]interface{}
	records = append(records, map[string]interface{}{
		"统计": map[string]interface{}{
			"总工单数":  total,
			"待处理":   pending,
			"处理中":   processing,
			"已完成":   completed,
		},
	})

	for _, m := range maintenances {
		records = append(records, map[string]interface{}{
			"工单号":   m.TicketNo,
			"房间号":   m.RoomNo,
			"类型":    typeText[m.Type],
			"优先级":   priorityText[m.Priority],
			"状态":    statusText[m.Status],
			"描述":    m.Description,
			"创建时间": m.CreatedAt.Format("2006-01-02 15:04"),
			"负责人":   m.Assignee,
		})
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	return string(data), nil
}