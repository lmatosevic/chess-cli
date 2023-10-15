package model

type ListResponse[T any] struct {
	Items       []T `json:"items"`
	ResultCount int `json:"resultCount"`
	TotalCount  int `json:"totalCount"`
}
