package services

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

const userFrontCardBucketName = "kyc"

type UserService interface {
	CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error)
}

type userService struct {
	userRepo ports.UserRepository
	storage  ports.Storage
}

func NewUserService(userRepo ports.UserRepository, storage ports.Storage) UserService {
	return &userService{
		userRepo: userRepo,
		storage:  storage,
	}
}

func (s *userService) CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error) {
	if req == nil || req.FrontCardImage == nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "front_card_image is required", entities.ErrBadRequest)
	}

	frontCardFile, err := req.FrontCardImage.Open()
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "failed to open front_card_image", err)
	}
	defer frontCardFile.Close()

	frontCardBytes, err := io.ReadAll(frontCardFile)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to read front_card_image", err)
	}

	frontCardURL, err := s.storage.UploadFile(userFrontCardBucketName, buildFrontCardObjectPath(req.IDCard, req.FrontCardImage.Filename), frontCardBytes)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to upload front_card_image", err)
	}

	user := &entities.User{
		IDCard:       req.IDCard,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		FrontCardURL: frontCardURL,
		BackCardURL:  "dasdasd",
		BirthDate:    req.BirthDate,
		Status:       string(entities.UserStatusPending),
		ExpireDate:   req.ExpireDate,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	dto := &response.UserResponse{
		ID:           user.ID,
		IDCard:       user.IDCard,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		BirthDate:    user.BirthDate,
		Status:       user.Status,
		ExpireDate:   user.ExpireDate,
		FrontCardURL: user.FrontCardURL,
		BackCardURL:  user.BackCardURL,
	}

	return dto, nil
}

func buildFrontCardObjectPath(idCard, filename string) string {
	cleanIDCard := sanitizeStorageSegment(idCard)
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(filename)))
	if ext == "" {
		ext = ".img"
	}

	return fmt.Sprintf("users/%s/front_card_%d%s", cleanIDCard, time.Now().UnixNano(), ext)
}

func sanitizeStorageSegment(value string) string {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.NewReplacer("/", "_", "\\", "_", " ", "_").Replace(cleaned)
	if cleaned == "" {
		return "unknown"
	}
	return cleaned
}
