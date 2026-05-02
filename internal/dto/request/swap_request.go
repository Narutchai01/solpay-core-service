package request

type SwapQuoteRequest struct {
	Slippage string `json:"slippage" query:"slippage"`
	AmountIn string `json:"amountIn" query:"amountIn"`
}

type SwapUnsignedTransactionRequest struct {
	Slippage  float64 `json:"slippage" query:"slippage"`
	AmountIn  string  `json:"amountIn" query:"amountIn"`
	PoolID    string  `json:"poolId" query:"poolId"`
	InputMint string  `json:"inputMint" query:"inputMint"`
}
