package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/dto/request"
	"github.com/Narutchai01/solpay-core-service/internal/dto/response"
)

type swapRepository struct {
	baseURL string
}

func NewSwapRepository(baseURL string) ports.SwapRepository {
	return &swapRepository{baseURL: baseURL}
}

func (r *swapRepository) GetSwapQuote(query request.SwapQuoteRequest) (response.SwapQuoteFullResponse, error) {
	// The path was previously /api/v1/quote but it returned 404.
	// Changing to /quote based on common conventions for local services.
	u, err := url.Parse(r.baseURL + "/api/v1/quote")
	if err != nil {
		return response.SwapQuoteFullResponse{}, fmt.Errorf("parse url: %w", err)
	}

	params := url.Values{}
	if query.Slippage != "" {
		params.Add("slippage", query.Slippage)
	}
	if query.AmountIn != "" {
		params.Add("amountIn", query.AmountIn)
	}
	u.RawQuery = params.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return response.SwapQuoteFullResponse{}, fmt.Errorf("http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response.SwapQuoteFullResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.SwapQuoteFullResponse{}, fmt.Errorf("read body: %w", err)
	}

	var quoteResp response.SwapQuoteFullResponse
	if err := json.Unmarshal(body, &quoteResp); err != nil {
		return response.SwapQuoteFullResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return quoteResp, nil
}
