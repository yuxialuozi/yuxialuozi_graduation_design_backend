package service

import (
	"errors"

	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) List(page, pageSize int, keyword, role string) ([]model.User, int64, error) {
	return s.userRepo.ListFiltered(page, pageSize, keyword, role)
}

func (s *UserService) GetByID(id uint) (interface{}, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Create(username, password, nickname, role string, tenantID uint) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密失败")
	}
	return s.userRepo.CreateUser(username, string(hashedPassword), nickname, role, tenantID)
}

func (s *UserService) Update(id uint, nickname, phone, email, role, status string, tenantID uint) error {
	updates := map[string]interface{}{}
	if nickname != "" {
		updates["nickname"] = nickname
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if email != "" {
		updates["email"] = email
	}
	if role != "" {
		updates["role"] = role
	}
	if status != "" {
		updates["status"] = status
	}
	if tenantID > 0 {
		updates["tenant_id"] = tenantID
	}
	if len(updates) == 0 {
		return nil
	}
	return s.userRepo.UpdateFields(id, updates)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserService) ResetPassword(id uint) (string, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return "", errors.New("用户不存在")
	}
	// 设为默认密码 123456
	defaultPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("密码加密失败")
	}
	user.Password = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return "", errors.New("重置密码失败")
	}
	return defaultPassword, nil
}