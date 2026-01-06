package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListSubtasks(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Task], error) {
	path := fmt.Sprintf("/tasks/%s/subtasks", taskGID)

	sep := "?"
	if limit > 0 {
		path += fmt.Sprintf("%slimit=%d", sep, limit)
		sep = "&"
	}
	if offset != "" {
		path += fmt.Sprintf("%soffset=%s", sep, offset)
	}

	var response struct {
		Data     []models.Task    `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Task]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) AddSubtask(ctx context.Context, parentGID string, name string) (*models.Task, error) {
	payload := struct {
		Data struct {
			Name string `json:"name"`
		} `json:"data"`
	}{Data: struct {
		Name string `json:"name"`
	}{Name: name}}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	var response struct {
		Data models.Task `json:"data"`
	}

	if err := c.post(ctx, fmt.Sprintf("/tasks/%s/subtasks", parentGID), bytes.NewReader(body), &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
