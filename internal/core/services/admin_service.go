package services

import (
	"context"
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/config"
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"github.com/Narutchai01/solpay-core-service/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	CreateAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, error)
	LoginAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, string, error)
	GetProfile(ctx context.Context, userID string) (*entities.AdminEntity, error)
}

type adminService struct {
	adminRepo ports.AdminRepository
	cfg       *config.Config
}

func NewAdminService(adminRepo ports.AdminRepository, cfg *config.Config) AdminService {
	return &adminService{
		adminRepo: adminRepo,
		cfg:       cfg,
	}
}

func (s *adminService) CreateAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, error) {
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	existing, err := s.adminRepo.GetAdminByUsername(ctx, req.Username)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, entities.NewAppError(
			entities.ErrTypeConflict,
			"username already exists",
			nil,
		)
	}

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

func (s *adminService) LoginAdmin(ctx context.Context, req request.CreateAdminRequest) (*entities.AdminEntity, string, error) {
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	admin, err := s.adminRepo.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password))
	if err != nil {
		return nil, "", entities.NewAppError(
			entities.ErrTypeBadRequest,
			"invalid username or password",
			nil,
		)
	}

	token, err := utils.GenerateJWT(admin.ID.String(), admin.Username, s.cfg.SECRET_JWT)
	if err != nil {
		return nil, "", entities.NewAppError(
			entities.ErrTypeInternal,
			"failed to generate token",
			err,
		)
	}

	return admin, token, nil
}

func (s *adminService) GetProfile(ctx context.Context, userID string) (*entities.AdminEntity, error) {
	admin, err := s.adminRepo.GetAdminByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return admin, nil
}
