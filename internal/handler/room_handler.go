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
