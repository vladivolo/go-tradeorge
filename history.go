package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type Trade struct {
	Date     int64  `json:"date"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type HistoryTradesService struct {
	c *Client

	symbol *string
}

// Do send request
func (s *HistoryTradesService) Do(ctx context.Context, opts ...RequestOption) ([]Trade, error) {
	if s.symbol == nil {
		return nil, fmt.Errorf("symbol not init")
	}

	r := &request{
		method:   "GET",
		endpoint: fmt.Sprintf("/history/%s", *s.symbol),
		secType:  secTypeNone,
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := []Trade{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *HistoryTradesService) Symbol(symbol string) *HistoryTradesService {
	s.symbol = &symbol

	return s
}
