package service

import (
	"errors"

	"yuxialuozi_graduation_design_backend/internal/model"
	"yuxialuozi_graduation_design_backend/internal/repository"
)

type RoomService struct {
	roomRepo   *repository.RoomRepository
	tenantRepo *repository.TenantRepository
}

func NewRoomService(roomRepo *repository.RoomRepository, tenantRepo *repository.TenantRepository) *RoomService {
	return &RoomService{
		roomRepo:   roomRepo,
		tenantRepo: tenantRepo,
	}
}

func (s *RoomService) Create(room *model.Room) error {
	return s.roomRepo.Create(room)
}

func (s *RoomService) GetByID(id uint) (*model.Room, error) {
	return s.roomRepo.FindByID(id)
}

func (s *RoomService) Update(room *model.Room) error {
	return s.roomRepo.Update(room)
}

func (s *RoomService) Delete(id uint) error {
	return s.roomRepo.Delete(id)
}

func (s *RoomService) List(page, pageSize int, keyword, building, status string) ([]model.Room, int64, error) {
	return s.roomRepo.List(page, pageSize, keyword, building, status)
}

func (s *RoomService) AssignTenant(roomID uint, tenantID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return err
	}

	if room.Status == "occupied" && room.TenantID != nil && *room.TenantID != tenantID {
		return errors.New("房间已被占用")
	}

	_, err = s.tenantRepo.FindByID(tenantID)
	if err != nil {
		return errors.New("租户不存在")
	}

	room.TenantID = &tenantID
	room.Status = "occupied"
	return s.roomRepo.Update(room)
}

func (s *RoomService) ReleaseTenant(roomID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return err
	}

	room.TenantID = nil
	room.Status = "vacant"
	return s.roomRepo.Update(room)
}

func (s *RoomService) GetBuildings() ([]string, error) {
	return s.roomRepo.GetBuildings()
}
