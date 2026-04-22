package api

import (
	"context"
	"fmt"

	"github.com/whoaa512/asana-cli/internal/models"
)

func (c *HTTPClient) ListWorkspaceCustomFields(ctx context.Context, workspaceGID string, limit int) (*models.ListResponse[models.CustomField], error) {
	path := fmt.Sprintf("/workspaces/%s/custom_fields?opt_fields=gid,name,type,enum_options.gid,enum_options.name,enum_options.color,enum_options.enabled", workspaceGID)
	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}

	var response struct {
		Data     []models.CustomField `json:"data"`
		NextPage *models.PageInfo     `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.CustomField]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetCustomField(ctx context.Context, gid string) (*models.CustomField, error) {
	path := "/custom_fields/" + gid + "?opt_fields=gid,name,type,enum_options.gid,enum_options.name,enum_options.color,enum_options.enabled"

	var response struct {
		Data models.CustomField `json:"data"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *HTTPClient) GetProjectCustomFieldSettings(ctx context.Context, projectGID string) (*models.ListResponse[models.CustomFieldSetting], error) {
	path := fmt.Sprintf(
		"/projects/%s/custom_field_settings?opt_fields=gid,is_important,custom_field.gid,custom_field.name,custom_field.type,custom_field.enum_options.gid,custom_field.enum_options.name,custom_field.enum_options.color,custom_field.enum_options.enabled",
		projectGID,
	)

	var response struct {
		Data     []models.CustomFieldSetting `json:"data"`
		NextPage *models.PageInfo            `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.CustomFieldSetting]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}
