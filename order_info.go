package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type Info struct {
	Date      int64  `json:"date"`
	Type      string `json:"type"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Market    string `json:"market"`
	Fulfilled string `json:"fulfilled"`
}

type OrderInfoService struct {
	c *Client

	uuid *string
}

// Do send request
func (s *OrderInfoService) Do(ctx context.Context, opts ...RequestOption) (*Info, error) {
	if s.uuid == nil {
		return nil, fmt.Errorf("uuid not init")
	}

	r := &request{
		method:   "GET",
		endpoint: fmt.Sprintf("/account/order/%s", *s.uuid),
		secType:  secTypeAPIKey,
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		Info
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf(res.Error)
	}

	return &res.Info, nil
}

func (s *OrderInfoService) Uuid(uuid string) *OrderInfoService {
	s.uuid = &uuid

	return s
}
