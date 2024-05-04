package justest_test

import (
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestWith(t *testing.T) {
	t.Parallel()
	t.Run("panics on nil T", func(t *testing.T) {
		t.Parallel()
		defer VerifyTestOutcome(t, ExpectPanic, `^given T instance must not be nil$`)
		With(nil)
	})
}

func TestAssertionOrFail(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals                   []any
		matcherAndVerifierFactory func() (MatcherFunc, func(*testing.T, *testCase))
		expectedOutcome           TestOutcomeExpectation
		expectedOutcomePattern    string
		expectedLogMessages       []FormatAndArgs
	}
	testCases := map[string]testCase{
		"Matcher receives no actuals when none are supplied": {
			actuals: []any{},
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				var actualsProvidedToMatcher []any
				matcherFunc := MatcherFunc(func(t T, actuals ...any) { actualsProvidedToMatcher = actuals })
				verifierFunc := func(t *testing.T, tc *testCase) {
					if !cmp.Equal(tc.actuals, actualsProvidedToMatcher) {
						t.Fatalf("Incorrect actuals given to Matcher: %s", cmp.Diff(tc.actuals, actualsProvidedToMatcher))
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess, // verifier can still fail the test
		},
		"Matcher receives correct single actuals": {
			actuals: []any{1},
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				var actualsProvidedToMatcher []any
				matcherFunc := MatcherFunc(func(t T, actuals ...any) { actualsProvidedToMatcher = actuals })
				verifierFunc := func(t *testing.T, tc *testCase) {
					if !cmp.Equal(tc.actuals, actualsProvidedToMatcher) {
						t.Fatalf("Incorrect actuals given to Matcher: %s", cmp.Diff(tc.actuals, actualsProvidedToMatcher))
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess, // verifier can still fail the test
		},
		"Matcher receives correct multiple actuals": {
			actuals: []any{1, 2, 3},
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				var actualsProvidedToMatcher []any
				matcherFunc := MatcherFunc(func(t T, actuals ...any) {
					actualsProvidedToMatcher = actuals
				})
				verifierFunc := func(t *testing.T, tc *testCase) {
					if !cmp.Equal(tc.actuals, actualsProvidedToMatcher) {
						t.Fatalf("Incorrect actuals given to Matcher: %s", cmp.Diff(tc.actuals, actualsProvidedToMatcher))
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess, // verifier can still fail the test
		},
		"Matcher success is propagated": {
			actuals: []any{1},
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				return func(t T, a ...any) {}, nil
			},
			expectedOutcome: ExpectSuccess,
		},
		"Matcher failure is propagated": {
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				return func(t T, a ...any) { t.Fatalf("expected failure") }, nil
			},
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^expected failure.*`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			matcherFunc, verifierFunc := tc.matcherAndVerifierFactory()
			With(NewMockT(t)).Verify(tc.actuals...).Will(matcherFunc).OrFail()
			if verifierFunc != nil {
				verifierFunc(t, &tc)
			}
		})
	}
}

func TestAssertionFor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals                   []any
		duration, interval        time.Duration
		matcherAndVerifierFactory func() (MatcherFunc, func(*testing.T, *testCase))
		expectedOutcome           TestOutcomeExpectation
		expectedOutcomePattern    string
		expectedLogMessages       []FormatAndArgs
	}
	testCases := map[string]testCase{
		"Constant success is propagated": {
			actuals:                   []any{1},
			duration:                  1 * time.Second,
			interval:                  100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) { return func(t T, actuals ...any) {}, nil },
			expectedOutcome:           ExpectSuccess,
		},
		"At least one success is required": {
			actuals:  []any{1},
			duration: 1 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				return func(t T, actuals ...any) { time.Sleep(5 * time.Second) }, nil
			},
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^Timed out after \d+s waiting for assertion to pass \(tick never finished once\).*`,
		},
		"Failure is immediately propagated": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				invocations := 0
				matcherFunc := func(t T, actuals ...any) {
					invocations++
					if invocations == 2 {
						t.Fatalf("expected failure: %d", invocations)
					}
				}
				return matcherFunc, nil
			},
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^expected failure: 2.*`,
		},
		"No parallel ticks allowed": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				running := false
				matcherFunc := MatcherFunc(func(t T, actuals ...any) {
					if running {
						t.Fatalf("parallel invocation detected")
					}
					running = true
					defer func() { running = false }()
					time.Sleep(150 * time.Millisecond)
				})
				return matcherFunc, nil
			},
			expectedOutcome: ExpectSuccess,
		},
		"Last failure is specified": {
			actuals:  []any{1},
			duration: 1 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				return func(t T, actuals ...any) { t.Fatalf("failure") }, nil
			},
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^failure\nAssertion failed after \d+(?:\.\d+)?m?s and did not pass repeatedly for \d+s.*`,
		},
		"Matcher cleanups are called between intervals": {
			actuals:  []any{1},
			duration: 5 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				cleanup1CallTime := time.Time{}
				cleanup2CallTime := time.Time{}
				matcherFunc := MatcherFunc(func(t T, actuals ...any) {
					t.Cleanup(func() { cleanup1CallTime = time.Now(); time.Sleep(1 * time.Second) })
					t.Cleanup(func() { cleanup2CallTime = time.Now(); time.Sleep(1 * time.Second) })
				})
				verifierFunc := func(t *testing.T, tc *testCase) {
					if cleanup1CallTime.IsZero() {
						t.Fatalf("Cleanup 1 was not called")
					}
					if cleanup2CallTime.IsZero() {
						t.Fatalf("Cleanup 2 was not called")
					}
					if cleanup1CallTime.Before(cleanup2CallTime) {
						t.Fatalf("Cleanup 1 (%s) was called after cleanup 2 (%s)", cleanup1CallTime, cleanup2CallTime)
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			matcherFunc, verifierFunc := tc.matcherAndVerifierFactory()
			With(NewMockT(t)).Verify(tc.actuals...).Will(matcherFunc).For(tc.duration, tc.interval)
			if verifierFunc != nil {
				verifierFunc(t, &tc)
			}
		})
	}
}

func TestAssertionWithin(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals                   []any
		duration, interval        time.Duration
		matcherAndVerifierFactory func() (MatcherFunc, func(*testing.T, *testCase))
		expectedOutcome           TestOutcomeExpectation
		expectedOutcomePattern    string
		expectedLogMessages       []FormatAndArgs
	}
	testCases := map[string]testCase{
		"Success within duration is propagated": {
			actuals:  []any{1},
			duration: 5 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				var firstCall time.Time
				matcherFunc := MatcherFunc(func(t T, actuals ...any) { firstCall = time.Now(); time.Sleep(100 * time.Millisecond) })
				verifierFunc := func(t *testing.T, tc *testCase) {
					elapsedDuration := time.Since(firstCall)
					if elapsedDuration > 1*time.Second {
						t.Fatalf("Assertion should have succeeded much faster than 1 second: %s", elapsedDuration)
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess, // verifier can still fail the test
		},
		"Success after interim failures but within duration is propagated": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				invocation := 0
				return func(t T, actuals ...any) {
					invocation++
					time.Sleep(1 * time.Second)
					if invocation < 3 {
						t.Fatalf("interim failure %d", invocation)
					}
				}, nil
			},
			expectedOutcome: ExpectSuccess,
		},
		"No parallel ticks allowed": {
			actuals:  []any{1},
			duration: 10 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				invocations := 0
				matcherFunc := MatcherFunc(func(t T, actuals ...any) { invocations++; time.Sleep(time.Second) })
				verifierFunc := func(t *testing.T, tc *testCase) {
					if invocations != 1 {
						t.Fatalf("%d invocations occurred, but exactly one was expected", invocations)
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess, // verifier can still fail the test
		},
		"Success beyond duration yields timeout failure": {
			actuals:  []any{1},
			duration: 1 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				return func(t T, actuals ...any) { time.Sleep(5 * time.Second) }, nil
			},
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^Timed out after \d+(\.\d+)?s waiting for assertion to pass \(tick never finished once\)\n.*`,
		},
		"Last failure is specified on timeout failure": {
			actuals:  []any{1},
			duration: 1 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				return func(t T, actuals ...any) { t.Fatalf("failure") }, nil
			},
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^failure\nTimed out after \d+(\.\d+)?m?s waiting for assertion to pass.*`,
		},
		"Matcher cleanups are called between intervals": {
			actuals:  []any{1},
			duration: 5 * time.Second,
			interval: 100 * time.Millisecond,
			matcherAndVerifierFactory: func() (MatcherFunc, func(*testing.T, *testCase)) {
				cleanup1CallTime := time.Time{}
				cleanup2CallTime := time.Time{}
				matcherFunc := MatcherFunc(func(t T, actuals ...any) {
					t.Cleanup(func() { cleanup1CallTime = time.Now(); time.Sleep(1 * time.Second) })
					t.Cleanup(func() { cleanup2CallTime = time.Now(); time.Sleep(1 * time.Second) })
				})
				verifierFunc := func(t *testing.T, tc *testCase) {
					if cleanup1CallTime.IsZero() {
						t.Fatalf("Cleanup 1 was not called")
					}
					if cleanup2CallTime.IsZero() {
						t.Fatalf("Cleanup 2 was not called")
					}
					if cleanup1CallTime.Before(cleanup2CallTime) {
						t.Fatalf("Cleanup 1 (%s) was called after cleanup 2 (%s)", cleanup1CallTime, cleanup2CallTime)
					}
				}
				return matcherFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			matcherFunc, verifierFunc := tc.matcherAndVerifierFactory()
			With(NewMockT(t)).Verify(tc.actuals...).Will(matcherFunc).Within(tc.duration, tc.interval)
			if verifierFunc != nil {
				verifierFunc(t, &tc)
			}
		})
	}
}
