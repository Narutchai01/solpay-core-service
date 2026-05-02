package response

type SwapQuoteFullResponse struct {
	Status string        `json:"status"`
	Data   SwapQuoteData `json:"data"`
}

type SwapQuoteData struct {
	PoolID            string            `json:"poolId"`
	InputMint         string            `json:"inputMint"`
	OutputMint        string            `json:"outputMint"`
	Slippage          float64           `json:"slippage"`
	AmountInRequested string            `json:"amountInRequested"`
	CurrentPrice      string            `json:"currentPrice"`
	ExecutionPrice    string            `json:"executionPrice"`
	PriceImpact       string            `json:"priceImpact"`
	ExecutionPriceX64 string            `json:"executionPriceX64"`
	RemainingAccounts []string          `json:"remainingAccounts"`
	RealAmountIn      SwapAmountDetails `json:"realAmountIn"`
	AmountOut         SwapAmountDetails `json:"amountOut"`
	MinAmountOut      SwapAmountDetails `json:"minAmountOut"`
	Fee               SwapAmountDetails `json:"fee"`
}

type SwapAmountDetails struct {
	RawAmount     string `json:"rawAmount"`
	DecimalAmount string `json:"decimalAmount"`
}

type SwapUnsignedTransactionFullResponse struct {
	Status string                      `json:"status"`
	Data   SwapUnsignedTransactionData `json:"data"`
}

type SwapUnsignedTransactionData struct {
	TxID                 string `json:"txId"`
	Transaction          string `json:"transaction"`
	Blockhash            string `json:"blockhash"`
	LastValidBlockHeight int64  `json:"lastValidBlockHeight"`
}
