package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) parseTimeRange(c *gin.Context) (time.Time, time.Time) {
	var req dto.ReportQueryRequest
	c.ShouldBindQuery(&req)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	end := now

	if req.Start != "" {
		if t, err := time.Parse("2006-01-02", req.Start); err == nil {
			start = t
		}
	}
	if req.End != "" {
		if t, err := time.Parse("2006-01-02", req.End); err == nil {
			end = t.Add(24*time.Hour - time.Second)
		}
	}

	return start, end
}

// GetIncome godoc
// @Summary 收入统计
// @Description 获取指定时间范围内的收入统计数据
// @Tags 报表统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start query string false "开始日期 (YYYY-MM-DD)"
// @Param end query string false "结束日期 (YYYY-MM-DD)"
// @Param groupBy query string false "分组方式" Enums(month, type)
// @Success 200 {object} response.Response{data=service.IncomeReport} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /reports/income [get]
func (h *ReportHandler) GetIncome(c *gin.Context) {
	start, end := h.parseTimeRange(c)
	groupBy := c.Query("groupBy")

	report, err := h.reportService.GetIncomeReport(start, end, groupBy)
	if err != nil {
		response.InternalError(c, "获取收入统计失败")
		return
	}

	response.Success(c, report)
}

// GetOccupancy godoc
// @Summary 出租率统计
// @Description 获取当前出租率统计数据
// @Tags 报表统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start query string false "开始日期 (YYYY-MM-DD)"
// @Param end query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=service.OccupancyReport} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /reports/occupancy [get]
func (h *ReportHandler) GetOccupancy(c *gin.Context) {
	start, end := h.parseTimeRange(c)

	report, err := h.reportService.GetOccupancyReport(start, end)
	if err != nil {
		response.InternalError(c, "获取出租率统计失败")
		return
	}

	response.Success(c, report)
}

// GetFeeComposition godoc
// @Summary 费用构成
// @Description 获取指定时间范围内的费用构成数据
// @Tags 报表统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start query string false "开始日期 (YYYY-MM-DD)"
// @Param end query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=[]repository.FeeComposition} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /reports/fees/composition [get]
func (h *ReportHandler) GetFeeComposition(c *gin.Context) {
	start, end := h.parseTimeRange(c)

	composition, err := h.reportService.GetFeeComposition(start, end)
	if err != nil {
		response.InternalError(c, "获取费用构成失败")
		return
	}

	response.Success(c, composition)
}

// GetMaintenanceStats godoc
// @Summary 维修统计
// @Description 获取指定时间范围内的维修统计数据
// @Tags 报表统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start query string false "开始日期 (YYYY-MM-DD)"
// @Param end query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=service.MaintenanceReport} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /reports/maintenance/stats [get]
func (h *ReportHandler) GetMaintenanceStats(c *gin.Context) {
	start, end := h.parseTimeRange(c)

	stats, err := h.reportService.GetMaintenanceStats(start, end)
	if err != nil {
		response.InternalError(c, "获取维修统计失败")
		return
	}

	response.Success(c, stats)
}

// GetTenantRanking godoc
// @Summary 租户排行
// @Description 获取缴费金额租户排行榜
// @Tags 报表统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "返回数量" default(10)
// @Param start query string false "开始日期 (YYYY-MM-DD)"
// @Param end query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=[]repository.TenantFeeRanking} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /reports/tenants/ranking [get]
func (h *ReportHandler) GetTenantRanking(c *gin.Context) {
	var req dto.ReportQueryRequest
	c.ShouldBindQuery(&req)

	if req.Limit <= 0 {
		req.Limit = 10
	}

	start, end := h.parseTimeRange(c)

	ranking, err := h.reportService.GetTenantRanking(req.Limit, start, end)
	if err != nil {
		response.InternalError(c, "获取租户排行失败")
		return
	}

	response.Success(c, ranking)
}

// GetDashboard godoc
// @Summary 仪表盘数据
// @Description 获取仪表盘汇总数据
// @Tags 报表统计
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.DashboardData} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /reports/dashboard [get]
func (h *ReportHandler) GetDashboard(c *gin.Context) {
	data, err := h.reportService.GetDashboardData()
	if err != nil {
		response.InternalError(c, "获取仪表盘数据失败")
		return
	}

	response.Success(c, data)
}
