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
