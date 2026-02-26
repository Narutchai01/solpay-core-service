package models

type MetaDataSolana struct {
	Transaction_type string `json:"transaction_type"`
	Amount_THB       int    `json:"amount_thb"`
	Amount_USDC      int    `json:"amount_usdc"`
	AccountID        int    `json:"account_id"`
}

type SolanaTxMessage struct {
	TxID     string         `json:"tx_id"`
	Base64Tx string         `json:"base64_tx"`
	MetaData MetaDataSolana `json:"metadata"`
}
