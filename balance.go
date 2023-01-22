package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type CoinBalance struct {
	Balance   string `json:"balance"`
	Available string `json:"available"`
}

type CoinBalanceService struct {
	c *Client

	asset *string
}

// Do send request
func (s *CoinBalanceService) Do(ctx context.Context, opts ...RequestOption) (*CoinBalance, error) {
	if s.asset == nil {
		return nil, fmt.Errorf("asset not init")
	}

	r := &request{
		method:   "POST",
		endpoint: "/account/balance",
		secType:  secTypeAPIKey,
	}

	r.setFormParam("currency", string(*s.asset))

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
		CoinBalance
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf(res.Error)
	}

	return &res.CoinBalance, nil
}

func (s *CoinBalanceService) Asset(asset string) *CoinBalanceService {
	s.asset = &asset

	return s
}
