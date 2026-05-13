package repositories

import (
	"errors"
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
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
	var existing entities.User
	if err := r.db.Where("account_id = ?", user.AccountID).First(&existing).Error; err == nil {
		user.ID = existing.ID
		return r.db.Save(user).Error
	}
	return r.db.Create(user).Error
}

func (r *gormUserRepository) GetUserByIDCard(idCard string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("id_card = ?", idCard).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) GetUserByAccountID(accountID uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("account_id = ?", accountID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) UpdateUser(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *gormUserRepository) GetUsers(query request.UserQuery) ([]*entities.User, error) {
	var users []*entities.User
	dbQuery := r.db.Model(&entities.User{})

	status := strings.ToUpper(strings.TrimSpace(query.Status))
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}

	page := query.Page
	if page <= 0 {
		page = 1
	}

	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	if err := dbQuery.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *gormUserRepository) CountUsers(query request.UserQuery) (int, error) {
	var count int64
	dbQuery := r.db.Model(&entities.User{})

	status := strings.ToUpper(strings.TrimSpace(query.Status))
	if status != "" {
		dbQuery = dbQuery.Where("status = ?", status)
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
