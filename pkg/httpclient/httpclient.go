package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type HttpClientInterface interface {
	Get(ctx context.Context, path string, options *RequestOptions) (*http.Response, error)
	Post(ctx context.Context, path string, jsonBody []byte, options *RequestOptions) (*http.Response, error)
	Put(ctx context.Context, path string, body any, options *RequestOptions) (*http.Response, error)
	Delete(ctx context.Context, path string, options *RequestOptions) (*http.Response, error)
	BuildGetRequestPath(queryParams url.Values) string
}

type HttpClient struct {
	client  *http.Client
	baseURL string
}

type RequestOptions struct {
	Headers map[string]string
	Timeout time.Duration
}

func NewHttpClient(baseURL string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		client:  &http.Client{Timeout: timeout},
		baseURL: baseURL,
	}
}

// GET-запрос
func (c *HttpClient) Get(ctx context.Context, path string, options *RequestOptions) (*http.Response, error) {
	url := c.BuildURL(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	c.applyOptions(req, options)

	return c.client.Do(req)
}

// POST-запрос с JSON телом
func (c *HttpClient) Post(ctx context.Context, path string, jsonBody []byte, options *RequestOptions) (*http.Response, error) {
	url := c.BuildURL(path)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.applyOptions(req, options)

	return c.client.Do(req)
}

// PUT-запрос с JSON телом
func (c *HttpClient) Put(ctx context.Context, path string, body any, options *RequestOptions) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON body: %w", err)
	}

	url := c.BuildURL(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.applyOptions(req, options)

	return c.client.Do(req)
}

// DELETE-запрос
func (c *HttpClient) Delete(ctx context.Context, path string, options *RequestOptions) (*http.Response, error) {
	url := c.BuildURL(path)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	c.applyOptions(req, options)

	return c.client.Do(req)
}

func (c *HttpClient) BuildQueryPath(queryParams url.Values) string {
	return fmt.Sprintf("?%s", queryParams.Encode())
}

// Построение полного URL
func (c *HttpClient) BuildURL(path string) string {
	return fmt.Sprintf("%s%s", c.baseURL, path)
}

// Применение опций к запросу
func (c *HttpClient) applyOptions(req *http.Request, options *RequestOptions) {
	if options == nil {
		return
	}

	if options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	if options.Timeout > 0 {
		c.client.Timeout = options.Timeout
	}
}
