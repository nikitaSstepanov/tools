package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	prefix string
}

func New(prefixParts ...string) *Client {
	prefix := ""
	
	if len(prefixParts) != 0 {
		prefix = strings.Join(prefixParts, "/")
	}

	return &Client{
		prefix: prefix,
	}
}

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

func (c *Client) GetWithBind(url string, to interface{}) error {
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

func (c *Client) PostWithBind(url string, contentType string, data interface{}, to interface{}) error {
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

