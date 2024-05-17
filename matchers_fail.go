package justest

import (
	"reflect"
	"regexp"
)

//go:noinline
func Fail(patterns ...string) Matcher {
	const failureFormatMsg = `Error message did not match any pattern:
	Patterns: %v
	Error:    %s`
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
				msg := last.(error).Error()
				for _, pattern := range patterns {
					re := regexp.MustCompile(pattern)
					if re.MatchString(msg) {
						return
					}
				}
				t.Fatalf(failureFormatMsg, patterns, msg)
			} else {
				return
			}
		}

		t.Fatalf("No error occurred")
	})
}
