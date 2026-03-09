package models

// SolanaMetaData holds transaction metadata for the Solana blockchain worker.
type SolanaMetaData struct {
	TransactionType string `json:"transaction_type"`
	AmountTHB       int    `json:"amount_thb"`
	AmountUSDC      int    `json:"amount_usdc"`
	AccountID       int    `json:"account_id"`
}

// SolanaTxMessage is the message sent to the Solana worker queue.
type SolanaTxMessage struct {
	TxID     string         `json:"tx_id"`
	Base64Tx string         `json:"base64_tx"`
	MetaData SolanaMetaData `json:"metadata"`
}
