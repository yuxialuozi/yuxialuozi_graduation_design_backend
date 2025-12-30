package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type ContractHandler struct {
	contractService *service.ContractService
}

func NewContractHandler(contractService *service.ContractService) *ContractHandler {
	return &ContractHandler{contractService: contractService}
}

// List godoc
// @Summary 获取合同列表
// @Description 分页获取合同列表，支持关键字搜索、状态筛选和日期范围
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键字"
// @Param status query string false "状态筛选" Enums(draft, active, expired, terminated)
// @Param startDateFrom query string false "开始日期起始 (YYYY-MM-DD)"
// @Param startDateTo query string false "开始日期结束 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /contracts [get]
func (h *ContractHandler) List(c *gin.Context) {
	var req dto.ContractListRequest
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

	var startDateFrom, startDateTo *time.Time
	if req.StartDateFrom != "" {
		t, err := time.Parse("2006-01-02", req.StartDateFrom)
		if err == nil {
			startDateFrom = &t
		}
	}
	if req.StartDateTo != "" {
		t, err := time.Parse("2006-01-02", req.StartDateTo)
		if err == nil {
			startDateTo = &t
		}
	}

	contracts, total, err := h.contractService.List(req.Page, req.PageSize, req.Keyword, req.Status, startDateFrom, startDateTo)
	if err != nil {
		response.InternalError(c, "获取合同列表失败")
		return
	}

	response.Success(c, dto.NewPageResult(contracts, total, req.Page, req.PageSize))
}

// GetByID godoc
// @Summary 获取合同详情
// @Description 根据 ID 获取合同详细信息
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同 ID"
// @Success 200 {object} response.Response{data=model.Contract} "获取成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 404 {object} response.Response "合同不存在"
// @Router /contracts/{id} [get]
func (h *ContractHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	contract, err := h.contractService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "合同不存在")
		return
	}

	response.Success(c, contract)
}

// Create godoc
// @Summary 创建合同
// @Description 创建新的合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateContractRequest true "创建合同请求"
// @Success 200 {object} response.Response{data=model.Contract} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "创建失败"
// @Router /contracts [post]
func (h *ContractHandler) Create(c *gin.Context) {
	var req dto.CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	contract := &model.Contract{
		TenantID:   req.TenantID,
		ContractNo: req.ContractNo,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Amount:     req.Amount,
		Status:     req.Status,
	}

	if contract.Status == "" {
		contract.Status = "draft"
	}

	if err := h.contractService.Create(contract); err != nil {
		response.InternalError(c, "创建合同失败")
		return
	}

	response.Success(c, contract)
}

// Update godoc
// @Summary 更新合同
// @Description 更新合同信息
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同 ID"
// @Param request body dto.UpdateContractRequest true "更新合同请求"
// @Success 200 {object} response.Response{data=model.Contract} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "合同不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /contracts/{id} [put]
func (h *ContractHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	contract, err := h.contractService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "合同不存在")
		return
	}

	var req dto.UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.TenantID > 0 {
		contract.TenantID = req.TenantID
	}
	if req.ContractNo != "" {
		contract.ContractNo = req.ContractNo
	}
	if !req.StartDate.IsZero() {
		contract.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		contract.EndDate = req.EndDate
	}
	if req.Amount > 0 {
		contract.Amount = req.Amount
	}
	if req.Status != "" {
		contract.Status = req.Status
	}

	if err := h.contractService.Update(contract); err != nil {
		response.InternalError(c, "更新合同失败")
		return
	}

	response.Success(c, contract)
}

// Delete godoc
// @Summary 删除合同
// @Description 删除指定合同
// @Tags 合同管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "合同 ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "删除失败"
// @Router /contracts/{id} [delete]
func (h *ContractHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.contractService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除合同失败")
		return
	}

	response.Success(c, nil)
}
