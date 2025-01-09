package httper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type ClientCfg struct {
	Prefix  string        `yaml:"prefix" env:"HTTP_CLIENT_PREFIX" env-default:""`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_CLIENT_TIMEOUT" env-default:"5s"`
}

type Client struct {
	prefix string
	client *http.Client
}

func NewClient(cfg *ClientCfg) *Client {
	return &Client{
		prefix: cfg.Prefix,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *Client) Get(url string) (*Resp, error) {
	if c.prefix != "" {
		url = c.prefix + url
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return newResp(resp, body), nil
}

func (c *Client) GetJson(url string, to interface{}) (*Resp, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.ByteBody, to)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) PostWithJson(url string, data interface{}) (*Resp, error) {
	if c.prefix != "" {
		url = c.prefix + url
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	resp, err := c.client.Post(url, string(JsonType), reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return newResp(resp, body), nil
}

func (c *Client) PostWithJsonBind(url string, data interface{}, to interface{}) (*Resp, error) {
	resp, err := c.PostWithJson(url, data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.ByteBody, to)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) Do(req *Req) (*Resp, error) {
	if c.prefix != "" {
		newUrl, err := url.Parse(c.prefix + fmt.Sprintf("%v", req.URL))
		if err != nil {
			return nil, err
		}

		req.URL = newUrl
	}

	resp, err := c.client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if req.NeedUnmarshal {
		req.Unmarshal(body)
	}

	return newResp(resp, body), nil
}
