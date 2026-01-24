package models

type Workspace struct {
	GID            string   `json:"gid"`
	Name           string   `json:"name"`
	IsOrganization bool     `json:"is_organization"`
	EmailDomains   []string `json:"email_domains,omitempty"`
}

func (w Workspace) GetName() string { return w.Name }
func (w Workspace) GetGID() string  { return w.GID }
