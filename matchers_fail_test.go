package justest_test

import (
	"fmt"
	. "github.com/arikkfir/justest"
	"testing"
)

func TestFail(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals              []any
		expectedOutcome      TestOutcomeExpectation
		expectFailurePattern string
	}
	testCases := map[string]testCase{
		"Fails if no actuals": {
			actuals:              []any{},
			expectedOutcome:      ExpectFailure,
			expectFailurePattern: `No error occurred`,
		},
		"Fails if last actual is nil": {
			actuals:              []any{1, 2, nil},
			expectedOutcome:      ExpectFailure,
			expectFailurePattern: `No error occurred`,
		},
		"Succeeds if last actual is an error": {
			actuals:         []any{1, fmt.Errorf("expected error")},
			expectedOutcome: ExpectSuccess,
		},
		"Fails if last actual is non-nil and not an error": {
			actuals:              []any{1, 2, 3},
			expectedOutcome:      ExpectFailure,
			expectFailurePattern: `No error occurred`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectFailurePattern)
			With(NewMockT(t)).Verify(tc.actuals...).Will(Fail()).OrFail()
		})
	}
}
