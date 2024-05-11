package justest

import (
	"reflect"
	"regexp"
)

//go:noinline
func Fail(patterns ...string) Matcher {
	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()

		l := len(actuals)
		if l == 0 {
			t.Fatalf("No error occurred")
		}

		last := actuals[l-1]
		if last == nil {
			t.Fatalf("No error occurred")
		}

		lastRT := reflect.TypeOf(last)
		if lastRT.AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
			if len(patterns) > 0 {
				for _, pattern := range patterns {
					re := regexp.MustCompile(pattern)
					if re.MatchString(last.(error).Error()) {
						return
					}
				}
				t.Fatalf("Error message did not match any of these patterns: %v", patterns)
			} else {
				return
			}
		}

		t.Fatalf("No error occurred")
	})
}
