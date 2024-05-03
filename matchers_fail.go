package justest

import "reflect"

//go:noinline
func Fail() Matcher {
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
			// ok, no-op
			return
		}

		t.Fatalf("No error occurred")
	})
}
