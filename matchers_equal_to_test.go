package justest_test

import (
	. "github.com/arikkfir/justest"
	"regexp"
	"testing"
)

func TestEqualTo(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals, expected      []any
		comparatorFactory      func(t T, tc *testCase) (Comparator, func())
		expectedOutcome        TestOutcomeExpectation
		expectedFailurePattern string
	}
	testCases := map[string]testCase{
		"Single item equality succeeds": {
			actuals:         []any{1},
			expected:        []any{1},
			expectedOutcome: ExpectSuccess,
		},
		"Multi-item equality succeeds": {
			actuals:         []any{1, 2, 3},
			expected:        []any{1, 2, 3},
			expectedOutcome: ExpectSuccess,
		},
		"Single item difference fails": {
			actuals:                []any{1},
			expected:               []any{2},
			expectedOutcome:        ExpectFailure,
			expectedFailurePattern: regexp.QuoteMeta(`Unexpected difference ("-" lines are expected values; "+" lines are actual values):`) + "\n.*",
		},
		"Multi-item difference fails": {
			actuals:                []any{1, 2, 2},
			expected:               []any{1, 2, 3},
			expectedOutcome:        ExpectFailure,
			expectedFailurePattern: regexp.QuoteMeta(`Unexpected difference ("-" lines are expected values; "+" lines are actual values):`) + "\n.*",
		},
		"Length validated": {
			actuals:                []any{1},
			expected:               []any{1, 2},
			expectedOutcome:        ExpectFailure,
			expectedFailurePattern: `^Unexpected difference: received 1 actual values and 2 expected values.*`,
		},
		"Custom comparator success propagates": {
			actuals:  []any{1},
			expected: []any{1},
			comparatorFactory: func(t T, tc *testCase) (Comparator, func()) {
				comparatorCalled := false
				comparatorFunc := func(t T, expected any, actual any) {
					comparatorCalled = true
					if expected != actual {
						t.Fatalf("Expected & actual value differ: expected %+v, got %+v", expected, actual)
					}
					if expected != tc.expected[0] {
						t.Fatalf("Incorrect 'expected' value provided - should be %+v, got %+v", tc.expected[0], expected)
					}
					if actual != tc.actuals[0] {
						t.Fatalf("Incorrect 'actual' value provided - should be %+v, got %+v", tc.actuals[0], actual)
					}
				}
				verifierFunc := func() {
					if !comparatorCalled {
						t.Fatalf("Comparator was not called")
					}
				}
				return comparatorFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess,
		},
		"Custom comparator invoked once per pair": {
			actuals:  []any{1, 2, 3},
			expected: []any{1, 2, 3},
			comparatorFactory: func(t T, tc *testCase) (Comparator, func()) {
				type invocation struct {
					expected, actual any
				}
				var invocations []invocation
				comparatorFunc := func(t T, expected any, actual any) {
					invocations = append(invocations, invocation{expected, actual})
				}
				verifierFunc := func() {
					if len(invocations) != len(tc.expected) {
						t.Fatalf("Expected %d invocations, got %d", len(tc.expected), len(invocations))
					}
					for i, inv := range invocations {
						if inv.expected != tc.expected[i] {
							t.Fatalf("Incorrect 'expected' value provided - should be %+v, got %+v", tc.expected[i], inv.expected)
						}
						if inv.actual != tc.actuals[i] {
							t.Fatalf("Incorrect 'actual' value provided - should be %+v, got %+v", tc.actuals[i], inv.actual)
						}
					}
				}
				return comparatorFunc, verifierFunc
			},
			expectedOutcome: ExpectSuccess,
		},
		"Custom comparator failure propagates": {
			actuals:  []any{1},
			expected: []any{2},
			comparatorFactory: func(t T, tc *testCase) (Comparator, func()) {
				comparatorCalled := false
				comparatorFunc := func(t T, expected any, actual any) {
					comparatorCalled = true
					if expected != actual {
						t.Fatalf("Expected & actual value differ: expected %+v, got %+v", expected, actual)
					}
					if expected != tc.expected[0] {
						t.Fatalf("Incorrect 'expected' value provided - should be %+v, got %+v", tc.expected[0], expected)
					}
					if actual != tc.actuals[0] {
						t.Fatalf("Incorrect 'actual' value provided - should be %+v, got %+v", tc.actuals[0], actual)
					}
				}
				verifierFunc := func() {
					if !comparatorCalled {
						t.Fatalf("Comparator was not called")
					}
				}
				return comparatorFunc, verifierFunc
			},
			expectedOutcome:        ExpectFailure,
			expectedFailurePattern: `Expected & actual value differ: expected \d+, got \d+.*`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedFailurePattern)
			matcher := EqualTo(tc.expected...)
			var verifierFunc func()
			if tc.comparatorFactory != nil {
				comparatorFunc, v := tc.comparatorFactory(t, &tc)
				if comparatorFunc != nil {
					matcher = matcher.Using(comparatorFunc)
				}
				verifierFunc = v
			}
			With(NewMockT(t)).Verify(tc.actuals...).Will(matcher).OrFail()
			if verifierFunc != nil {
				verifierFunc()
			}
		})
	}
}
