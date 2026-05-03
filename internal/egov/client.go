// Package egov is a thin client for the e-Gov 法令API v2
// (https://laws.e-gov.go.jp/api/2/swagger-ui).
//
// It wraps the JSON endpoints (laws, law_revisions, law_data, keyword) and
// returns the response body verbatim as json.RawMessage so callers can decide
// how to present it.
package egov

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// BaseURL is the production e-Gov 法令API v2 root.
	BaseURL = "https://laws.e-gov.go.jp/api/2"

	defaultUserAgent = "egov-go-client/0.1.0"
	defaultTimeout   = 60 * time.Second
)

// Client calls the e-Gov 法令API v2.
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
}

// Options configures a Client. Zero values fall back to sensible defaults.
type Options struct {
	HTTPClient *http.Client
	BaseURL    string
	UserAgent  string
}

// NewClient returns a Client. opts may be nil.
func NewClient(opts *Options) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    BaseURL,
		userAgent:  defaultUserAgent,
	}
	if opts == nil {
		return c
	}
	if opts.HTTPClient != nil {
		c.httpClient = opts.HTTPClient
	}
	if opts.BaseURL != "" {
		c.baseURL = opts.BaseURL
	}
	if opts.UserAgent != "" {
		c.userAgent = opts.UserAgent
	}
	return c
}

// APIError represents a non-2xx response from the e-Gov API.
type APIError struct {
	Status int
	Body   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("e-Gov API HTTP %d: %s", e.Status, e.Body)
}

// get issues a GET request to baseURL+path with the given query string and
// returns the response body when the status is 2xx.
func (c *Client) get(ctx context.Context, path string, q *queryBuilder) (json.RawMessage, error) {
	q.set("response_format", "json")
	full := c.baseURL + path
	if encoded := q.encode(); encoded != "" {
		full += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, full, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{Status: resp.StatusCode, Body: string(body)}
	}
	if !json.Valid(body) {
		return nil, fmt.Errorf("invalid JSON from e-Gov API: %s", string(body))
	}
	return body, nil
}
