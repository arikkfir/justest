package justest

import (
	"reflect"
)

//go:noinline
func BeGreaterThan(min any) Matcher {
	if min == nil {
		panic("expected a non-nil minimum value")
	}

	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()
		for _, actual := range actuals {
			v := NumericValueExtractor.MustExtractValue(t, actual)
			actualValue := reflect.ValueOf(v)

			minimumValue := reflect.ValueOf(min)
			if actualValue.Kind() != minimumValue.Kind() {
				t.Fatalf("Expected actual value to be of type '%T', but it is of type '%T'", min, v)
			}

			resultValues := getNumericCompareFuncFor(t, v).Call([]reflect.Value{actualValue, minimumValue})
			if resultValues[0].Int() <= 0 {
				t.Fatalf("Expected actual value %v to be greater than %v", v, min)
			}
		}
	})
}
