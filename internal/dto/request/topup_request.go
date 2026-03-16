package request

import "encoding/json"

type TopUpRequest struct {
	QuoteID string `json:"quote_id" validate:"required"`
	TxHash  string `json:"tx_hash" validate:"required"`
}

func (r *TopUpRequest) UnmarshalJSON(data []byte) error {
	type payload struct {
		QuoteID      string `json:"quote_id"`
		QuoteIDCamel string `json:"quoteID"`
		TxHash       string `json:"tx_hash"`
		TxHashCamel  string `json:"txHash"`
	}

	var p payload
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	if p.QuoteID == "" {
		p.QuoteID = p.QuoteIDCamel
	}
	if p.TxHash == "" {
		p.TxHash = p.TxHashCamel
	}

	r.QuoteID = p.QuoteID
	r.TxHash = p.TxHash

	return nil
}
