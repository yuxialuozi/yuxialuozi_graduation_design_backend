package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type RoomHandler struct {
	roomService *service.RoomService
}

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{roomService: roomService}
}

// List godoc
// @Summary 获取房间列表
// @Description 分页获取房间列表，支持关键字搜索、楼栋筛选和状态筛选
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键字"
// @Param building query string false "楼栋筛选"
// @Param status query string false "状态筛选" Enums(vacant, occupied, maintenance)
// @Success 200 {object} response.Response{data=dto.PageResult} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /rooms [get]
func (h *RoomHandler) List(c *gin.Context) {
	var req dto.RoomListRequest
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

	rooms, total, err := h.roomService.List(req.Page, req.PageSize, req.Keyword, req.Building, req.Status)
	if err != nil {
		response.InternalError(c, "获取房间列表失败")
		return
	}

	response.Success(c, dto.NewPageResult(rooms, total, req.Page, req.PageSize))
}

// GetByID godoc
// @Summary 获取房间详情
// @Description 根据 ID 获取房间详细信息
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间 ID"
// @Success 200 {object} response.Response{data=model.Room} "获取成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 404 {object} response.Response "房间不存在"
// @Router /rooms/{id} [get]
func (h *RoomHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	room, err := h.roomService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "房间不存在")
		return
	}

	response.Success(c, room)
}

// Create godoc
// @Summary 创建房间
// @Description 创建新的房间
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRoomRequest true "创建房间请求"
// @Success 200 {object} response.Response{data=model.Room} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "创建失败"
// @Router /rooms [post]
func (h *RoomHandler) Create(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	room := &model.Room{
		RoomNo:      req.RoomNo,
		Building:    req.Building,
		Floor:       req.Floor,
		Area:        req.Area,
		MonthlyRent: req.MonthlyRent,
		Status:      req.Status,
	}

	if room.Status == "" {
		room.Status = "vacant"
	}

	if err := h.roomService.Create(room); err != nil {
		response.InternalError(c, "创建房间失败")
		return
	}

	response.Success(c, room)
}

// Update godoc
// @Summary 更新房间
// @Description 更新房间信息
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间 ID"
// @Param request body dto.UpdateRoomRequest true "更新房间请求"
// @Success 200 {object} response.Response{data=model.Room} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "房间不存在"
// @Failure 500 {object} response.Response "更新失败"
// @Router /rooms/{id} [put]
func (h *RoomHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	room, err := h.roomService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "房间不存在")
		return
	}

	var req dto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if req.RoomNo != "" {
		room.RoomNo = req.RoomNo
	}
	if req.Building != "" {
		room.Building = req.Building
	}
	if req.Floor > 0 {
		room.Floor = req.Floor
	}
	if req.Area > 0 {
		room.Area = req.Area
	}
	if req.MonthlyRent > 0 {
		room.MonthlyRent = req.MonthlyRent
	}
	if req.Status != "" {
		room.Status = req.Status
	}

	if err := h.roomService.Update(room); err != nil {
		response.InternalError(c, "更新房间失败")
		return
	}

	response.Success(c, room)
}

// Delete godoc
// @Summary 删除房间
// @Description 删除指定房间
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间 ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "无效的 ID"
// @Failure 500 {object} response.Response "删除失败"
// @Router /rooms/{id} [delete]
func (h *RoomHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.roomService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除房间失败")
		return
	}

	response.Success(c, nil)
}

// AssignTenant godoc
// @Summary 分配租户
// @Description 将租户分配到指定房间
// @Tags 房间管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "房间 ID"
// @Param request body dto.AssignTenantRequest true "分配租户请求"
// @Success 200 {object} response.Response "分配成功"
// @Failure 400 {object} response.Response "请求参数错误或房间已被占用"
// @Router /rooms/{id}/assign [post]
func (h *RoomHandler) AssignTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req dto.AssignTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err := h.roomService.AssignTenant(uint(id), req.TenantID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, nil)
}
