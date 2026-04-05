package services

import (
	"context"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type AdminService interface {
	CreateAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, error)
}

type adminService struct {
	adminRepo ports.AdminRepository
	uow       ports.UnitOfWork
}

func NewAdminService(adminRepo ports.AdminRepository, uow ports.UnitOfWork) AdminService {
	return &adminService{
		adminRepo: adminRepo,
		uow:       uow,
	}
}

func (s *adminService) CreateAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, error) {
	result, err := s.uow.Execute(ctx, func(txCtx context.Context) (any, error) {
		admin := &entities.AdminEntity{
			Username: req.Username,
			Password: req.Password,
		}
		if err := s.adminRepo.CreateAdmin(txCtx, admin); err != nil {
			return nil, err
		}
		return admin, nil
	})
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to create admin", err)
	}
	return result.(*entities.AdminEntity), nil
}
