package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type Ticker struct {
	Initialprice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
}

type TickerService struct {
	c *Client

	symbol *string
}

// Do send request
func (s *TickerService) Do(ctx context.Context, opts ...RequestOption) (*Ticker, error) {
	if s.symbol == nil {
		return nil, fmt.Errorf("symbol not init")
	}

	r := &request{
		method:   "GET",
		endpoint: fmt.Sprintf("/ticker/%s", *s.symbol),
		secType:  secTypeNone,
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Ticker
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if res.Success != true {
		return nil, fmt.Errorf(res.Error)
	}

	return &res.Ticker, nil
}

func (s *TickerService) Symbol(symbol string) *TickerService {
	s.symbol = &symbol

	return s
}
