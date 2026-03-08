package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type PaymentRepository interface {
	CreateRecipient(recipent *entities.Recipient) error
	GetRecipentByNumber(number string) (entities.Recipient, error)
}
