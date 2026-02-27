package models

type RegistrySearch[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}

type SearchRequest struct {
	Limit  *int `query:"limit"`
	Offset *int `query:"offset"`
}
