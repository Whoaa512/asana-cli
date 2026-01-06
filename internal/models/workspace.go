package models

type Workspace struct {
	GID            string   `json:"gid"`
	Name           string   `json:"name"`
	IsOrganization bool     `json:"is_organization"`
	EmailDomains   []string `json:"email_domains,omitempty"`
}
