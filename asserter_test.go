package justest_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	. "github.com/arikkfir/justest"
)

func TestWith(t *testing.T) {
	t.Parallel()
	t.Run("panics on nil T", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Expected 'With(nil)' to panic, but it didn't")
			} else if r != "given T instance must not be nil" {
				t.Fatalf("Expected 'With(nil)' to panic with 'given T instance must not be nil', but got: %v", r)
			}
		}()
		With(nil)
	})
}

func TestCorrectActualsPassedToMatcher(t *testing.T) {
	t.Parallel()
	type testCase struct{ actuals []any }
	testCases := map[string]testCase{
		"No actuals":       {actuals: []any{}},
		"Single actual":    {actuals: []any{1}},
		"Multiple actuals": {actuals: []any{1, 2, 3}},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(SuccessVerifier())
			var actualsProvidedToMatcher []any
			With(mt).Verify(tc.actuals...).Will(MatcherFunc(func(t T, actuals ...any) { actualsProvidedToMatcher = actuals })).OrFail()
			if !cmp.Equal(tc.actuals, actualsProvidedToMatcher) {
				t.Fatalf("Incorrect actuals given to Matcher: %s", cmp.Diff(tc.actuals, actualsProvidedToMatcher))
			}
		})
	}
}

func TestMatcherFailureIsPropagated(t *testing.T) {
	t.Parallel()
	mt := NewMockT(t)
	defer mt.Verify(FailureVerifier(`^expected failure(?m:\n^.+:\d+\s+-->\s+.+$){2}$`))
	With(mt).Verify().Will(MatcherFunc(func(t T, a ...any) { t.Fatalf("expected failure") })).OrFail()
}

func TestAssertionFor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals            []any
		duration, interval time.Duration
		matcherFactory     func() MatcherFunc
		verifier           TestOutcomeVerifier
	}
	testCases := map[string]testCase{
		"Constant success is propagated": {
			actuals:        []any{1},
			duration:       1 * time.Second,
			interval:       100 * time.Millisecond,
			matcherFactory: func() MatcherFunc { return func(_ T, actuals ...any) {} },
			verifier:       SuccessVerifier(),
		},
		"At least one success is required": {
			actuals:        []any{1},
			duration:       1 * time.Second,
			interval:       100 * time.Millisecond,
			matcherFactory: func() MatcherFunc { return func(t T, actuals ...any) { time.Sleep(5 * time.Second) } },
			verifier:       FailureVerifier(`^Timed out after \d+s waiting for assertion to pass \(tick never finished once\).*`),
		},
		"Failure is immediately propagated": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherFactory: func() MatcherFunc {
				invocations := 0
				matcherFunc := func(t T, actuals ...any) {
					invocations++
					if invocations == 2 {
						t.Fatalf("expected failure: %d", invocations)
					}
				}
				return matcherFunc
			},
			verifier: FailureVerifier(`^expected failure: 2.*`),
		},
		"No parallel ticks allowed": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherFactory: func() MatcherFunc {
				running := false
				matcherFunc := MatcherFunc(func(t T, actuals ...any) {
					if running {
						t.Fatalf("parallel invocation detected")
					}
					running = true
					defer func() { running = false }()
					time.Sleep(150 * time.Millisecond)
				})
				return matcherFunc
			},
			verifier: SuccessVerifier(),
		},
		"Last failure is specified": {
			actuals:        []any{1},
			duration:       1 * time.Second,
			interval:       100 * time.Millisecond,
			matcherFactory: func() MatcherFunc { return func(t T, actuals ...any) { t.Fatalf("failure") } },
			verifier:       FailureVerifier(`^failure\nAssertion failed after \d+(?:\.\d+)?m?s and did not pass repeatedly for \d+s.*`),
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			With(mt).Verify(tc.actuals...).Will(tc.matcherFactory()).For(tc.duration, tc.interval)
		})
	}
	t.Run("Matcher cleanups are called between intervals", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(SuccessVerifier())

		cleanup1CallTime := time.Time{}
		cleanup2CallTime := time.Time{}

		With(mt).Verify(1).
			Will(MatcherFunc(func(t T, actuals ...any) {
				t.Cleanup(func() { cleanup1CallTime = time.Now(); time.Sleep(1 * time.Second) })
				t.Cleanup(func() { cleanup2CallTime = time.Now(); time.Sleep(1 * time.Second) })
			})).
			For(5*time.Second, 100*time.Millisecond)

		if cleanup1CallTime.IsZero() {
			t.Fatalf("Cleanup 1 was not called")
		}
		if cleanup2CallTime.IsZero() {
			t.Fatalf("Cleanup 2 was not called")
		}
		if cleanup1CallTime.Before(cleanup2CallTime) {
			t.Fatalf("Cleanup 1 (%s) was called after cleanup 2 (%s)", cleanup1CallTime, cleanup2CallTime)
		}
	})
}

