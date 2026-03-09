package repositories

import (
	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewGormPaymentRepository(db *gorm.DB) ports.PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) CreateRecipient(recipient *entities.Recipient) error {
	return r.db.Create(recipient).Error
}

func (r *PaymentRepository) GetRecipentByNumber(number string) (entities.Recipient, error) {
	recipent := entities.Recipient{}
	if err := r.db.Where("number = ?", number).First(&recipent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entities.Recipient{}, entities.ErrNotFound
		}
		return entities.Recipient{}, err
	}
	return recipent, nil
}

func (r *PaymentRepository) CreateLogPayment(payment *entities.LogPayment) error {
	return r.db.Create(payment).Error
}
