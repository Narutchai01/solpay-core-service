package services

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	CreateAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, error)
}

type adminService struct {
	adminRepo ports.AdminRepository
}

func NewAdminService(adminRepo ports.AdminRepository) AdminService {
	return &adminService{
		adminRepo: adminRepo,
	}
}

func (s *adminService) CreateAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to hash password", err)
	}
	admin := &entities.AdminEntity{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := s.adminRepo.CreateAdmin(ctx, admin); err != nil {
		return nil, err
	}
	return admin, nil
}
