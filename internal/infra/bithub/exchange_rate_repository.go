package bithub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Narutchai01/solpay-core-service/internal/core/ports"
	"github.com/Narutchai01/solpay-core-service/internal/entities"
)

type exchangeRateRepository struct {
	baseURL string
}

func NewExchangeRateRepository(baseURL string) ports.ExchangeRepository {
	return &exchangeRateRepository{
		baseURL: baseURL,
	}
}

func (r *exchangeRateRepository) GetExchangeRate(symbol string) (*[]entities.ExchangeRate, error) {
	resp, err := http.Get(r.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var allRates []entities.ExchangeRate
	if err := json.Unmarshal(body, &allRates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exchange rate data: %w", err)
	}

	// Filter by symbol if provided
	if symbol != "" {
		upperSymbol := strings.ToUpper(symbol)
		filtered := make([]entities.ExchangeRate, 0)
		for _, rate := range allRates {
			if strings.ToUpper(rate.Symbol) == upperSymbol {
				filtered = append(filtered, rate)
			}
		}
		return &filtered, nil
	}

	return &allRates, nil
}
