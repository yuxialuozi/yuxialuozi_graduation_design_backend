package repository

import (
	"gorm.io/gorm"

	"yuxialuozi_graduation_design_backend/internal/model"
)

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *RoomRepository) FindByID(id uint) (*model.Room, error) {
	var room model.Room
	if err := r.db.Preload("Tenant").First(&room, id).Error; err != nil {
		return nil, err
	}
	if room.Tenant != nil {
		room.TenantName = room.Tenant.Name
	}
	return &room, nil
}

func (r *RoomRepository) FindByRoomNo(roomNo string) (*model.Room, error) {
	var room model.Room
	if err := r.db.Where("room_no = ?", roomNo).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) Update(room *model.Room) error {
	return r.db.Save(room).Error
}

func (r *RoomRepository) Delete(id uint) error {
	return r.db.Delete(&model.Room{}, id).Error
}

func (r *RoomRepository) List(page, pageSize int, keyword, building, status string) ([]model.Room, int64, error) {
	var rooms []model.Room
	var total int64

	query := r.db.Model(&model.Room{}).Preload("Tenant")

	if keyword != "" {
		query = query.Where("room_no ILIKE ?", "%"+keyword+"%")
	}
	if building != "" {
		query = query.Where("building = ?", building)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("room_no ASC").Find(&rooms).Error; err != nil {
		return nil, 0, err
	}

	for i := range rooms {
		if rooms[i].Tenant != nil {
			rooms[i].TenantName = rooms[i].Tenant.Name
		}
	}

	return rooms, total, nil
}

func (r *RoomRepository) CountByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Room{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RoomRepository) CountTotal() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Room{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RoomRepository) GetBuildings() ([]string, error) {
	var buildings []string
	if err := r.db.Model(&model.Room{}).Distinct("building").Pluck("building", &buildings).Error; err != nil {
		return nil, err
	}
	return buildings, nil
}
