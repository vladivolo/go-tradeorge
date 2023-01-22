package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type CreateOrderResponse struct {
	Uuid         string `json:""`
	Bnewbalavail string `json:"bnewbalavail"`
	Snewbalavail string `json:"snewbalavail"`
}

type CreateOrderService struct {
	c *Client

	ops      *string
	market   *string
	quantity *float64
	price    *float64
}

// Do send request
func (s *CreateOrderService) Do(ctx context.Context, opts ...RequestOption) (*CreateOrderResponse, error) {
	if s.ops == nil {
		return nil, fmt.Errorf("ops not init")
	}

	if s.market == nil {
		return nil, fmt.Errorf("market not init")
	}

	if s.quantity == nil {
		return nil, fmt.Errorf("quantity not init")
	}

	if s.price == nil {
		return nil, fmt.Errorf("price not init")
	}

	r := &request{
		method:   "POST",
		endpoint: fmt.Sprintf("/order/%s", *s.ops),
		secType:  secTypeAPIKey,
	}

	r.setFormParam("market", string(*s.market))
	r.setFormParam("price", fmt.Sprintf("%.8f", *s.price))
	r.setFormParam("quantity", fmt.Sprintf("%.8f", *s.quantity))

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		CreateOrderResponse
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf(res.Error)
	}

	return &res.CreateOrderResponse, nil
}

func (s *CreateOrderService) Operation(ops string) *CreateOrderService {
	s.ops = &ops

	return s
}

func (s *CreateOrderService) Market(market string) *CreateOrderService {
	s.market = &market

	return s
}

func (s *CreateOrderService) Quantity(quantity float64) *CreateOrderService {
	s.quantity = &quantity

	return s
}

func (s *CreateOrderService) Price(price float64) *CreateOrderService {
	s.price = &price

	return s
}
