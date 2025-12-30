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
