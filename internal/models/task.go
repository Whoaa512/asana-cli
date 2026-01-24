package models

type Task struct {
	GID          string          `json:"gid"`
	Name         string          `json:"name"`
	Notes        string          `json:"notes,omitempty"`
	Completed    bool            `json:"completed"`
	CompletedAt  string          `json:"completed_at,omitempty"`
	DueOn        string          `json:"due_on,omitempty"`
	Assignee     *AsanaResource  `json:"assignee,omitempty"`
	Projects     []AsanaResource `json:"projects,omitempty"`
	Parent       *AsanaResource  `json:"parent,omitempty"`
	Tags         []AsanaResource `json:"tags,omitempty"`
	Dependencies *[]Task         `json:"dependencies,omitempty"`
}

func (t Task) GetName() string { return t.Name }
func (t Task) GetGID() string  { return t.GID }

type TaskCreateRequest struct {
	Name      string   `json:"name"`
	Notes     string   `json:"notes,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	DueOn     string   `json:"due_on,omitempty"`
	Projects  []string `json:"projects,omitempty"`
	Parent    string   `json:"parent,omitempty"`
	Workspace string   `json:"workspace,omitempty"`
}

type TaskUpdateRequest struct {
	Name      *string `json:"name,omitempty"`
	Notes     *string `json:"notes,omitempty"`
	Assignee  *string `json:"assignee,omitempty"`
	DueOn     *string `json:"due_on,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}
