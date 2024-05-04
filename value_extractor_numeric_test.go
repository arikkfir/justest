package justest_test

import (
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestNumericValueExtractor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actual                 any
		expectedOutcome        TestOutcomeExpectation
		expected               any
		expectedOutcomePattern string
	}
	testCases := map[string]testCase{
		"string fails":         {actual: "1", expectedOutcome: ExpectFailure, expectedOutcomePattern: `Unsupported actual value: 1`},
		"chan int":             {actual: ChanOf(1), expectedOutcome: ExpectSuccess, expected: 1},
		"chan string fails":    {actual: ChanOf("1"), expectedOutcome: ExpectFailure, expectedOutcomePattern: `Unsupported actual value: 1`},
		"float32":              {actual: float32(1.1), expectedOutcome: ExpectSuccess, expected: float32(1.1)},
		"float64":              {actual: 1.1, expectedOutcome: ExpectSuccess, expected: 1.1},
		"func int":             {actual: func(t T) any { return 1 }, expectedOutcome: ExpectSuccess, expected: 1},
		"func string fails":    {actual: func(t T) any { return "1" }, expectedOutcome: ExpectFailure, expectedOutcomePattern: `Unsupported actual value: 1`},
		"int":                  {actual: 1, expectedOutcome: ExpectSuccess, expected: 1},
		"int8":                 {actual: int8(1), expectedOutcome: ExpectSuccess, expected: int8(1)},
		"int16":                {actual: int16(1), expectedOutcome: ExpectSuccess, expected: int16(1)},
		"int32":                {actual: int32(1), expectedOutcome: ExpectSuccess, expected: int32(1)},
		"int64":                {actual: int64(1), expectedOutcome: ExpectSuccess, expected: int64(1)},
		"pointer int":          {actual: Ptr(1), expectedOutcome: ExpectSuccess, expected: 1},
		"pointer string fails": {actual: Ptr("1"), expectedOutcome: ExpectFailure, expectedOutcomePattern: `Unsupported actual value: 1`},
		"uint":                 {actual: 1, expectedOutcome: ExpectSuccess, expected: 1},
		"uint8":                {actual: uint8(1), expectedOutcome: ExpectSuccess, expected: uint8(1)},
		"uint16":               {actual: uint16(1), expectedOutcome: ExpectSuccess, expected: uint16(1)},
		"uint32":               {actual: uint32(1), expectedOutcome: ExpectSuccess, expected: uint32(1)},
		"uint64":               {actual: uint64(1), expectedOutcome: ExpectSuccess, expected: uint64(1)},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			mt := NewMockT(t)
			v := NewNumericValueExtractor().MustExtractValue(mt, tc.actual)
			if !cmp.Equal(tc.expected, v) {
				t.Fatalf("Expected '%v', got '%v'", tc.expected, v)
			}
		})
	}
}
