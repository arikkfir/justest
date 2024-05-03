package justest

import (
	"reflect"
)

var (
	emptyValueExtractor ValueExtractor
	lengthExtractor     Extractor = func(t T, v any) (any, bool) {
		GetHelper(t).Helper()
		return reflect.ValueOf(v).Len(), true
	}
)

func init() {
	emptyValueExtractor = NewValueExtractor(ExtractorUnsupported)
	emptyValueExtractor[reflect.Array] = lengthExtractor
	emptyValueExtractor[reflect.Chan] = lengthExtractor
	emptyValueExtractor[reflect.Map] = lengthExtractor
	emptyValueExtractor[reflect.Pointer] = NewPointerExtractor(emptyValueExtractor, true)
	emptyValueExtractor[reflect.Slice] = lengthExtractor
	emptyValueExtractor[reflect.String] = lengthExtractor
}

//go:noinline
func BeEmpty() Matcher {
	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()
		for _, actual := range actuals {
			length := emptyValueExtractor.MustExtractValue(t, actual).(int)
			if length != 0 {
				t.Fatalf("Expected '%+v' to be empty, but it is not (has a length of %d)", actual, length)
			}
		}
	})
}
