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

type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
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
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role, s.config.JWT.Secret, expire)
	if err != nil {
		return nil, errors.New("生成 token 失败")
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) GetCurrentUser(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *AuthService) CreateDefaultAdmin() error {
	_, err := s.userRepo.FindByUsername("admin")
	if err == nil {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
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
