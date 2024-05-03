package justest

import (
	"fmt"
)

type T interface {
	Name() string
	Cleanup(f func())
	Fatalf(format string, args ...any)
	Failed() bool
	Log(args ...any)
	Logf(format string, args ...any)
}

type HasParent interface{ GetParent() T }

type noOpHelper struct{}

//go:noinline
func (n *noOpHelper) Helper() {}

//go:noinline
func GetHelper(t T) interface{ Helper() } {
	var candidate any = t
	for candidate != nil {
		if h, ok := candidate.(interface{ Helper() }); ok {
			return h
		} else if hp, ok := candidate.(HasParent); ok {
			candidate = hp.GetParent()
		} else {
			panic(fmt.Sprintf("unsupported T instance: %+v (%T)", candidate, candidate))
		}
	}
	return &noOpHelper{}
}
