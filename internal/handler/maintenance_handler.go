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
