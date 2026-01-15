package api

import (
	"context"
	"fmt"
	"net/url"

	"github.com/whoaa512/asana-cli/internal/models"
)

type TeamListOptions struct {
	Organization string
	Limit        int
	Offset       string
}

func (c *HTTPClient) ListTeams(ctx context.Context, opts TeamListOptions) (*models.ListResponse[models.Team], error) {
	if opts.Organization == "" {
		return nil, fmt.Errorf("organization is required")
	}

	path := fmt.Sprintf("/organizations/%s/teams", opts.Organization)

	params := url.Values{}
	params.Set("opt_fields", "name")
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset != "" {
		params.Set("offset", opts.Offset)
	}
	path += "?" + params.Encode()

	var response struct {
		Data     []models.Team    `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Team]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

type UserTeamListOptions struct {
	UserGID      string
	Organization string
	Limit        int
	Offset       string
}

func (c *HTTPClient) ListUserTeams(ctx context.Context, opts UserTeamListOptions) (*models.ListResponse[models.Team], error) {
	if opts.Organization == "" {
		return nil, fmt.Errorf("organization is required")
	}

	userGID := opts.UserGID
	if userGID == "" {
		userGID = "me"
	}

	path := fmt.Sprintf("/users/%s/teams", userGID)

	params := url.Values{}
	params.Set("organization", opts.Organization)
	params.Set("opt_fields", "name")
	if opts.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset != "" {
		params.Set("offset", opts.Offset)
	}
	path += "?" + params.Encode()

	var response struct {
		Data     []models.Team    `json:"data"`
		NextPage *models.PageInfo `json:"next_page,omitempty"`
	}

	if err := c.get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &models.ListResponse[models.Team]{
		Data:     response.Data,
		NextPage: response.NextPage,
	}, nil
}

func (c *HTTPClient) GetTeam(ctx context.Context, gid string) (*models.Team, error) {
	var response struct {
		Data models.Team `json:"data"`
	}

	if err := c.get(ctx, "/teams/"+gid, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
