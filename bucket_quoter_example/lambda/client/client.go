package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const (
	DefaultHTTPHost = "https://localhost:8443"
)

type Client struct {
	httpc *resty.Client
}

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func NewClientWithResty(httpc *resty.Client, opts ...ClientOpt) (*Client, error) {
	c := &Client{
		httpc: httpc,
	}
	for _, o := range opts {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	if c.httpc.HostURL == "" {
		c.httpc.SetBaseURL(DefaultHTTPHost)
	}
	c.httpc.SetHeader("Host", c.httpc.HostURL)

	c.httpc.SetDoNotParseResponse(true)

	return c, nil
}

func NewClient(opts ...ClientOpt) (*Client, error) {
	return NewClientWithResty(resty.New(), opts...)
}

func (c Client) GetLimitResponse(ctx context.Context) (*Response, error) {
	req := c.httpc.R()

	resp, err := req.SetContext(ctx).Get("/limiter")
	if err != nil {
		return nil, fmt.Errorf("lambda: %w", err)
	}

	// read all and close body for proper Keep-Alive connection reuse
	defer func() {
		_, _ = io.Copy(ioutil.Discard, resp.RawBody())
		_ = resp.RawBody().Close()
	}()

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("laas: bad HTTP status code %d", resp.StatusCode())
	}

	var result Response

	dec := json.NewDecoder(resp.RawBody())
	if err := dec.Decode(&result); err != nil {
		return nil, fmt.Errorf("lambda: %w", err)
	}

	return &result, nil
}
