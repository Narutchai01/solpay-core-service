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
	"github.com/google/uuid"
)

const userFrontCardBucketName = "kyc"

type UserService interface {
	ApprovalStatus(req *request.ApprovalStatus) (*response.UserResponse, error)
	CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error)
	GetUsers(req request.UserQuery) (*response.UserListResponse, error)
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

	if req.BackCardImage == nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "back_card_image is required", entities.ErrBadRequest)
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

	backCardFile, err := req.BackCardImage.Open()
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "failed to open back_card_image", err)
	}
	defer backCardFile.Close()

	backCardBytes, err := io.ReadAll(backCardFile)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to read back_card_image", err)
	}

	backCardURL, err := s.storage.UploadFile(userFrontCardBucketName, buildFrontCardObjectPath(req.IDCard, req.BackCardImage.Filename), backCardBytes)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to upload back_card_image", err)
	}

	user := &entities.User{
		IDCard:       req.IDCard,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		FrontCardURL: frontCardURL,
		BackCardURL:  backCardURL,
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

	return fmt.Sprintf("users/%s/%d%s", cleanIDCard, time.Now().UnixNano(), ext)
}

func sanitizeStorageSegment(value string) string {
	cleaned := strings.TrimSpace(value)
	cleaned = strings.NewReplacer("/", "_", "\\", "_", " ", "_").Replace(cleaned)
	if cleaned == "" {
		return "unknown"
	}
	return cleaned
}

func (s *userService) ApprovalStatus(req *request.ApprovalStatus) (*response.UserResponse, error) {
	user, err := s.userRepo.GetUserByIDCard(req.IDCard)
	if err != nil {
		return nil, err
	}

	user.Status = req.Status
	if req.Status == string(entities.UserStatusApproved) && user.KYCToken == "" {
		kycToken, err := uuid.NewV7()
		if err != nil {
			return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to generate kyc token", err)
		}
		user.KYCToken = kycToken.String()
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
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

func (s *userService) GetUsers(req request.UserQuery) (*response.UserListResponse, error) {
	users, err := s.userRepo.GetUsers(req)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.userRepo.CountUsers(req)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	var userResponses []*response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &response.UserResponse{
			ID:           user.ID,
			IDCard:       user.IDCard,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			BirthDate:    user.BirthDate,
			Status:       user.Status,
			ExpireDate:   user.ExpireDate,
			FrontCardURL: user.FrontCardURL,
			BackCardURL:  user.BackCardURL,
		})
	}

	return &response.UserListResponse{
		Rows:       userResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
