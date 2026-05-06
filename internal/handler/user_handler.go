package handler

import (
	"strconv"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/repository"
	"yuxialuozi_graduation_design_backend/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建新用户（仅用于测试）
// @Tags 用户
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateUserRequest true "用户信息"
// @Success 200 {object} response.Response "创建成功"
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	err = h.userRepo.CreateUser(req.Username, string(hashedPassword), req.Nickname, req.Role, req.TenantID)
	if err != nil {
		response.BadRequest(c, "创建用户失败，用户名可能已存在")
		return
	}

	response.Success(c, nil)
}

// List godoc
// @Summary 获取用户列表
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param keyword query string false "关键词搜索"
// @Param role query string false "角色筛选"
// @Success 200 {object} response.Response
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	var req dto.UserListRequest
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

	users, total, err := h.userRepo.ListFiltered(req.Page, req.PageSize, req.Keyword, req.Role)
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Success(c, dto.NewPageResult(users, total, req.Page, req.PageSize))
}

// GetByID godoc
// @Summary 获取用户详情
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// Update godoc
// @Summary 更新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body dto.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.TenantID > 0 {
		updates["tenant_id"] = req.TenantID
	}

	if err := h.userRepo.UpdateFields(uint(id), updates); err != nil {
		response.InternalError(c, "更新用户失败")
		return
	}

	response.Success(c, nil)
}

// Delete godoc
// @Summary 删除用户
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除用户失败")
		return
	}

	response.Success(c, nil)
}

// ResetPassword godoc
// @Summary 重置用户密码
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Router /users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	defaultPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "密码加密失败")
		return
	}

	user.Password = string(hashedPassword)
	if err := h.userRepo.Update(user); err != nil {
		response.InternalError(c, "重置密码失败")
		return
	}

	response.Success(c, gin.H{"newPassword": defaultPassword})
}