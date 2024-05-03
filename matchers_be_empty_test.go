package justest_test

import (
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"regexp"
	"testing"
)

func TestBeEmpty(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actual               any
		expectedOutcome      TestOutcomeExpectation
		expectFailurePattern string
	}
	//goland:noinspection GoRedundantConversion
	testCases := map[string]testCase{
		"Empty array matches":    {actual: [0]int{}, expectedOutcome: ExpectSuccess},
		"Non-empty array fails":  {actual: [3]int{1, 2, 3}, expectedOutcome: ExpectFailure, expectFailurePattern: regexp.QuoteMeta(`Expected '[1 2 3]' to be empty, but it is not (has a length of 3)`)},
		"Empty chan matches":     {actual: ChanOf[int](), expectedOutcome: ExpectSuccess},
		"Non-empty chan fails":   {actual: ChanOf[int](1, 2, 3), expectedOutcome: ExpectFailure, expectFailurePattern: `Expected '.+' to be empty, but it is not \(has a length of 3\)`},
		"Empty map matches":      {actual: map[int]int{}, expectedOutcome: ExpectSuccess},
		"Non-empty map fails":    {actual: map[int]int{1: 1, 2: 2, 3: 3}, expectedOutcome: ExpectFailure, expectFailurePattern: regexp.QuoteMeta(`Expected 'map[1:1 2:2 3:3]' to be empty, but it is not (has a length of 3)`)},
		"Empty slice matches":    {actual: []int{}, expectedOutcome: ExpectSuccess},
		"Non-empty slice fails":  {actual: []int{1, 2, 3}, expectedOutcome: ExpectFailure, expectFailurePattern: regexp.QuoteMeta(`Expected '[1 2 3]' to be empty, but it is not (has a length of 3)`)},
		"Empty string matches":   {actual: "", expectedOutcome: ExpectSuccess},
		"Non-empty string fails": {actual: "abc", expectedOutcome: ExpectFailure, expectFailurePattern: regexp.QuoteMeta(`Expected 'abc' to be empty, but it is not (has a length of 3)`)},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectFailurePattern)
			With(NewMockT(t)).Verify(tc.actual).Will(BeEmpty()).OrFail()
		})
	}
}
