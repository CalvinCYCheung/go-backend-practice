package errorsmodel

import "fmt"

type ErrorModel[T any] struct {
	Err  error
	Data T
}

func (e *ErrorModel[T]) Error() string {
	return fmt.Sprintf("error %v", e.Err)
}

func NewErrorModel[T any](err error, t T) *ErrorModel[T] {
	return &ErrorModel[T]{
		Err:  err,
		Data: t,
	}
}
