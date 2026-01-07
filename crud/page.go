package crud

type Page[T any] struct {
	PageNum   int64 `json:"page_num"`
	PageSize  int64 `json:"page_size"`
	PageCount int64 `json:"page_count"`
	Total     int64 `json:"total"`
	Items     *[]T  `json:"items"`
}
