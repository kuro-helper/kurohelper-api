package dto

type TResponse[T any] struct {
	Message string
	Data    T
}
