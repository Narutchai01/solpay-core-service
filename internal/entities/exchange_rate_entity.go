package entities

type ExchangeRate struct {
	Symbol        string `json:"symbol"`
	BaseVolume    string `json:"base_volume"`
	High24_Hr     string `json:"high_24_hr"`
	HighestBid    string `json:"highest_bid"`
	Last          string `json:"last"`
	Low24_Hr      string `json:"low_24_hr"`
	LowestAsk     string `json:"lowest_ask"`
	PercentChange string `json:"percent_change"`
	QuoteVolume   string `json:"quote_volume"`
}
