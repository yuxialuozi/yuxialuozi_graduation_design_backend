package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/middleware"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

// TenantPortalHandler 租户端 API 处理器
type TenantPortalHandler struct {
	tenantPortalService *service.TenantPortalService
}

func NewTenantPortalHandler(tenantPortalService *service.TenantPortalService) *TenantPortalHandler {
	return &TenantPortalHandler{tenantPortalService: tenantPortalService}
}

// GetProfile godoc
// @Summary 获取租户个人信息
// @Description 获取当前登录租户的详细信息（租户信息）
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.TenantProfileResponse} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "租户不存在"
// @Router /tenant/profile [get]
func (h *TenantPortalHandler) GetProfile(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	profile, err := h.tenantPortalService.GetProfile(userID, tenantID)
	if err != nil {
		response.NotFound(c, "租户不存在")
		return
	}

	response.Success(c, profile)
}

// GetContracts godoc
// @Summary 获取租户合同列表
// @Description 获取当前登录租户的合同列表
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Router /tenant/contracts [get]
func (h *TenantPortalHandler) GetContracts(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	result, err := h.tenantPortalService.GetContracts(tenantID, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(c, "获取合同列表失败")
		return
	}

	response.Success(c, result)
}

// GetRooms godoc
// @Summary 获取租户房间列表
// @Description 获取当前登录租户的房间列表
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Router /tenant/rooms [get]
func (h *TenantPortalHandler) GetRooms(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	result, err := h.tenantPortalService.GetRooms(tenantID, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(c, "获取房间列表失败")
		return
	}

	response.Success(c, result)
}

// GetFees godoc
// @Summary 获取租户账单列表
// @Description 获取当前登录租户的账单列表
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query string false "状态筛选"
// @Param feeType query string false "费用类型筛选"
// @Param period query string false "账期筛选"
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Router /tenant/fees [get]
func (h *TenantPortalHandler) GetFees(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.FeeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	result, err := h.tenantPortalService.GetFees(tenantID, req.Page, req.PageSize, req.Status, req.FeeType, req.Period)
	if err != nil {
		response.InternalError(c, "获取账单列表失败")
		return
	}

	response.Success(c, result)
}

// PayFee godoc
// @Summary 缴纳账单
// @Description 租户缴纳账单
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "账单ID"
// @Param request body dto.PayFeeRequest false "缴费请求"
// @Success 200 {object} response.Response "缴费成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "账单不存在"
// @Router /tenant/fees/{id}/pay [post]
func (h *TenantPortalHandler) PayFee(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var idDto dto.IDRequest
	if err := c.ShouldBindUri(&idDto); err != nil {
		response.BadRequest(c, "无效的账单ID")
		return
	}

	var req dto.PayFeeRequest
	// 只有当请求有body时才尝试解析JSON，避免空body导致解析失败
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "请求参数格式错误，请检查JSON格式")
			return
		}
	}

	// 设置默认缴费时间
	paidDate := time.Now()
	if req.PaidDate != "" {
		parsed, err := time.Parse("2006-01-02", req.PaidDate)
		if err == nil {
			paidDate = parsed
		}
	}

	err := h.tenantPortalService.PayFee(idDto.ID, tenantID, &paidDate)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetMaintenance godoc
// @Summary 获取维修工单列表
// @Description 获取当前登录租户的维修工单列表
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query string false "状态筛选"
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Router /tenant/maintenance [get]
func (h *TenantPortalHandler) GetMaintenance(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)

	var req dto.MaintenanceListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	result, err := h.tenantPortalService.GetMaintenance(tenantID, req.Page, req.PageSize, req.Status)
	if err != nil {
		response.InternalError(c, "获取维修工单列表失败")
		return
	}

	response.Success(c, result)
}

// CreateMaintenance godoc
// @Summary 提交维修工单
// @Description 租户提交维修工单
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMaintenanceRequest true "维修工单请求"
// @Success 200 {object} response.Response "提交成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /tenant/maintenance [post]
func (h *TenantPortalHandler) CreateMaintenance(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	username := middleware.GetUsername(c)

	var req dto.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err := h.tenantPortalService.CreateMaintenance(tenantID, username, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetDashboard godoc
// @Summary 获取租户仪表盘数据
// @Description 获取当前登录租户的仪表盘统计数据
// @Tags 租户端
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=service.TenantDashboardResponse} "获取成功"
// @Router /tenant/dashboard [get]
func (h *TenantPortalHandler) GetDashboard(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	username := middleware.GetUsername(c)

	dashboard, err := h.tenantPortalService.GetDashboard(tenantID, username)
	if err != nil {
		response.InternalError(c, "获取仪表盘数据失败")
		return
	}

	response.Success(c, dashboard)
}