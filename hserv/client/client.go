package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Config struct {
	Prefix  string        `yaml:"prefix" env:"HTTP_CLIENT_PREFIX" env-default:""`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_CLIENT_TIMEOUT" env-default:"5s"`
}

type Client struct {
	prefix string
	client *http.Client
}

func New(cfg *Config) *Client {
	return &Client{
		prefix: cfg.Prefix,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// "OPTIONS"                ; Section 9.2
//| "GET"                    ; Section 9.3
//| "HEAD"                   ; Section 9.4
//| "POST"                   ; Section 9.5
//| "PUT"                    ; Section 9.6
//| "DELETE"                 ; Section 9.7
//| "TRACE"                  ; Section 9.8
//| "CONNECT"

func (c *Client) Get(url string) ([]byte, error) {
	if c.prefix != "" {
		url = c.prefix + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) GetWithJsonBind(url string, to interface{}) error {
	body, err := c.Get(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, to)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Post(url string, contentType string, data interface{}) ([]byte, error) {
	if c.prefix != "" {
		url = c.prefix + url
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(body)

	resp, err := http.Post(url, contentType, reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) PostWithJsonBind(url string, contentType string, data interface{}, to interface{}) error {
	body, err := c.Post(url, contentType, data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, to)
	if err != nil {
		return err
	}

	return nil
}
