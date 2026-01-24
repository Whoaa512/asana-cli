package models

type Team struct {
	GID          string `json:"gid"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
}

func (t Team) GetName() string { return t.Name }
func (t Team) GetGID() string  { return t.GID }
