package justest_test

import (
	"fmt"
	"regexp"
	"testing"

	. "github.com/arikkfir/justest"
)

func TestFail(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals  []any
		verifier TestOutcomeVerifier
	}
	testCases := map[string]testCase{
		"Fails if no actuals": {
			actuals:  []any{},
			verifier: FailureVerifier(`No error occurred`),
		},
		"Fails if last actual is nil": {
			actuals:  []any{1, 2, nil},
			verifier: FailureVerifier(`No error occurred`),
		},
		"Succeeds if last actual is an error": {
			actuals:  []any{1, fmt.Errorf("expected error")},
			verifier: SuccessVerifier(),
		},
		"Fails if last actual is non-nil and not an error": {
			actuals:  []any{1, 2, 3},
			verifier: FailureVerifier(`No error occurred`),
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			With(mt).VerifyThat(tc.actuals...).Will(Fail()).Now()
		})
	}
	t.Run("Succeeds if error matches one of the patterns", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(SuccessVerifier())
		With(mt).VerifyThat(fmt.Errorf("expected error")).Will(Fail(`^expected error$`)).Now()
	})
	t.Run("Fails if error matches none of the patterns", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(FailureVerifier(`.*` + regexp.QuoteMeta(`[^abc$ ^def$ ^ghi$]`) + `\n.*expected error`))
		With(mt).VerifyThat(fmt.Errorf("expected error")).Will(Fail(`^abc$`, `^def$`, `^ghi$`)).Now()
	})
}
