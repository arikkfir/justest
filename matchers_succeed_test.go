package justest_test

import (
	"fmt"
	"testing"

	. "github.com/arikkfir/justest"
)

func TestSucceed(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals  []any
		verifier TestOutcomeVerifier
	}
	testCases := map[string]testCase{
		"Succeeds if no actuals":           {actuals: []any{}, verifier: SuccessVerifier()},
		"Succeeds if last actual is nil":   {actuals: []any{1, 2, nil}, verifier: SuccessVerifier()},
		"Fails if last actual is an error": {actuals: []any{"abc", fmt.Errorf("expected error")}, verifier: FailureVerifier(`Error occurred: expected error`)},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			With(mt).VerifyThat(tc.actuals...).Will(Succeed()).Now()
		})
	}
}
