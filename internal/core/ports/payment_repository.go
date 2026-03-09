package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

// PaymentRepository defines operations for payment-related persistence.
type PaymentRepository interface {
	CreateRecipient(recipient *entities.Recipient) error
	GetRecipientByNumber(number string) (entities.Recipient, error)
	CreateLogPayment(payment *entities.LogPayment) error
}
