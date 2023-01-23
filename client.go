package tradeorge

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	secTypeNone secType = iota
	secTypeAPIKey
)

type doFunc func(req *http.Request) (*http.Response, error)

// Client define API client
type Client struct {
	APIKey     string
	APISecret  string
	BaseURL    string
	UserAgent  string
	HTTPClient *http.Client
	Debug      bool
	Logger     *log.Logger

	do doFunc
}

// NewClient initialize an API client instance with API key and secret key.
// You should always call this function before using this SDK.
// Services will be created by the form client.NewXXXService().
func NewClient(apiKey string, apiSecret string) *Client {
	return &Client{
		APIKey:     apiKey,
		APISecret:  apiSecret,
		BaseURL:    "https://tradeogre.com/api/v1",
		UserAgent:  "go-tradeorge/golang",
		HTTPClient: http.DefaultClient,
		Debug:      false,
		Logger:     log.New(os.Stderr, "go-tradeorge ", log.LstdFlags),
	}
}

func (c *Client) debug(format string, v ...interface{}) {
	if c.Debug {
		c.Logger.Printf(format, v...)
	}
}

func (c *Client) parseRequest(r *request, opts ...RequestOption) (err error) {
	// set request options from user
	for _, opt := range opts {
		opt(r)
	}
	err = r.validate()
	if err != nil {
		return err
	}

	fullURL := fmt.Sprintf("%s%s", c.BaseURL, r.endpoint)
	queryString := r.query.Encode()
	body := &bytes.Buffer{}
	bodyString := r.form.Encode()
	header := http.Header{}

	header.Set("accept", "application/json")

	if r.secType == secTypeAPIKey {
		header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.APIKey+":"+c.APISecret)))
	}

	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}

	if queryString != "" {
		fullURL = fmt.Sprintf("%s?%s", fullURL, queryString)
	}
	c.debug("full url: %s, body: %s", fullURL, bodyString)

	r.fullURL = fullURL
	r.header = header
	r.body = body
	return nil
}

func (c *Client) callAPI(ctx context.Context, r *request, opts ...RequestOption) (data []byte, err error) {
	err = c.parseRequest(r, opts...)
	if err != nil {
		return []byte{}, err
	}
	req, err := http.NewRequest(r.method, r.fullURL, r.body)
	if err != nil {
		return []byte{}, err
	}
	req = req.WithContext(ctx)
	req.Header = r.header
	c.debug("request: %#v", req)
	f := c.do
	if f == nil {
		f = c.HTTPClient.Do
	}
	res, err := f(req)
	if err != nil {
		return []byte{}, err
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		cerr := res.Body.Close()
		// Only overwrite the retured error if the original error was nil and an
		// error occurred while closing the body.
		if err == nil && cerr != nil {
			err = cerr
		}
	}()
	c.debug("response: %#v", res)
	c.debug("response body: %s", string(data))
	c.debug("response status code: %d", res.StatusCode)

	if res.StatusCode >= 400 {
		apiErr := new(APIError)
		e := json.Unmarshal(data, apiErr)
		if e != nil {
			c.debug("failed to unmarshal json: %s", e)
		}
		return nil, apiErr
	}
	return data, nil
}

// Get list of all avialable markets
func (c *Client) NewAvailableMarketsService() *AvailableMarketsService {
	return &AvailableMarketsService{c: c}
}

// Get orders book of market
func (c *Client) NewOrderBookService() *OrderBookService {
	return &OrderBookService{c: c}
}

// Get ticker of market
func (c *Client) NewTickerService() *TickerService {
	return &TickerService{c: c}
}

// Get history of market
func (c *Client) NewHistoryTradesService() *HistoryTradesService {
	return &HistoryTradesService{c: c}
}

// Get active orders of market
func (c *Client) NewActiveOrdersService() *ActiveOrdersService {
	return &ActiveOrdersService{c: c}
}

// Get order info
func (c *Client) NewOrderInfoService() *OrderInfoService {
	return &OrderInfoService{c: c}
}

// Get Asset balance
func (c *Client) NewCoinBalanceService() *CoinBalanceService {
	return &CoinBalanceService{c: c}
}

// Get Assets balances
func (c *Client) NewBalancesService() *BalancesService {
	return &BalancesService{c: c}
}

// Create new order
func (c *Client) NewCreateOrderService() *CreateOrderService {
	return &CreateOrderService{c: c}
}

// Cancel order
func (c *Client) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c}
}
