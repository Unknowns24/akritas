package dto

type ListResponseDTO[T any] struct {
	Data   []T       `json:"data"`
	Paging PagingDTO `json:"paging"`
}
