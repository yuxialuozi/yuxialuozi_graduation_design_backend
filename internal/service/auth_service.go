package service

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"yuxialuozi_graduation_design_backend/internal/config"
	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
	"yuxialuozi_graduation_design_backend/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response for login API.
// @Description LoginResponse represents the response for login API with JWT token and user info.
type LoginResponse struct {
	Token    string              `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserInfo *model.UserResponse `json:"userInfo"` // 改为 userInfo 字段，匹配前端期望
	// 保留 user 字段做兼容
	User *model.UserResponse `json:"-"` // 不在JSON中输出，但保留用于内部使用
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	expire, _ := time.ParseDuration(s.config.JWT.Expire)
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, user.TenantID, s.config.JWT.Secret, expire)
	if err != nil {
		return nil, errors.New("生成 token 失败")
	}

	// Convert model.User to model.UserResponse for Swagger
	userResponse := &model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Role:        user.Role,
		Permissions: []string(user.Permissions),
		TenantID:    user.TenantID,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}

	return &LoginResponse{
		Token:    token,
		UserInfo: userResponse,
		User:     userResponse,
	}, nil
}

func (s *AuthService) GetCurrentUser(userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return &model.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Role:        user.Role,
		Permissions: []string(user.Permissions),
		TenantID:    user.TenantID,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

func (s *AuthService) CreateDefaultAdmin() error {
	_, err := s.userRepo.FindByUsername("admin")
	if err == nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username:    "admin",
		Password:    string(hashedPassword),
		Nickname:    "管理员",
		Role:        "admin",
		Permissions: []string{"*"},
	}

	return s.userRepo.Create(admin)
}
