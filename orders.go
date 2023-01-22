package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

type origOrderBook struct {
	Success string            `json:"success"`
	Error   string            `json:"error"`
	Buy     map[string]string `json:"buy"`
	Sell    map[string]string `json:"sell"`
}

type Order struct {
	Price  float64
	Volume float64
}

type Orders []Order

type OrderBook struct {
	Buy  []Order
	Sell []Order
}

type OrderBookService struct {
	c *Client

	symbol *string
}

// Do send request
func (s *OrderBookService) Do(ctx context.Context, opts ...RequestOption) (*OrderBook, error) {
	if s.symbol == nil {
		return nil, fmt.Errorf("symbol not init")
	}

	r := &request{
		method:   "GET",
		endpoint: fmt.Sprintf("/orders/%s", *s.symbol),
		secType:  secTypeNone,
	}

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return nil, err
	}

	res := origOrderBook{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if res.Success != "true" {
		return nil, fmt.Errorf(res.Error)
	}

	ob := OrderBook{
		Buy:  make([]Order, 0, len(res.Buy)),
		Sell: make([]Order, 0, len(res.Sell)),
	}

	for price_str, volume_str := range res.Buy {
		price, err := strconv.ParseFloat(price_str, 64)
		if err != nil {
			continue
		}

		volume, err := strconv.ParseFloat(volume_str, 64)
		if err != nil {
			continue
		}

		ob.Buy = append(ob.Buy, Order{
			Price:  price,
			Volume: volume,
		})
	}

	sort.Slice(ob.Buy, func(i, j int) bool { return ob.Buy[i].Price > ob.Buy[j].Price })

	for price_str, volume_str := range res.Sell {
		price, err := strconv.ParseFloat(price_str, 64)
		if err != nil {
			continue
		}

		volume, err := strconv.ParseFloat(volume_str, 64)
		if err != nil {
			continue
		}

		ob.Sell = append(ob.Sell, Order{
			Price:  price,
			Volume: volume,
		})
	}

	sort.Slice(ob.Sell, func(i, j int) bool { return ob.Sell[i].Price < ob.Sell[j].Price })

	return &ob, err
}

func (s *OrderBookService) Symbol(symbol string) *OrderBookService {
	s.symbol = &symbol

	return s
}
