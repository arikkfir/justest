package justest

import (
	"reflect"
)

var (
	nilValueExtractor = NewValueExtractor(ExtractSameValue)
)

//go:noinline
func BeNil() Matcher {
	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()
		for _, actual := range actuals {
			v := nilValueExtractor.MustExtractValue(t, actual)
			switch rv := reflect.ValueOf(v); rv.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
				if !rv.IsNil() {
					t.Fatalf("Expected actual to be nil, but it is not: %+v", v)
				}
			default:
				if v != nil {
					t.Fatalf("Expected actual to be nil, but it is not: %+v", v)
				}
			}
		}
	})
}
