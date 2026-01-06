package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

type HTTPClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
	debug      bool
	debugOut   io.Writer
}

type Option func(*HTTPClient)

func WithDebug(w io.Writer) Option {
	return func(c *HTTPClient) {
		c.debug = true
		c.debugOut = w
	}
}

func WithBaseURL(url string) Option {
	return func(c *HTTPClient) {
		c.baseURL = url
	}
}

func NewHTTPClient(cfg *config.Config, opts ...Option) *HTTPClient {
	c := &HTTPClient{
		baseURL: BaseURL,
		token:   cfg.AccessToken,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *HTTPClient) GetMe(ctx context.Context) (*models.User, error) {
	var response struct {
		Data models.User `json:"data"`
	}

	if err := c.get(ctx, "/users/me", &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) get(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

func (c *HTTPClient) do(ctx context.Context, method, path string, body io.Reader, result any) error {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.NewGeneralError("failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.debug && c.debugOut != nil {
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] %s %s\n", method, url)
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Authorization: Bearer %s...\n", truncateToken(c.token))
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.NewNetworkError("request failed", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if c.debug && c.debugOut != nil {
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Response: %d %s (%s)\n", resp.StatusCode, resp.Status, time.Since(start))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewNetworkError("failed to read response", err)
	}

	if c.debug && c.debugOut != nil {
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Body: %s\n", truncateBody(string(respBody)))
	}

	if err := c.checkError(resp.StatusCode, respBody); err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return errors.NewGeneralError("failed to parse response", err)
		}
	}

	return nil
}

func (c *HTTPClient) checkError(statusCode int, body []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	var apiErr struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	_ = json.Unmarshal(body, &apiErr)

	msg := "unknown error"
	if len(apiErr.Errors) > 0 {
		msg = apiErr.Errors[0].Message
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return errors.NewAuthError(msg)
	case http.StatusForbidden:
		return errors.NewAuthError(msg)
	case http.StatusNotFound:
		return errors.NewNotFoundError("resource")
	case http.StatusTooManyRequests:
		return errors.NewRateLimitedError("")
	default:
		return errors.NewGeneralError(fmt.Sprintf("API error %d: %s", statusCode, msg), nil)
	}
}

func truncateToken(token string) string {
	if len(token) > 8 {
		return token[:8]
	}
	return token
}

func truncateBody(body string) string {
	if len(body) > 500 {
		return body[:500] + "..."
	}
	return body
}

var _ Client = (*HTTPClient)(nil)
