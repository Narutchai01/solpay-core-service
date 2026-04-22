package request

import "mime/multipart"

type CreateUserRequest struct {
	IDCard         string                `json:"id_card" form:"id_card" validate:"required" binding:"required"`
	FirstName      string                `json:"first_name" form:"first_name" validate:"required" binding:"required"`
	LastName       string                `json:"last_name" form:"last_name" validate:"required" binding:"required"`
	BirthDate      string                `json:"birth_date" form:"birth_date" validate:"required" binding:"required"`
	ExpireDate     string                `json:"expire_date" form:"expire_date" validate:"required" binding:"required"`
	FrontCardImage *multipart.FileHeader `json:"-" form:"front_card_image" validate:"required"`
}
