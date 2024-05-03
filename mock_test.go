package justest_test

import (
	"fmt"
	. "github.com/arikkfir/justest"
	"github.com/arikkfir/justest/internal"
	"regexp"
)

type MockT struct {
	Parent      T
	Cleanups    []func()
	Failures    []internal.FormatAndArgs
	LogMessages []internal.FormatAndArgs
}

//go:noinline
func NewMockT(parent T) *MockT {
	return &MockT{Parent: parent}
}

//go:noinline
func (t *MockT) GetParent() T { return t.Parent }

//go:noinline
func (t *MockT) Name() string {
	return t.Parent.Name()
}

//go:noinline
func (t *MockT) Cleanup(f func()) { GetHelper(t).Helper(); t.Cleanups = append(t.Cleanups, f) }

//go:noinline
func (t *MockT) Fatalf(format string, args ...any) {
	GetHelper(t).Helper()
	t.Failures = append(t.Failures, internal.FormatAndArgs{Format: &format, Args: args})
	panic(t)
}

//go:noinline
func (t *MockT) Failed() bool { GetHelper(t).Helper(); return len(t.Failures) > 0 }

//go:noinline
func (t *MockT) Log(args ...any) {
	GetHelper(t).Helper()
	t.LogMessages = append(t.LogMessages, internal.FormatAndArgs{Args: args})
}

//go:noinline
func (t *MockT) Logf(format string, args ...any) {
	GetHelper(t).Helper()
	t.LogMessages = append(t.LogMessages, internal.FormatAndArgs{Format: &format, Args: args})
}

type TestOutcomeExpectation string

const (
	ExpectFailure TestOutcomeExpectation = "expect failure"
	ExpectPanic   TestOutcomeExpectation = "expect panic"
	ExpectSuccess TestOutcomeExpectation = "expect success"
)

//go:noinline
func VerifyTestOutcome(t T, expectedOutcome TestOutcomeExpectation, pattern string) {
	GetHelper(t).Helper()
	if t.Failed() {
		// If the given T has already failed, there's no point verifying the mock T (
		// which would be a potential panic recovered)
		return
	}

	switch expectedOutcome {
	case ExpectFailure:
		if r := recover(); r == nil {
			t.Fatalf("Expected test failure did not happen")
		} else if mt, ok := r.(*MockT); !ok {
			t.Fatalf("Unexpected panic '%+v' happened instead of an expected test failure", mt)
		} else {
			for _, f := range mt.Failures {
				if actualMsg := f.String(); !regexp.MustCompile(pattern).MatchString(actualMsg) {
					t.Fatalf(""+
						"Unexpected test failure:\n"+
						"--> Expected pattern: %s\n"+
						"--> Message:          %s", pattern, actualMsg)
				}
			}
		}
	case ExpectPanic:
		if r := recover(); r == nil {
			t.Fatalf("Expected panic did not happen")
		} else if msg := fmt.Sprintf("%+v", r); !regexp.MustCompile(pattern).MatchString(msg) {
			t.Fatalf("Expected panic matching '%s', but got: %s", pattern, msg)
		}
	case ExpectSuccess:
		if r := recover(); r == nil {
			// ok; no-op
		} else if mt, ok := r.(*MockT); !ok {
			t.Fatalf("Unexpected panic '%+v' happened", mt)
		} else {
			msg := ""
			for _, f := range mt.Failures {
				msg = fmt.Sprintf("%s\n%s", msg, f.String())
			}
			t.Fatalf("Test failure(s) happened when no test failures were expected:%s", msg)
		}
	}
}
