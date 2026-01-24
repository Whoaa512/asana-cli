package models

type Project struct {
	GID       string         `json:"gid"`
	Name      string         `json:"name"`
	Archived  bool           `json:"archived"`
	Color     string         `json:"color,omitempty"`
	Notes     string         `json:"notes,omitempty"`
	Workspace *AsanaResource `json:"workspace,omitempty"`
}

func (p Project) GetName() string { return p.Name }
func (p Project) GetGID() string  { return p.GID }

type ProjectCreateRequest struct {
	Name      string `json:"name"`
	Workspace string `json:"workspace,omitempty"`
	Team      string `json:"team,omitempty"`
	Notes     string `json:"notes,omitempty"`
	Color     string `json:"color,omitempty"`
}
