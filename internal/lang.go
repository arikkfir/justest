package internal

import (
	"reflect"
)

//go:noinline
func IsErrorType(ot reflect.Type) bool {
	return ot.Kind() == reflect.Interface && ot.Implements(reflect.TypeOf((*error)(nil)).Elem())
}

//go:noinline
func ChanOf[T any](values ...T) chan T {
	c := make(chan T, len(values))
	for _, v := range values {
		c <- v
	}
	return c
}

//go:noinline
func Ptr[T any](x T) *T {
	return &[]T{x}[0]
}
