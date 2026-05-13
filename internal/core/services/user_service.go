package services

import (
	"context"
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
	CreateUser(req *request.CreateUserRequest, accountID uint) (*response.UserResponse, error)
	GetUsers(req request.UserQuery) (*response.UserListResponse, error)
}

type userService struct {
	userRepo    ports.UserRepository
	accountRepo ports.AccountRepository
	storage     ports.Storage
}

func NewUserService(userRepo ports.UserRepository, accountRepo ports.AccountRepository, storage ports.Storage) UserService {
	return &userService{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		storage:     storage,
	}
}

func (s *userService) CreateUser(req *request.CreateUserRequest, accountID uint) (*response.UserResponse, error) {
	if req == nil || req.FrontCardImage == nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "front_card_image is required", entities.ErrBadRequest)
	}

	if req.BackCardImage == nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "back_card_image is required", entities.ErrBadRequest)
	}

	if req.FaceImage == nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "face_image is required", entities.ErrBadRequest)
	}

	// 1. Pre-check for duplicate IDCard or AccountID
	existingByID, errID := s.userRepo.GetUserByIDCard(req.IDCard)
	if errID == nil && existingByID.Status != string(entities.UserStatusRejected) {
		return nil, entities.NewAppError(entities.ErrTypeConflict, "ID Card already registered and is not rejected", entities.ErrConflict)
	}

	existingByAcc, errAcc := s.userRepo.GetUserByAccountID(accountID)
	if errAcc == nil && existingByAcc.Status != string(entities.UserStatusRejected) {
		return nil, entities.NewAppError(entities.ErrTypeConflict, "Account already has a KYC record that is not rejected", entities.ErrConflict)
	}

	// 2. Upload images only if we are clear to proceed
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

	faceCardFile, err := req.FaceImage.Open()
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeBadRequest, "failed to open face_image", err)
	}
	defer faceCardFile.Close()

	faceCardBytes, err := io.ReadAll(faceCardFile)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to read face_image", err)
	}

	faceCardURL, err := s.storage.UploadFile(userFrontCardBucketName, buildFrontCardObjectPath(req.IDCard, req.FaceImage.Filename), faceCardBytes)
	if err != nil {
		return nil, entities.NewAppError(entities.ErrTypeInternal, "failed to upload face_image", err)
	}

	// 3. Create (which handles Upsert in the repository)
	var user *entities.User
	if errAcc == nil && existingByAcc.Status == string(entities.UserStatusRejected) {
		user = existingByAcc
		user.IDCard = req.IDCard
		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.FrontCardURL = frontCardURL
		user.BackCardURL = backCardURL
		user.FaceURL = faceCardURL
		user.BirthDate = req.BirthDate
		user.Status = string(entities.UserStatusPending)
		user.ExpireDate = req.ExpireDate
	} else {
		user = &entities.User{
			IDCard:       req.IDCard,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			AccountID:    accountID,
			FrontCardURL: frontCardURL,
			BackCardURL:  backCardURL,
			FaceURL:      faceCardURL,
			BirthDate:    req.BirthDate,
			Status:       string(entities.UserStatusPending),
			ExpireDate:   req.ExpireDate,
		}
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return response.FormatUserResponse(user), nil
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

		account, err := s.accountRepo.GetAccountByID(int(user.AccountID))
		if err != nil {
			return nil, err
		}

		account.KycToken = &user.KYCToken
		if err := s.accountRepo.UpdateAccount(context.Background(), int(user.AccountID), account); err != nil {
			return nil, err
		}
	}

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return response.FormatUserResponse(user), nil
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
		userResponses = append(userResponses, response.FormatUserResponse(user))
	}

	return &response.UserListResponse{
		Rows:       userResponses,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}
