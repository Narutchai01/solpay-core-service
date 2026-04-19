package request

import "encoding/json"

type OffChainRequest struct {
	CategoryID int    `json:"category_id" validate:"number default=1"`
	QuoteID    string `json:"quote_id" validate:"required"`
}

func (r *OffChainRequest) UnmarshalJSON(data []byte) error {
	type payload struct {
		QuoteID         string `json:"quote_id"`
		QuoteIDCamel    string `json:"quoteID"`
		CategoryID      *int   `json:"category_id"`
		CategoryIDCamel *int   `json:"categoryID"`
	}

	var p payload
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	if p.QuoteID == "" {
		p.QuoteID = p.QuoteIDCamel
	}

	r.QuoteID = p.QuoteID

	categoryID := p.CategoryID
	if categoryID == nil {
		categoryID = p.CategoryIDCamel
	}

	if categoryID == nil || *categoryID == 0 {
		r.CategoryID = 1
	} else {
		r.CategoryID = *categoryID
	}

	return nil
}
