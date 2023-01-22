package tradeorge

import (
	"context"
	"encoding/json"
	"fmt"
)

type CancelOrderService struct {
	c *Client

	uuid *string
}

// Do send request
func (s *CancelOrderService) Do(ctx context.Context, opts ...RequestOption) error {
	if s.uuid == nil {
		return fmt.Errorf("uuid not init")
	}

	r := &request{
		method:   "POST",
		endpoint: "/order/cancel",
		secType:  secTypeAPIKey,
	}

	r.setFormParam("uuid", string(*s.uuid))

	data, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return err
	}

	res := struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}{}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	if !res.Success {
		return fmt.Errorf(res.Error)
	}

	return nil
}

func (s *CancelOrderService) Uuid(uuid string) *CancelOrderService {
	s.uuid = &uuid

	return s
}
