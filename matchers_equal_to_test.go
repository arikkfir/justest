package justest_test

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	. "github.com/arikkfir/justest"
)

func TestEqualTo(t *testing.T) {
	t.Parallel()
	t.Run("Builtin comparators", func(t *testing.T) {
		type testCase struct {
			actuals, expected []any
			verifier          TestOutcomeVerifier
		}
		testCases := map[string]testCase{
			"Single item equality succeeds": {
				actuals:  []any{1},
				expected: []any{1},
				verifier: SuccessVerifier(),
			},
			"Multi-item equality succeeds": {
				actuals:  []any{1, 2, 3},
				expected: []any{1, 2, 3},
				verifier: SuccessVerifier(),
			},
			"Single item difference fails": {
				actuals:  []any{1},
				expected: []any{2},
				verifier: FailureVerifier(regexp.QuoteMeta(`Unexpected difference ("-" lines are expected values; "+" lines are actual values):`) + "\n.*"),
			},
			"Multi-item difference fails": {
				actuals:  []any{1, 2, 2},
				expected: []any{1, 2, 3},
				verifier: FailureVerifier(regexp.QuoteMeta(`Unexpected difference ("-" lines are expected values; "+" lines are actual values):`) + "\n.*"),
			},
			"Length validated": {
				actuals:  []any{1},
				expected: []any{1, 2},
				verifier: FailureVerifier(`^Unexpected difference: received 1 actual values and 2 expected values.*`),
			},
			"Failure diff uses opts": {
				actuals: []any{
					struct {
						Public  string
						private string
					}{
						Public:  "public-value",
						private: "private-value",
					},
				},
				expected: []any{
					struct {
						Public  string
						private string
					}{Public: "incorrect-value", private: "private-value"},
					cmpopts.IgnoreUnexported(struct {
						Public  string
						private string
					}{})},
				verifier: FailureVerifier(regexp.QuoteMeta(`Unexpected difference ("-" lines are expected values; "+" lines are actual values):`) + "\n.*"),
			},
		}
		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				mt := NewMockT(t)
				defer mt.Verify(tc.verifier)
				With(mt).Verify(tc.actuals...).Will(EqualTo(tc.expected...)).OrFail()
			})
		}
	})
	t.Run("Custom comparators", func(t *testing.T) {
		type testCase struct {
			actuals, expected []any
			comparatorFactory func(t T, tc *testCase) (Comparator, func())
			outcomeVerifier   TestOutcomeVerifier
		}
		testCases := map[string]testCase{
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
				outcomeVerifier: SuccessVerifier(),
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
				outcomeVerifier: SuccessVerifier(),
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
				outcomeVerifier: FailureVerifier(`Expected & actual value differ: expected \d+, got \d+.*`),
			},
		}
		for name, tc := range testCases {
			tc := tc
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				mt := NewMockT(t)
				defer mt.Verify(tc.outcomeVerifier)
				comparatorFunc, verifierFunc := tc.comparatorFactory(mt, &tc)
				With(mt).Verify(tc.actuals...).Will(EqualTo(tc.expected...).Using(comparatorFunc)).OrFail()
				if verifierFunc != nil {
					verifierFunc()
				}
			})
		}
	})
}
