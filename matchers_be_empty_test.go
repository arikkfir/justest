package justest_test

import (
	"regexp"
	"testing"

	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
)

func TestBeEmpty(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actual   any
		verifier TestOutcomeVerifier
	}
	//goland:noinspection GoRedundantConversion
	testCases := map[string]testCase{
		"Empty array matches":    {actual: [0]int{}, verifier: SuccessVerifier()},
		"Non-empty array fails":  {actual: [3]int{1, 2, 3}, verifier: FailureVerifier(regexp.QuoteMeta(`Expected '[1 2 3]' to be empty, but it is not (has a length of 3)`))},
		"Empty chan matches":     {actual: ChanOf[int](), verifier: SuccessVerifier()},
		"Non-empty chan fails":   {actual: ChanOf[int](1, 2, 3), verifier: FailureVerifier(`Expected '.+' to be empty, but it is not \(has a length of 3\)`)},
		"Empty map matches":      {actual: map[int]int{}, verifier: SuccessVerifier()},
		"Non-empty map fails":    {actual: map[int]int{1: 1, 2: 2, 3: 3}, verifier: FailureVerifier(regexp.QuoteMeta(`Expected 'map[1:1 2:2 3:3]' to be empty, but it is not (has a length of 3)`))},
		"Empty slice matches":    {actual: []int{}, verifier: SuccessVerifier()},
		"Non-empty slice fails":  {actual: []int{1, 2, 3}, verifier: FailureVerifier(regexp.QuoteMeta(`Expected '[1 2 3]' to be empty, but it is not (has a length of 3)`))},
		"Empty string matches":   {actual: "", verifier: SuccessVerifier()},
		"Non-empty string fails": {actual: "abc", verifier: FailureVerifier(regexp.QuoteMeta(`Expected 'abc' to be empty, but it is not (has a length of 3)`))},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			With(mt).VerifyThat(tc.actual).Will(BeEmpty()).Now()
		})
	}
}
