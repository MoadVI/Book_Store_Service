package models

type SearchCriteria struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`

	MinPrice *float64 `json:"min_price"`
	MaxPrice *float64 `json:"max_price"`

	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}
