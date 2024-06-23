package justest_test

import (
	"testing"

	. "github.com/arikkfir/justest"
)

func TestNot(t *testing.T) {
	type testCase struct {
		actuals  []any
		matcher  MatcherFunc
		verifier TestOutcomeVerifier
	}
	testCases := map[string]testCase{
		"Failed matcher succeeds": {
			actuals:  []any{"foo-bar"},
			matcher:  MatcherFunc(func(t T, actual ...any) { t.Fatalf("should be ignored"); panic("unreachable") }),
			verifier: SuccessVerifier(),
		},
		"Successful matcher fails": {
			actuals:  []any{"foo-bar"},
			matcher:  MatcherFunc(func(t T, actual ...any) {}),
			verifier: FailureVerifier(`Expected mismatch did not happen`),
		},
		"Panicking matcher re-panics": {
			actuals:  []any{"foo-bar"},
			matcher:  MatcherFunc(func(t T, actual ...any) { panic("panic propagated") }),
			verifier: PanicVerifier(`panic propagated`),
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			With(mt).VerifyThat(tc.actuals...).Will(Not(tc.matcher)).Now()
		})
	}
}
