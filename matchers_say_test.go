package justest_test

import (
	"bytes"
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"regexp"
	"testing"
)

func TestSay(t *testing.T) {
	type testCase struct {
		actuals              []any
		expected             string
		expectedOutcome      TestOutcomeExpectation
		expectFailurePattern string
	}
	testCases := map[string]testCase{
		"*bytes.Buffer actual": {
			actuals:         []any{bytes.NewBuffer([]byte("abc"))},
			expected:        "^abc$",
			expectedOutcome: ExpectSuccess,
		},
		"string actual": {
			actuals:         []any{"abc"},
			expected:        "^abc$",
			expectedOutcome: ExpectSuccess,
		},
		"*string actual": {
			actuals:         []any{Ptr("abc")},
			expected:        "^abc$",
			expectedOutcome: ExpectSuccess,
		},
		"*[]byte actual": {
			actuals:         []any{Ptr([]byte("abc"))},
			expected:        "^abc$",
			expectedOutcome: ExpectSuccess,
		},
		"[]byte actual": {
			actuals:         []any{[]byte("abc")},
			expected:        "^abc$",
			expectedOutcome: ExpectSuccess,
		},
		"Non-bytes slice fails": {
			actuals:              []any{[]int{1, 2, 3}},
			expected:             "^abc$",
			expectedOutcome:      ExpectFailure,
			expectFailurePattern: `Unsupported type '\[\]int' for Say matcher: \[1 2 3\]`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectFailurePattern)
			With(NewMockT(t)).Verify(tc.actuals...).Will(Say(tc.expected)).OrFail()
			With(NewMockT(t)).Verify(tc.actuals...).Will(Say(regexp.MustCompile(tc.expected))).OrFail()
		})
	}
}
