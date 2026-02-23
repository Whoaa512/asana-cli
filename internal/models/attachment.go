package models

type Attachment struct {
	GID          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
	CreatedAt    string `json:"created_at,omitempty"`
	DownloadURL  string `json:"download_url,omitempty"`
	Host         string `json:"host,omitempty"`
	ViewURL      string `json:"view_url,omitempty"`
	Size         int64  `json:"size,omitempty"`
}
