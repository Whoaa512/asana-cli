package models

type AsanaResource struct {
	GID  string `json:"gid"`
	Name string `json:"name,omitempty"`
}

type PageInfo struct {
	Offset string `json:"offset,omitempty"`
	URI    string `json:"uri,omitempty"`
}

type ListResponse[T any] struct {
	Data     []T       `json:"data"`
	NextPage *PageInfo `json:"next_page,omitempty"`
}
