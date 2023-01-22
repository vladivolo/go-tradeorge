package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type Balances map[string]string

type BalancesService struct {
	c *Client
}

// Do send request
func (s *BalancesService) Do(ctx context.Context, opts ...RequestOption) (Balances, error) {
	r := &request{
		method:   "GET",
		endpoint: "/account/balances",
		secType:  secTypeAPIKey,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	var res struct {
		Success  bool   `json:"success"`
		Error    string `json:"error"`
		Balances `json:"balances"`
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf(res.Error)
	}

	return res.Balances, nil
}
