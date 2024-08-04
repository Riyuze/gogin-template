package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gogin-template/baselib/dto"
	"gogin-template/bootstrap"
	"net/http"
)

type SampleClient struct {
	url      string
	endpoint string
	client   *bootstrap.HttpClient
	cfg      *bootstrap.Container
}

func NewSampleClient(url string, endpoint string, cfg *bootstrap.Container) *SampleClient {
	client := bootstrap.NewHttpClient(cfg, cfg.GetTracer())
	return &SampleClient{
		url:      url,
		endpoint: endpoint,
		client:   client,
		cfg:      cfg,
	}
}

func (c *SampleClient) Post(ctx context.Context, auth string, clientId string,
	cookiesDomain string, cookiesPrefix string) (*dto.Response[bool], error) {

	data := map[string]interface{}{
		"ClientId":      clientId,
		"CookiesDomain": cookiesDomain,
		"CookiesPrefix": cookiesPrefix,
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	addHeaderFunc := func(r *http.Request) {
		//* Do something with the request object
		r.Header.Add("Authorization", auth)
		r.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.client.Post(ctx, fmt.Sprintf("%s%s", c.url, c.endpoint), bytes.NewReader(jsonBytes), addHeaderFunc)

	if err != nil {
		return nil, err
	}

	var res dto.Response[bool]
	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return nil, err
	}

	return &res, err
}
