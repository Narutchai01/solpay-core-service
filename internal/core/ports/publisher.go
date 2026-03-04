package ports

type Publisher interface {
	Publish(queueName string, message []byte) error
}
