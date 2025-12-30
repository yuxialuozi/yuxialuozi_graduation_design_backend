package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type FeeHandler struct {
	feeService *service.FeeService
}

func NewFeeHandler(feeService *service.FeeService) *FeeHandler {
	return &FeeHandler{feeService: feeService}
}

// List godoc
// @Summary 获取费用列表
// @Description 分页获取费用列表，支持多条件筛选
// @Tags 费用管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param tenantId query int false "租户 ID"
// @Param roomNo query string false "房间号"
// @Param feeType query string false "费用类型" Enums(rent, water, electricity, property, other)
// @Param status query string false "状态" Enums(unpaid, paid, overdue)
// @Param period query string false "账期 (如: 2024-03)"
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /fees [get]
func (h *FeeHandler) List(c *gin.Context) {
	var req dto.FeeListRequest
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

	fees, total, err := h.feeService.List(req.Page, req.PageSize, req.TenantID, req.RoomNo, req.FeeType, req.Status, req.Period)
	if err != nil {
		response.InternalError(c, "获取费用列表失败")
		return
	}

	response.Success(c, dto.NewPageResult(fees, total, req.Page, req.PageSize))
}

// GetByID godoc
// @Summary 获取费用详情
// @Description 根据 ID 获取费用详细信息
// @Tags 费用管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "费用 ID"
// @Success 200 {object} response.Response{data=model.Fee} "获取成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 404 {object} response.Response "费用记录不存在"
// @Router /fees/{id} [get]
func (h *FeeHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	fee, err := h.feeService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "费用记录不存在")
		return
	}

	response.Success(c, fee)
}

// Create godoc
// @Summary 创建费用
// @Description 创建新的费用记录
// @Tags 费用管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateFeeRequest true "创建费用请求"
// @Success 200 {object} response.Response{data=model.Fee} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "创建失败"
// @Router /fees [post]
func (h *FeeHandler) Create(c *gin.Context) {
	var req dto.CreateFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	fee := &model.Fee{
		TenantID: req.TenantID,
		RoomNo:   req.RoomNo,
		FeeType:  req.FeeType,
		Amount:   req.Amount,
		Period:   req.Period,
		DueDate:  req.DueDate,
		Status:   req.Status,
	}

	if fee.Status == "" {
		fee.Status = "unpaid"
	}

	if err := h.feeService.Create(fee); err != nil {
		response.InternalError(c, "创建费用记录失败")
		return
	}

	response.Success(c, fee)
}

// Update godoc
// @Summary 更新费用
// @Description 更新费用记录
// @Tags 费用管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "费用 ID"
// @Param request body dto.UpdateFeeRequest true "更新费用请求"
// @Success 200 {object} response.Response{data=model.Fee} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "费用记录不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /fees/{id} [put]
func (h *FeeHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	fee, err := h.feeService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "费用记录不存在")
		return
	}

	var req dto.UpdateFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.TenantID > 0 {
		fee.TenantID = req.TenantID
	}
	if req.RoomNo != "" {
		fee.RoomNo = req.RoomNo
	}
	if req.FeeType != "" {
		fee.FeeType = req.FeeType
	}
	if req.Amount > 0 {
		fee.Amount = req.Amount
	}
	if req.Period != "" {
		fee.Period = req.Period
	}
	if !req.DueDate.IsZero() {
		fee.DueDate = req.DueDate
	}
	if req.Status != "" {
		fee.Status = req.Status
	}

	if err := h.feeService.Update(fee); err != nil {
		response.InternalError(c, "更新费用记录失败")
		return
	}

	response.Success(c, fee)
}

// Delete godoc
// @Summary 删除费用
// @Description 删除指定费用记录
// @Tags 费用管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "费用 ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "删除失败"
// @Router /fees/{id} [delete]
func (h *FeeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.feeService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除费用记录失败")
		return
	}

	response.Success(c, nil)
}

// Pay godoc
// @Summary 确认缴费
// @Description 确认费用已缴纳
// @Tags 费用管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "费用 ID"
// @Param request body dto.PayFeeRequest false "缴费请求"
// @Success 200 {object} response.Response "缴费成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "确认缴费失败"
// @Router /fees/{id}/pay [post]
func (h *FeeHandler) Pay(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req dto.PayFeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.PaidDate = nil
	}

	if err := h.feeService.Pay(uint(id), req.PaidDate); err != nil {
		response.InternalError(c, "确认缴费失败")
		return
	}

	response.Success(c, nil)
}
