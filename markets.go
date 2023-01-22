package tradeorge

import (
	"context"
	"encoding/json"
)

type MarketInfo struct {
	Symbol       string
	Initialprice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
}

type AvailableMarketsService struct {
	c *Client
}

// Do send request
func (s *AvailableMarketsService) Do(ctx context.Context, opts ...RequestOption) ([]MarketInfo, error) {
	r := &request{
		method:   "GET",
		endpoint: "/markets",
		secType:  secTypeNone,
	}
	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	var res []map[string]MarketInfo

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	ri := make([]MarketInfo, 0, len(res))
	for _, i := range res {
		for symbol, info := range i {
			info.Symbol = symbol
			ri = append(ri, info)
		}
	}

	return ri, err
}
