package request

import "encoding/json"

type TopUpRequest struct {
	QuoteID     string   `json:"quote_id" validate:"required"`
	TxHash      string   `json:"tx_hash" validate:"required"`
	MaxSlippage *float64 `json:"max_slippage,omitempty"`
	CategoryID  int      `json:"category_id,omitempty"`
}

// DefaultSlippage is the default value for MaxSlippage if not provided.
const DefaultSlippage = 0.01

func (r *TopUpRequest) SetDefaultSlippage() {
	if r.MaxSlippage == nil {
		r.MaxSlippage = new(float64)
		*r.MaxSlippage = DefaultSlippage
	}
}

func (r *TopUpRequest) UnmarshalJSON(data []byte) error {
	type payload struct {
		QuoteID         string   `json:"quote_id"`
		QuoteIDCamel    string   `json:"quoteID"`
		TxHash          string   `json:"tx_hash"`
		TxHashCamel     string   `json:"txHash"`
		MaxSlippage     *float64 `json:"max_slippage"`
		MaxSlipCamel    *float64 `json:"maxSlippage"`
		CategoryID      *int     `json:"category_id"`
		CategoryIDCamel *int     `json:"categoryID"`
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
	if p.MaxSlippage == nil {
		p.MaxSlippage = p.MaxSlipCamel
	}

	categoryID := p.CategoryID
	if categoryID == nil {
		categoryID = p.CategoryIDCamel
	}

	r.QuoteID = p.QuoteID
	r.TxHash = p.TxHash
	r.MaxSlippage = p.MaxSlippage
	if categoryID == nil || *categoryID == 0 {
		r.CategoryID = 1
	} else {
		r.CategoryID = *categoryID
	}

	return nil
}
