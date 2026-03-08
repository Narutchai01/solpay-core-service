package ports

type WebSocketPort interface {
	NotifyTransactionStatus(txID string, payload []byte)
}
