package ports

import "github.com/Narutchai01/solpay-core-service/internal/entities"

type Publisher interface {
	Publish(queueName string, message []byte) error
	PublishTransactionMessage(queue string, tx *entities.TransactionEntity, status string)
}
