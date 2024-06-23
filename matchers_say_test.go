package justest_test

import (
	"bytes"
	"regexp"
	"testing"

	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
)

func TestSay(t *testing.T) {
	type testCase struct {
		actuals  []any
		expected string
		verifier TestOutcomeVerifier
	}
	testCases := map[string]testCase{
		"*bytes.Buffer actual":  {actuals: []any{bytes.NewBufferString("abc")}, expected: "^abc$", verifier: SuccessVerifier()},
		"string actual":         {actuals: []any{"abc"}, expected: "^abc$", verifier: SuccessVerifier()},
		"*string actual":        {actuals: []any{Ptr("abc")}, expected: "^abc$", verifier: SuccessVerifier()},
		"*[]byte actual":        {actuals: []any{Ptr([]byte("abc"))}, expected: "^abc$", verifier: SuccessVerifier()},
		"[]byte actual":         {actuals: []any{[]byte("abc")}, expected: "^abc$", verifier: SuccessVerifier()},
		"Non-bytes slice fails": {actuals: []any{[]int{1, 2, 3}}, expected: "^abc$", verifier: FailureVerifier(`Unsupported type '\[\]int' for Say matcher: \[1 2 3\]`)},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			With(mt).VerifyThat(tc.actuals...).Will(Say(tc.expected)).Now()
			With(mt).VerifyThat(tc.actuals...).Will(Say(regexp.MustCompile(tc.expected))).Now()
		})
	}
}
