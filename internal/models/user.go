package models

type User struct {
	GID          string          `json:"gid"`
	Name         string          `json:"name"`
	Email        string          `json:"email,omitempty"`
	Photo        *UserPhoto      `json:"photo,omitempty"`
	ResourceType string          `json:"resource_type,omitempty"`
	Workspaces   []AsanaResource `json:"workspaces,omitempty"`
}

type UserPhoto struct {
	Image21x21   string `json:"image_21x21,omitempty"`
	Image27x27   string `json:"image_27x27,omitempty"`
	Image36x36   string `json:"image_36x36,omitempty"`
	Image60x60   string `json:"image_60x60,omitempty"`
	Image128x128 string `json:"image_128x128,omitempty"`
}
