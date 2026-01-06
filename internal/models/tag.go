package models

type Tag struct {
	GID       string         `json:"gid"`
	Name      string         `json:"name"`
	Color     string         `json:"color,omitempty"`
	Workspace *AsanaResource `json:"workspace,omitempty"`
}

type TagCreateRequest struct {
	Name  string `json:"name"`
	Color string `json:"color,omitempty"`
}
