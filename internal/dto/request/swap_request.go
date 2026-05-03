package request

type SwapQuoteRequest struct {
	Slippage string `json:"slippage" query:"slippage"`
	AmountIn string `json:"amountIn" query:"amountIn"`
}
