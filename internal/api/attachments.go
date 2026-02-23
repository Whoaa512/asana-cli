package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListAttachments(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Attachment], error) {
	path := fmt.Sprintf("/tasks/%s/attachments", taskGID)

	sep := "?"
	if limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, limit)
		sep = "&"
	}
	if offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, offset)
	}

	var response struct {
		Data     []models.Attachment `json:"data"`
		NextPage *models.PageInfo    `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Attachment]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetAttachment(ctx context.Context, gid string) (*models.Attachment, error) {
	var response struct {
		Data models.Attachment `json:"data"`
	}

	if err := c.get(ctx, "/attachments/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) UploadAttachment(ctx context.Context, taskGID string, filePath string, name string) (*models.Attachment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.NewGeneralError("failed to open file", err)
	}
	defer func() { _ = file.Close() }()

	uploadName := name
	if uploadName == "" {
		uploadName = filepath.Base(filePath)
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	filePart, err := writer.CreateFormFile("file", uploadName)
	if err != nil {
		return nil, errors.NewGeneralError("failed to create multipart file part", err)
	}

	if _, err := io.Copy(filePart, file); err != nil {
		return nil, errors.NewGeneralError("failed to copy file to multipart body", err)
	}

	if err := writer.WriteField("name", uploadName); err != nil {
		return nil, errors.NewGeneralError("failed to write multipart field", err)
	}

	if err := writer.Close(); err != nil {
		return nil, errors.NewGeneralError("failed to finalize multipart body", err)
	}

	var response struct {
		Data models.Attachment `json:"data"`
	}

	path := fmt.Sprintf("/tasks/%s/attachments", taskGID)
	if err := c.doMultipart(ctx, http.MethodPost, path, bytes.NewReader(body.Bytes()), writer.FormDataContentType(), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) DeleteAttachment(ctx context.Context, gid string) error {
	return c.delete(ctx, "/attachments/"+gid)
}

func (c *HTTPClient) doMultipart(ctx context.Context, method, path string, body io.Reader, contentType string, result any) error {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return errors.NewGeneralError("failed to read request body", err)
		}
	}

	return c.doMultipartWithRetry(ctx, method, path, bodyBytes, contentType, result, 0)
}

func (c *HTTPClient) doMultipartWithRetry(ctx context.Context, method, path string, bodyBytes []byte, contentType string, result any, attempt int) error {
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
		req.Header.Set("Content-Type", contentType)
	}

	if c.debug && c.debugOut != nil {
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] %s %s\n", method, url)
		_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Authorization: Bearer %s...\n", truncateToken(c.token))
		if bodyBytes != nil {
			_, _ = fmt.Fprintf(c.debugOut, "[DEBUG] Request Body: <multipart %d bytes>\n", len(bodyBytes))
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

		return c.doMultipartWithRetry(ctx, method, path, bodyBytes, contentType, result, attempt+1)
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
