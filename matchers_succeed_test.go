package justest_test

import (
	"fmt"
	. "github.com/arikkfir/justest"
	"testing"
)

func TestSucceed(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actuals              []any
		expectedOutcome      TestOutcomeExpectation
		expectFailurePattern string
	}
	testCases := map[string]testCase{
		"Succeeds if no actuals": {
			actuals:         []any{},
			expectedOutcome: ExpectSuccess,
		},
		"Succeeds if last actual is nil": {
			actuals:         []any{1, 2, nil},
			expectedOutcome: ExpectSuccess,
		},
		"Fails if last actual is an error": {
			actuals:              []any{"abc", fmt.Errorf("expected error")},
			expectedOutcome:      ExpectFailure,
			expectFailurePattern: `Error occurred: expected error`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectFailurePattern)
			With(NewMockT(t)).Verify(tc.actuals...).Will(Succeed()).OrFail()
		})
	}
}
