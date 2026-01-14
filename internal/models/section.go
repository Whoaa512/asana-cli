package models

type Section struct {
	GID     string         `json:"gid"`
	Name    string         `json:"name"`
	Project *AsanaResource `json:"project,omitempty"`
}

type SectionCreateRequest struct {
	Name string `json:"name"`
}

type SectionUpdateRequest struct {
	Name string `json:"name"`
}

type SectionInsertRequest struct {
	Section       string  `json:"section"`
	BeforeSection *string `json:"before_section,omitempty"`
	AfterSection  *string `json:"after_section,omitempty"`
}
