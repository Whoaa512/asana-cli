package models

type Project struct {
	GID       string         `json:"gid"`
	Name      string         `json:"name"`
	Archived  bool           `json:"archived"`
	Color     string         `json:"color,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	Workspace *AsanaResource `json:"workspace,omitempty"`
}
