package justest

import (
	"reflect"
)

//go:noinline
func BeLessThan(max any) Matcher {
	if max == nil {
		panic("expected a non-nil maximum value")
	}

	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()
		for _, actual := range actuals {
			v := NumericValueExtractor.MustExtractValue(t, actual)
			actualValue := reflect.ValueOf(v)

			maximumValue := reflect.ValueOf(max)
			if actualValue.Kind() != maximumValue.Kind() {
				t.Fatalf("Expected actual value to be of type '%T', but it is of type '%T'", max, v)
			}

			resultValues := getNumericCompareFuncFor(t, v).Call([]reflect.Value{actualValue, maximumValue})
			if resultValues[0].Int() >= 0 {
				t.Fatalf("Expected actual value %v to be less than %v", v, max)
			}
		}
	})
}
