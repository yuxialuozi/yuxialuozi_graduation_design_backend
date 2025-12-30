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
