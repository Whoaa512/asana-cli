package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListDependencies(ctx context.Context, taskGID string) ([]models.Task, error) {
	var response struct {
		Data []models.Task `json:"data"`
	}

	if err := c.get(ctx, fmt.Sprintf("/tasks/%s/dependencies", taskGID), &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *HTTPClient) ListDependents(ctx context.Context, taskGID string) ([]models.Task, error) {
	var response struct {
		Data []models.Task `json:"data"`
	}

	if err := c.get(ctx, fmt.Sprintf("/tasks/%s/dependents", taskGID), &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (c *HTTPClient) AddDependency(ctx context.Context, taskGID string, dependsOnGID string) error {
	payload := struct {
		Data struct {
			Dependencies []string `json:"dependencies"`
		} `json:"data"`
	}{Data: struct {
		Dependencies []string `json:"dependencies"`
	}{Dependencies: []string{dependsOnGID}}}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	return c.post(ctx, fmt.Sprintf("/tasks/%s/addDependencies", taskGID), bytes.NewReader(body), nil)
}

func (c *HTTPClient) RemoveDependency(ctx context.Context, taskGID string, dependsOnGID string) error {
	payload := struct {
		Data struct {
			Dependencies []string `json:"dependencies"`
		} `json:"data"`
	}{Data: struct {
		Dependencies []string `json:"dependencies"`
	}{Dependencies: []string{dependsOnGID}}}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to encode request: %w", err)
	}

	return c.post(ctx, fmt.Sprintf("/tasks/%s/removeDependencies", taskGID), bytes.NewReader(body), nil)
}
