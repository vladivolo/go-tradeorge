package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type OrderInfo struct {
	Uuid     string `json:"uuid"`
	Date     int64  `json:"date"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	Market   string `json:"market"`
}

type ActiveOrders []OrderInfo

type ActiveOrdersService struct {
	c *Client

	symbol *string
}

// Do send request
func (s *ActiveOrdersService) Do(ctx context.Context, opts ...RequestOption) ([]OrderInfo, error) {
	if s.symbol == nil {
		return nil, fmt.Errorf("symbol not init")
	}

	r := &request{
		method:   "POST",
		endpoint: "/account/orders",
		secType:  secTypeAPIKey,
	}

	r.setFormParam("market", string(*s.symbol))

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := ActiveOrders{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *ActiveOrdersService) Symbol(symbol string) *ActiveOrdersService {
	s.symbol = &symbol

	return s
}
