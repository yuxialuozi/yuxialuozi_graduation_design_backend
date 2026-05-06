package handler

import (
	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/middleware"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary 用户登录
// @Description 使用用户名和密码登录系统
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录请求"
// @Success 200 {object} response.Response{data=service.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	result, err := h.authService.Login(&service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, result)
}

// GetCurrentUser godoc
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.UserResponse} "获取成功"
// @Failure 401 {object} response.Response "未登录"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// Logout godoc
// @Summary 退出登录
// @Description 退出当前用户登录状态
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response "退出成功"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}

// ChangePassword godoc
// @Summary 修改密码
// @Description 当前用户修改自己的密码
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "密码修改请求"
// @Success 200 {object} response.Response "修改成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "原密码错误"
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误，新密码至少6位")
		return
	}

	err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "密码修改成功")
}

// GetProfile godoc
// @Summary 获取个人资料
// @Description 获取当前登录用户的详细资料
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Router /profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

// UpdateProfile godoc
// @Summary 更新个人资料
// @Description 更新当前登录用户的资料（昵称、手机、邮箱）
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateProfileRequest true "个人资料"
// @Success 200 {object} response.Response
// @Router /profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	err := h.authService.UpdateProfile(userID, req.Nickname, req.Phone, req.Email)
	if err != nil {
		response.InternalError(c, "更新资料失败")
		return
	}

	response.Success(c, "资料更新成功")
}
