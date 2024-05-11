package justest_test

import (
	"fmt"
	"regexp"
	"sync"

	. "github.com/arikkfir/justest"
	"github.com/arikkfir/justest/internal"
)

type TestOutcomeVerifier func(t *MockT, recovered any)

type MockT struct {
	Parent      T
	Cleanups    []func()
	LogMessages []internal.FormatAndArgs
	Failures    []internal.FormatAndArgs
	mutex       sync.Mutex
}

//go:noinline
func NewMockT(parent T) *MockT {
	return &MockT{Parent: parent}
}

//go:noinline
func (t *MockT) GetParent() T { return t.Parent }

//go:noinline
func (t *MockT) Name() string { GetHelper(t).Helper(); return t.Parent.Name() }

//go:noinline
func (t *MockT) Cleanup(f func()) {
	GetHelper(t).Helper()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Cleanups = append(t.Cleanups, f)
}

//go:noinline
func (t *MockT) Fatalf(format string, args ...any) {
	GetHelper(t).Helper()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Failures = append(t.Failures, internal.FormatAndArgs{Format: &format, Args: args})
	panic(t)
}

//go:noinline
func (t *MockT) Failed() bool {
	GetHelper(t).Helper()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return len(t.Failures) > 0
}

//go:noinline
func (t *MockT) Log(args ...any) {
	GetHelper(t).Helper()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.LogMessages = append(t.LogMessages, internal.FormatAndArgs{Args: args})
}

//go:noinline
func (t *MockT) Logf(format string, args ...any) {
	GetHelper(t).Helper()
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.LogMessages = append(t.LogMessages, internal.FormatAndArgs{Format: &format, Args: args})
}

func (t *MockT) Verify(verifiers ...TestOutcomeVerifier) {
	GetHelper(t).Helper()
	if root := GetRoot(t); root != nil && root.Failed() {
		root.Log("Root T has already failed, no point verifying anything else")
		return
	}

	r := recover()
	if _, ok := r.(*MockT); ok && r == t {
		// MockT failed, it's not a real panic
		r = nil
	}

	for _, v := range verifiers {
		v(t, r)
	}
}

func FailureVerifier(patterns ...string) TestOutcomeVerifier {
	compiledPatterns := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		compiledPatterns[i] = regexp.MustCompile(pattern)
	}
	return func(t *MockT, recovered any) {
		GetHelper(t).Helper()
		if recovered != nil {
			GetRoot(t).Fatalf("Expected test to fail, but it resulted in a panic: %+v", recovered)
		}

		if len(t.Failures) == 0 {
			GetRoot(t).Fatalf("No failures recorded, but expected at least one failure matching one of: %v", patterns)
		}

		for _, failure := range t.Failures {
			failureMatched := false
			for _, pattern := range compiledPatterns {
				if failure.MatchesRegexp(pattern) {
					failureMatched = true
					break
				}
			}
			if !failureMatched {
				GetRoot(t).Fatalf(""+
					"Unexpected test failure:\n"+
					"--> Expected one of these patterns: %v\n"+
					"--> Actual failure:                 %s", patterns, failure)
			}
		}
	}
}

func PanicVerifier(patterns ...string) TestOutcomeVerifier {
	compiledPatterns := make([]*regexp.Regexp, len(patterns))
	for i, pattern := range patterns {
		compiledPatterns[i] = regexp.MustCompile(pattern)
	}
	return func(t *MockT, recovered any) {
		GetHelper(t).Helper()
		if recovered == nil {
			GetRoot(t).Fatalf("Expected test to result in a panic, but it did not")
		} else {
			panicMsg := fmt.Sprintf("%+v", recovered)
			for _, pattern := range compiledPatterns {
				if pattern.MatchString(panicMsg) {
					return
				}
			}
			GetRoot(t).Fatalf(""+
				"Expected test to result in a different panic:\n"+
				"--> Expected one of these patterns: %v\n"+
				"--> Actual panic:                   %s", patterns, panicMsg)
		}
	}
}

func SuccessVerifier() TestOutcomeVerifier {
	return func(t *MockT, recovered any) {
		GetHelper(t).Helper()
		if recovered != nil {
			GetRoot(t).Fatalf("Expected test to succeed, but it resulted in a panic: %+v", recovered)
		}

		if len(t.Failures) > 0 {
			GetRoot(t).Fatalf("Expected test to succeed, but failures were recorded: %+v", t.Failures)
		}
	}
}
