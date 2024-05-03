package justest

import (
	"reflect"
)

//go:noinline
func BeBetween(min, max any) Matcher {
	if min == nil {
		panic("expected a non-nil minimum value")
	}

	if max == nil {
		panic("expected a non-nil maximum value")
	}

	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()
		for _, actual := range actuals {
			v := NumericValueExtractor.MustExtractValue(t, actual)
			actualValue := reflect.ValueOf(v)

			minimumValue := reflect.ValueOf(min)
			if actualValue.Type().Kind() != minimumValue.Type().Kind() {
				t.Fatalf("Expected actual value to be of type '%T', but it is of type '%T'", min, v)
			}

			maximumValue := reflect.ValueOf(max)
			if actualValue.Type().Kind() != maximumValue.Type().Kind() {
				t.Fatalf("Expected actual value to be of type '%T', but it is of type '%T'", max, v)
			}

			cmpCompareFunctionValue := getNumericCompareFuncFor(t, v)

			resultActualVsMin := cmpCompareFunctionValue.Call([]reflect.Value{actualValue, minimumValue})
			if result := resultActualVsMin[0].Int(); result < 0 {
				t.Fatalf("Expected actual value %v to be between %v and %v", v, min, max)
			}

			resultActualVsMax := cmpCompareFunctionValue.Call([]reflect.Value{actualValue, maximumValue})
			if result := resultActualVsMax[0].Int(); result > 0 {
				t.Fatalf("Expected actual value %v to be between %v and %v", v, min, max)
			}
		}
	})
}
