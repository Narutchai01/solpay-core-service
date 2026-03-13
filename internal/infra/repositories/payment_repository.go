package repositories

import (
	"errors"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

// NewGormPaymentRepository creates a new GORM-backed PaymentRepository.
func NewGormPaymentRepository(database *gorm.DB) ports.PaymentRepository {
	return &paymentRepository{db: database}
}

func (r *paymentRepository) CreateRecipient(recipient *entities.Recipient) error {
	return r.db.Create(recipient).Error
}

func (r *paymentRepository) GetRecipientByNumber(number string) (entities.Recipient, error) {
	var recipient entities.Recipient
	if err := r.db.Where("number = ?", number).First(&recipient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Recipient{}, entities.ErrNotFound
		}
		return entities.Recipient{}, err
	}
	return recipient, nil
}

func (r *paymentRepository) CreateLogPayment(payment *entities.LogPayment) error {
	return r.db.Create(payment).Error
}
