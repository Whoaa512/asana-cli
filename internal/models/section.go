package models

type Section struct {
	GID     string         `json:"gid"`
	Name    string         `json:"name"`
	Project *AsanaResource `json:"project,omitempty"`
}
