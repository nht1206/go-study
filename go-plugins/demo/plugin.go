package demo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/nht1206/go-study/go-plugins/core"
)

type PluginCallerConfig struct {
	Endpoint string
	Timeout  time.Duration
}

type PluginCaller struct {
	endpoint   *url.URL
	httpClient *http.Client
}

func NewPluginCaller(config *PluginCallerConfig) (*PluginCaller, error) {
	if len(config.Endpoint) <= 0 {
		return nil, errors.New("missing plugin endpoint")
	}

	u, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, err
	}

	c := &PluginCaller{
		endpoint: u,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}

	return c, nil
}

func (c *PluginCaller) do(ctx context.Context, req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", core.DefaultContentType)
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return json.NewDecoder(res.Body).Decode(v)
}

func (p *PluginCaller) Call(ctx context.Context, plgReq *Request) (*Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(plgReq); err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, p.endpoint.String(), &buf)
	if err != nil {
		return nil, err
	}

	plgRes := Response{}
	err = p.do(ctx, r, &plgRes)
	if err != nil {
		return nil, err
	}

	return &plgRes, nil
}