func TestAssertionWithin(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals            []any
		duration, interval time.Duration
		matcherFactory     func() MatcherFunc
		verifier           TestOutcomeVerifier
	}
	testCases := map[string]testCase{
		"Success after interim failures but within duration is propagated": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherFactory: func() MatcherFunc {
				invocation := 0
				return func(t T, actuals ...any) {
					invocation++
					time.Sleep(1 * time.Second)
					if invocation < 3 {
						t.Fatalf("interim failure %d", invocation)
					}
				}
			},
			verifier: SuccessVerifier(),
		},
		"Success beyond duration yields timeout failure": {
			actuals:        []any{1},
			duration:       1 * time.Second,
			interval:       100 * time.Millisecond,
			matcherFactory: func() MatcherFunc { return func(t T, actuals ...any) { time.Sleep(5 * time.Second) } },
			verifier:       FailureVerifier(`^Timed out after \d+(\.\d+)?s waiting for assertion to pass \(tick never finished once\)\n.*`),
		},
		"Last failure is specified on timeout failure": {
			actuals:        []any{1},
			duration:       1 * time.Second,
			interval:       100 * time.Millisecond,
			matcherFactory: func() MatcherFunc { return func(t T, actuals ...any) { t.Fatalf("failure") } },
			verifier:       FailureVerifier(`^failure\nTimed out after \d+(\.\d+)?m?s waiting for assertion to pass.*`),
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			matcherFunc := tc.matcherFactory()
			With(mt).Verify(tc.actuals...).Will(matcherFunc).Within(tc.duration, tc.interval)
		})
	}
	t.Run("Success within duration is propagated", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(SuccessVerifier())
		var firstCall time.Time
		matcherFunc := MatcherFunc(func(t T, actuals ...any) { firstCall = time.Now(); time.Sleep(100 * time.Millisecond) })
		With(mt).Verify(1).Will(matcherFunc).Within(5*time.Second, 100*time.Millisecond)
		elapsedDuration := time.Since(firstCall)
		if elapsedDuration > 1*time.Second {
			t.Fatalf("Assertion should have succeeded much faster than 1 second: %s", elapsedDuration)
		}
	})
	t.Run("No parallel ticks allowed", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(SuccessVerifier())
		invocations := 0
		matcherFunc := MatcherFunc(func(t T, actuals ...any) { invocations++; time.Sleep(time.Second) })
		With(mt).Verify(1).Will(matcherFunc).Within(10*time.Second, 100*time.Millisecond)
		if invocations != 1 {
			t.Fatalf("%d invocations occurred, but exactly one was expected", invocations)
		}
	})
	t.Run("Matcher cleanups are called between intervals", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(SuccessVerifier())
		cleanup1CallTime := time.Time{}
		cleanup2CallTime := time.Time{}
		matcherFunc := MatcherFunc(func(t T, actuals ...any) {
			t.Cleanup(func() { cleanup1CallTime = time.Now(); time.Sleep(1 * time.Second) })
			t.Cleanup(func() { cleanup2CallTime = time.Now(); time.Sleep(1 * time.Second) })
		})
		With(mt).Verify(1).Will(matcherFunc).Within(5*time.Second, 100*time.Millisecond)
		if cleanup1CallTime.IsZero() {
			t.Fatalf("Cleanup 1 was not called")
		}
		if cleanup2CallTime.IsZero() {
			t.Fatalf("Cleanup 2 was not called")
		}
		if cleanup1CallTime.Before(cleanup2CallTime) {
			t.Fatalf("Cleanup 1 (%s) was called after cleanup 2 (%s)", cleanup1CallTime, cleanup2CallTime)
		}
	})
}
