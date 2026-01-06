package models

type Story struct {
	GID       string         `json:"gid"`
	CreatedAt string         `json:"created_at"`
	Type      string         `json:"type"`
	Text      string         `json:"text,omitempty"`
	CreatedBy *AsanaResource `json:"created_by,omitempty"`
}

type StoryCreateRequest struct {
	Text string `json:"text"`
}
