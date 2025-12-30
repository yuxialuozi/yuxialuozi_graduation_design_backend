package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type TenantHandler struct {
	tenantService *service.TenantService
}

func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService}
}

// List godoc
// @Summary 获取租户列表
// @Description 分页获取租户列表，支持关键字搜索和状态筛选
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键字"
// @Param status query string false "状态筛选" Enums(active, inactive)
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /tenants [get]
func (h *TenantHandler) List(c *gin.Context) {
	var req dto.TenantListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	tenants, total, err := h.tenantService.List(req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		response.InternalError(c, "获取租户列表失败")
		return
	}

	response.Success(c, dto.NewPageResult(tenants, total, req.Page, req.PageSize))
}

// GetByID godoc
// @Summary 获取租户详情
// @Description 根据 ID 获取租户详细信息
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "租户 ID"
// @Success 200 {object} response.Response{data=model.Tenant} "获取成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 404 {object} response.Response "租户不存在"
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	tenant, err := h.tenantService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "租户不存在")
		return
	}

	response.Success(c, tenant)
}

// Create godoc
// @Summary 创建租户
// @Description 创建新的租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTenantRequest true "创建租户请求"
// @Success 200 {object} response.Response{data=model.Tenant} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "创建失败"
// @Router /tenants [post]
func (h *TenantHandler) Create(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	tenant := &model.Tenant{
		Name:          req.Name,
		ContactPerson: req.ContactPerson,
		Phone:         req.Phone,
		Email:         req.Email,
		Status:        req.Status,
	}

	if tenant.Status == "" {
		tenant.Status = "active"
	}

	if err := h.tenantService.Create(tenant); err != nil {
		response.InternalError(c, "创建租户失败")
		return
	}

	response.Success(c, tenant)
}

// Update godoc
// @Summary 更新租户
// @Description 更新租户信息
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "租户 ID"
// @Param request body dto.UpdateTenantRequest true "更新租户请求"
// @Success 200 {object} response.Response{data=model.Tenant} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "租户不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /tenants/{id} [put]
func (h *TenantHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	tenant, err := h.tenantService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "租户不存在")
		return
	}

	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.ContactPerson != "" {
		tenant.ContactPerson = req.ContactPerson
	}
	if req.Phone != "" {
		tenant.Phone = req.Phone
	}
	if req.Email != "" {
		tenant.Email = req.Email
	}
	if req.Status != "" {
		tenant.Status = req.Status
	}

	if err := h.tenantService.Update(tenant); err != nil {
		response.InternalError(c, "更新租户失败")
		return
	}

	response.Success(c, tenant)
}

// Delete godoc
// @Summary 删除租户
// @Description 删除指定租户
// @Tags 租户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "租户 ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "删除失败"
// @Router /tenants/{id} [delete]
func (h *TenantHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.tenantService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除租户失败")
		return
	}

	response.Success(c, nil)
}
