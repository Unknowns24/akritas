package dto

type DataResponseDTO[T any] struct {
	Data T `json:"data"`
}
