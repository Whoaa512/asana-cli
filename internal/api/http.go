package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
	maxBackoff     = 30 * time.Second
)

type HTTPClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
	debug      bool
	debugOut   io.Writer
	rng        *rand.Rand
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
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
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

func (c *HTTPClient) post(ctx context.Context, path string, body io.Reader, result any) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

func (c *HTTPClient) put(ctx context.Context, path string, body io.Reader, result any) error {
	return c.do(ctx, http.MethodPut, path, body, result)
}

func (c *HTTPClient) delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *HTTPClient) do(ctx context.Context, method, path string, body io.Reader, result any) error {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return errors.NewGeneralError("failed to read request body", err)
		}
	}

	return c.doWithRetry(ctx, method, path, bodyBytes, result, 0)
}

func (c *HTTPClient) doWithRetry(ctx context.Context, method, path string, bodyBytes []byte, result any, attempt int) error {
	url := c.baseURL + path

	var body io.Reader
	if bodyBytes != nil {
		body = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return errors.NewGeneralError("failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if bodyBytes != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.debug && c.debugOut != nil {
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] %s %s\n", method, url)
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Authorization: Bearer %s...\n", truncateToken(c.token))
		if bodyBytes != nil {
			_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Request Body: %s\n", truncateBody(string(bodyBytes)))
		}
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

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		if attempt >= maxRetries {
			return errors.NewRateLimitedError(formatRetryAfter(retryAfter))
		}

		waitTime := c.calculateBackoff(attempt, retryAfter)
		if c.debug && c.debugOut != nil {
			_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Rate limited, retrying in %s (attempt %d/%d)\n", waitTime, attempt+1, maxRetries)
		}

		select {
		case <-ctx.Done():
			return errors.NewNetworkError("request cancelled", ctx.Err())
		case <-time.After(waitTime):
		}

		return c.doWithRetry(ctx, method, path, bodyBytes, result, attempt+1)
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

	msg := "unknown error"
	if err := json.Unmarshal(body, &apiErr); err != nil {
		if len(body) > 200 {
			msg = fmt.Sprintf("API error (non-JSON): %s...", string(body[:200]))
		} else if len(body) > 0 {
			msg = fmt.Sprintf("API error (non-JSON): %s", string(body))
		}
	} else if len(apiErr.Errors) > 0 {
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

func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(header); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return 0
}

func formatRetryAfter(d time.Duration) string {
	if d == 0 {
		return ""
	}
	return d.String()
}

func (c *HTTPClient) calculateBackoff(attempt int, retryAfter time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}
	backoff := initialBackoff * (1 << attempt)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	jitter := time.Duration(c.rng.Int63n(int64(backoff) / 4))
	return backoff + jitter
}

var _ Client = (*HTTPClient)(nil)
