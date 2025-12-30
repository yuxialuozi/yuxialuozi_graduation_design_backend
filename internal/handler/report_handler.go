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

func (h *ReportHandler) GetOccupancy(c *gin.Context) {
	start, end := h.parseTimeRange(c)

	report, err := h.reportService.GetOccupancyReport(start, end)
	if err != nil {
		response.InternalError(c, "获取出租率统计失败")
		return
	}

	response.Success(c, report)
}

func (h *ReportHandler) GetFeeComposition(c *gin.Context) {
	start, end := h.parseTimeRange(c)

	composition, err := h.reportService.GetFeeComposition(start, end)
	if err != nil {
		response.InternalError(c, "获取费用构成失败")
		return
	}

	response.Success(c, composition)
}

func (h *ReportHandler) GetMaintenanceStats(c *gin.Context) {
	start, end := h.parseTimeRange(c)

	stats, err := h.reportService.GetMaintenanceStats(start, end)
	if err != nil {
		response.InternalError(c, "获取维修统计失败")
		return
	}

	response.Success(c, stats)
}

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

func (h *ReportHandler) GetDashboard(c *gin.Context) {
	data, err := h.reportService.GetDashboardData()
	if err != nil {
		response.InternalError(c, "获取仪表盘数据失败")
		return
	}

	response.Success(c, data)
}
