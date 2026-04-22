package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &gormUserRepository{db: db}
}

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *gormUserRepository) GetUserByIDCard(idCard string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("id_card = ?", idCard).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}
