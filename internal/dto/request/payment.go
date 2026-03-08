package request

type CreateRecipient struct {
	Name   string `json:"name" validate:"required"`
	Number string `json:"number" validate:"required"`
}
