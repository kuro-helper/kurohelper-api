package dto

type TResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
