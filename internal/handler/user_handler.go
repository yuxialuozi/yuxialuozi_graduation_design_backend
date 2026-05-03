package handler

import (
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
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// 加密密码
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

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Nickname  string `json:"nickname"`
	Role      string `json:"role"`
	TenantID  uint   `json:"tenantId"`
}