package repositories

import (
	"bytes"
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
	u, err := url.Parse(r.baseURL + "/quote")
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

func (r *swapRepository) BuildSwapUnsignedTransaction(req request.SwapUnsignedTransactionRequest, walletAddress string) (response.SwapUnsignedTransactionFullResponse, error) {
	u, err := url.Parse(r.baseURL + "/swap")
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("parse url: %w", err)
	}

	payload := map[string]any{
		"wallet":    walletAddress,
		"poolId":    req.PoolID,
		"amountIn":  req.AmountIn,
		"slippage":  req.Slippage,
		"inputMint": req.InputMint,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("marshal body: %w", err)
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var swapResp response.SwapUnsignedTransactionFullResponse
	if err := json.Unmarshal(body, &swapResp); err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return swapResp, nil
}

func (r *swapRepository) BuildSwapInstruction(req request.SwapUnsignedTransactionRequest, walletAddress string) (response.SwapUnsignedTransactionFullResponse, error) {
	u, err := url.Parse(r.baseURL + "/swap/build")
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("parse url: %w", err)
	}

	payload := map[string]any{
		"wallet":    walletAddress,
		"poolId":    req.PoolID,
		"amountIn":  req.AmountIn,
		"slippage":  req.Slippage,
		"inputMint": req.InputMint,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("marshal body: %w", err)
	}

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var swapResp response.SwapUnsignedTransactionFullResponse
	if err := json.Unmarshal(body, &swapResp); err != nil {
		return response.SwapUnsignedTransactionFullResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return swapResp, nil
}
