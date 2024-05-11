package justest_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
)

func TestNumericValueExtractor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		actual   any
		verifier TestOutcomeVerifier
		expected any
	}
	testCases := map[string]testCase{
		"string fails":         {actual: "1", verifier: FailureVerifier(`Unsupported actual value: 1`)},
		"chan int":             {actual: ChanOf(1), verifier: SuccessVerifier(), expected: 1},
		"chan string fails":    {actual: ChanOf("1"), verifier: FailureVerifier(`Unsupported actual value: 1`)},
		"float32":              {actual: float32(1.1), verifier: SuccessVerifier(), expected: float32(1.1)},
		"float64":              {actual: 1.1, verifier: SuccessVerifier(), expected: 1.1},
		"func int":             {actual: func(t T) any { return 1 }, verifier: SuccessVerifier(), expected: 1},
		"func string fails":    {actual: func(t T) any { return "1" }, verifier: FailureVerifier(`Unsupported actual value: 1`)},
		"int":                  {actual: 1, verifier: SuccessVerifier(), expected: 1},
		"int8":                 {actual: int8(1), verifier: SuccessVerifier(), expected: int8(1)},
		"int16":                {actual: int16(1), verifier: SuccessVerifier(), expected: int16(1)},
		"int32":                {actual: int32(1), verifier: SuccessVerifier(), expected: int32(1)},
		"int64":                {actual: int64(1), verifier: SuccessVerifier(), expected: int64(1)},
		"pointer int":          {actual: Ptr(1), verifier: SuccessVerifier(), expected: 1},
		"pointer string fails": {actual: Ptr("1"), verifier: FailureVerifier(`Unsupported actual value: 1`)},
		"uint":                 {actual: 1, verifier: SuccessVerifier(), expected: 1},
		"uint8":                {actual: uint8(1), verifier: SuccessVerifier(), expected: uint8(1)},
		"uint16":               {actual: uint16(1), verifier: SuccessVerifier(), expected: uint16(1)},
		"uint32":               {actual: uint32(1), verifier: SuccessVerifier(), expected: uint32(1)},
		"uint64":               {actual: uint64(1), verifier: SuccessVerifier(), expected: uint64(1)},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.verifier)
			v := NewNumericValueExtractor().MustExtractValue(mt, tc.actual)
			if !cmp.Equal(tc.expected, v) {
				t.Fatalf("Expected '%v', got '%v'", tc.expected, v)
			}
		})
	}
}
