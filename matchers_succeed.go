package justest

import (
	"reflect"
)

var (
	succeedValueExtractor = NewValueExtractor(ExtractSameValue)
)

func init() {
	succeedValueExtractor[reflect.Chan] = NewChannelExtractor(succeedValueExtractor, true)
	succeedValueExtractor[reflect.Func] = NewFuncExtractor(succeedValueExtractor, true)
}

//go:noinline
func Succeed() Matcher {
	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()

		resolvedActuals := make([]any, 0, len(actuals))
		for _, actual := range actuals {
			resolvedActual, found := succeedValueExtractor.ExtractValue(t, actual)
			if found {
				resolvedActuals = append(resolvedActuals, resolvedActual)
			}
		}

		l := len(resolvedActuals)
		if l == 0 {
			return
		}

		last := succeedValueExtractor.MustExtractValue(t, resolvedActuals[l-1])
		if last == nil {
			return
		}

		lastRT := reflect.TypeOf(last)
		if lastRT.AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
			t.Fatalf("Error occurred: %+v", last)
		}
	})
}
