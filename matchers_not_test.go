package justest_test

import (
	. "github.com/arikkfir/justest"
	"testing"
)

func TestNot(t *testing.T) {
	type Verifier func(t T)
	type testCase struct {
		actuals              []any
		matcherGenerator     func(t T, tc *testCase) (Matcher, Verifier)
		expectedOutcome      TestOutcomeExpectation
		expectFailurePattern string
	}
	testCases := map[string]testCase{
		"Failed matcher succeeds": {
			actuals: []any{"foo-bar"},
			matcherGenerator: func(t T, tc *testCase) (Matcher, Verifier) {
				m := MatcherFunc(func(t T, actual ...any) { t.Fatalf("should be ignored"); panic("unreachable") })
				return m, nil
			},
			expectedOutcome: ExpectSuccess,
		},
		"Successful matcher fails": {
			actuals: []any{"foo-bar"},
			matcherGenerator: func(t T, tc *testCase) (Matcher, Verifier) {
				return MatcherFunc(func(t T, actual ...any) {}), nil
			},
			expectedOutcome:      ExpectFailure,
			expectFailurePattern: `Expected this matcher to fail, but it did not`,
		},
		"Panicking matcher re-panics": {
			actuals: []any{"foo-bar"},
			matcherGenerator: func(t T, tc *testCase) (Matcher, Verifier) {
				return MatcherFunc(func(t T, actual ...any) { panic("panic propagated") }), nil
			},
			expectedOutcome:      ExpectPanic,
			expectFailurePattern: `panic propagated`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectFailurePattern)
			mt := NewMockT(t)
			m, v := tc.matcherGenerator(mt, &tc)
			With(mt).Verify(tc.actuals...).Will(Not(m)).OrFail()
			if v != nil {
				v(mt)
			}
		})
	}
}
