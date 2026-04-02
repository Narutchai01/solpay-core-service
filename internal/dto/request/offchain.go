package request

import "encoding/json"

type OffChainRequest struct {
	QuoteID string `json:"quote_id" validate:"required"`
}

func (r *OffChainRequest) UnmarshalJSON(data []byte) error {
	type payload struct {
		QuoteID      string `json:"quote_id"`
		QuoteIDCamel string `json:"quoteID"`
	}

	var p payload
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	if p.QuoteID == "" {
		p.QuoteID = p.QuoteIDCamel
	}

	r.QuoteID = p.QuoteID

	return nil
}
