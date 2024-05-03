package justest

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
)

var (
	sayValueExtractor ValueExtractor
)

func init() {
	sayValueExtractor = NewValueExtractor(ExtractorUnsupported)
	sayValueExtractor[reflect.Chan] = NewChannelExtractor(sayValueExtractor, true)
	sayValueExtractor[reflect.Func] = NewFuncExtractor(sayValueExtractor, true)
	sayValueExtractor[reflect.Pointer] = func(t T, v any) (any, bool) {
		GetHelper(t).Helper()
		if bufferPointer, ok := v.(*bytes.Buffer); ok {
			return bufferPointer.String(), true
		} else if stringPointer, ok := v.(*string); ok {
			return *stringPointer, true
		} else if ba, ok := v.(*[]byte); ok {
			return string(*ba), true
		} else {
			t.Fatalf("Unsupported type '%T' for Say matcher: %+v", v, v)
			panic("unreachable")
		}
	}
	sayValueExtractor[reflect.Slice] = func(t T, v any) (any, bool) {
		GetHelper(t).Helper()
		if b, ok := v.([]byte); ok {
			return string(b), true
		} else {
			t.Fatalf("Unsupported type '%T' for Say matcher: %+v", v, v)
			panic("unreachable")
		}
	}
	sayValueExtractor[reflect.String] = ExtractSameValue
}

//go:noinline
func Say[Type string | *regexp.Regexp](expectation Type) Matcher {
	switch e := any(expectation).(type) {
	case string:
		re := regexp.MustCompile(e)
		return MatcherFunc(func(t T, actuals ...any) {
			GetHelper(t).Helper()
			for _, actual := range actuals {
				v := sayValueExtractor.MustExtractValue(t, actual)
				if !re.Match([]byte(v.(string))) {
					t.Fatalf("Expected actual value to match '%s', but it does not: %s", re, v)
				}
			}
		})

	case *regexp.Regexp:
		return MatcherFunc(func(t T, actuals ...any) {
			GetHelper(t).Helper()
			for _, actual := range actuals {
				v := sayValueExtractor.MustExtractValue(t, actual)
				if !e.Match([]byte(v.(string))) {
					t.Fatalf("Expected actual value to match '%s', but it does not: %s", e, v)
				}
			}
		})

	default:
		panic(fmt.Sprintf("unsupported type for Say matcher: %T", expectation))
	}
}
