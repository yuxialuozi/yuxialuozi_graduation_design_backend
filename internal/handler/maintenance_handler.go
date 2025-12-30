package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type MaintenanceHandler struct {
	maintenanceService *service.MaintenanceService
}

func NewMaintenanceHandler(maintenanceService *service.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{maintenanceService: maintenanceService}
}

// List godoc
// @Summary 获取维修工单列表
// @Description 分页获取维修工单列表，支持多条件筛选
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键字"
// @Param type query string false "维修类型" Enums(electrical, plumbing, appliance, furniture, other)
// @Param status query string false "状态" Enums(pending, processing, completed, cancelled)
// @Param priority query string false "优先级" Enums(low, medium, high, urgent)
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /maintenance [get]
func (h *MaintenanceHandler) List(c *gin.Context) {
	var req dto.MaintenanceListRequest
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

	maintenances, total, err := h.maintenanceService.List(req.Page, req.PageSize, req.Keyword, req.Type, req.Status, req.Priority)
	if err != nil {
		response.InternalError(c, "获取维修工单列表失败")
		return
	}

	response.Success(c, dto.NewPageResult(maintenances, total, req.Page, req.PageSize))
}

// GetByID godoc
// @Summary 获取维修工单详情
// @Description 根据 ID 获取维修工单详细信息
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单 ID"
// @Success 200 {object} response.Response{data=model.Maintenance} "获取成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 404 {object} response.Response "维修工单不存在"
// @Router /maintenance/{id} [get]
func (h *MaintenanceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	maintenance, err := h.maintenanceService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "维修工单不存在")
		return
	}

	response.Success(c, maintenance)
}

// Create godoc
// @Summary 创建维修工单
// @Description 创建新的维修工单
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMaintenanceRequest true "创建工单请求"
// @Success 200 {object} response.Response{data=model.Maintenance} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "创建失败"
// @Router /maintenance [post]
func (h *MaintenanceHandler) Create(c *gin.Context) {
	var req dto.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	maintenance := &model.Maintenance{
		TenantID:    req.TenantID,
		RoomNo:      req.RoomNo,
		Type:        req.Type,
		Description: req.Description,
		Priority:    req.Priority,
	}

	if maintenance.Priority == "" {
		maintenance.Priority = "medium"
	}

	if err := h.maintenanceService.Create(maintenance); err != nil {
		response.InternalError(c, "创建维修工单失败")
		return
	}

	response.Success(c, maintenance)
}

// Update godoc
// @Summary 更新维修工单
// @Description 更新维修工单信息
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单 ID"
// @Param request body dto.UpdateMaintenanceRequest true "更新工单请求"
// @Success 200 {object} response.Response{data=model.Maintenance} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "维修工单不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /maintenance/{id} [put]
func (h *MaintenanceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	maintenance, err := h.maintenanceService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "维修工单不存在")
		return
	}

	var req dto.UpdateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.TenantID > 0 {
		maintenance.TenantID = req.TenantID
	}
	if req.RoomNo != "" {
		maintenance.RoomNo = req.RoomNo
	}
	if req.Type != "" {
		maintenance.Type = req.Type
	}
	if req.Description != "" {
		maintenance.Description = req.Description
	}
	if req.Priority != "" {
		maintenance.Priority = req.Priority
	}
	if req.Status != "" {
		maintenance.Status = req.Status
	}
	if req.Assignee != "" {
		maintenance.Assignee = req.Assignee
	}

	if err := h.maintenanceService.Update(maintenance); err != nil {
		response.InternalError(c, "更新维修工单失败")
		return
	}

	response.Success(c, maintenance)
}

// Delete godoc
// @Summary 删除维修工单
// @Description 删除指定维修工单
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单 ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "删除失败"
// @Router /maintenance/{id} [delete]
func (h *MaintenanceHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.maintenanceService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除维修工单失败")
		return
	}

	response.Success(c, nil)
}

// Assign godoc
// @Summary 指派维修人员
// @Description 为维修工单指派维修人员
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单 ID"
// @Param request body dto.AssignMaintenanceRequest true "指派请求"
// @Success 200 {object} response.Response "指派成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "指派失败"
// @Router /maintenance/{id}/assign [post]
func (h *MaintenanceHandler) Assign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req dto.AssignMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.maintenanceService.Assign(uint(id), req.Assignee); err != nil {
		response.InternalError(c, "指派维修人员失败")
		return
	}

	response.Success(c, nil)
}

// Complete godoc
// @Summary 完成工单
// @Description 标记维修工单为已完成
// @Tags 维修管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "工单 ID"
// @Param request body dto.CompleteMaintenanceRequest false "完成请求"
// @Success 200 {object} response.Response "完成成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "完成失败"
// @Router /maintenance/{id}/complete [post]
func (h *MaintenanceHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req dto.CompleteMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.CompletedAt = nil
	}

	if err := h.maintenanceService.Complete(uint(id), req.CompletedAt); err != nil {
		response.InternalError(c, "完成工单失败")
		return
	}

	response.Success(c, nil)
}
