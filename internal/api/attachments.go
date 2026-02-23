package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListAttachments(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Attachment], error) {
	path := fmt.Sprintf("/tasks/%s/attachments", taskGID)

	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset != "" {
		params.Set("offset", offset)
	}
	if encoded := params.Encode(); encoded != "" {
		path += "?" + encoded
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

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

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
	if err := c.executeWithRetry(ctx, requestOpts{
		method:      "POST",
		path:        path,
		bodyBytes:   buf.Bytes(),
		contentType: writer.FormDataContentType(),
		debugBody:   fmt.Sprintf("<multipart %d bytes>", buf.Len()),
		result:      &response,
	}, 0); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) DeleteAttachment(ctx context.Context, gid string) error {
	return c.delete(ctx, "/attachments/"+gid)
}
