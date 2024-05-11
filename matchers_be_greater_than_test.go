package justest_test

import (
	"reflect"
	"testing"

	. "github.com/arikkfir/justest"
)

func TestBeGreaterThan(t *testing.T) {
	type testCase struct {
		verifier    TestOutcomeVerifier
		actual, min any
	}
	//goland:noinspection GoRedundantConversion
	testCases := map[reflect.Kind]map[string]testCase{
		reflect.Float32: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5.1 to be greater than 5.1`), actual: float32(5.1), min: float32(5.1)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: float32(5.1), min: float32(0.1)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5.1 to be greater than 6.1`), actual: float32(5.1), min: float32(6.1)},
		},
		reflect.Float64: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5.1 to be greater than 5.1`), actual: float64(5.1), min: float64(5.1)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: float64(5.1), min: float64(0.1)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5.1 to be greater than 6.1`), actual: float64(5.1), min: float64(6.1)},
		},
		reflect.Int: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: 5, min: 5},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: 5, min: 0},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: 5, min: 6},
		},
		reflect.Int8: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: int8(5), min: int8(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: int8(5), min: int8(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: int8(5), min: int8(6)},
		},
		reflect.Int16: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: int16(5), min: int16(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: int16(5), min: int16(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: int16(5), min: int16(6)},
		},
		reflect.Int32: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: int32(5), min: int32(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: int32(5), min: int32(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: int32(5), min: int32(6)},
		},
		reflect.Int64: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: int64(5), min: int64(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: int64(5), min: int64(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: int64(5), min: int64(6)},
		},
		reflect.Uint: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: uint(5), min: uint(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: uint(5), min: uint(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: uint(5), min: uint(6)},
		},
		reflect.Uint8: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: uint8(5), min: uint8(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: uint8(5), min: uint8(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: uint8(5), min: uint8(6)},
		},
		reflect.Uint16: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: uint16(5), min: uint16(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: uint16(5), min: uint16(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: uint16(5), min: uint16(6)},
		},
		reflect.Uint32: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: uint32(5), min: uint32(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: uint32(5), min: uint32(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: uint32(5), min: uint32(6)},
		},
		reflect.Uint64: {
			"EqualMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 5`), actual: uint64(5), min: uint64(5)},
			"AboveMin succeeds": {verifier: SuccessVerifier(), actual: uint64(5), min: uint64(0)},
			"BelowMin fails":    {verifier: FailureVerifier(`Expected actual value 5 to be greater than 6`), actual: uint64(5), min: uint64(6)},
		},
	}
	for kind, kindTestCases := range testCases {
		kind := kind
		kindTestCases := kindTestCases
		t.Run(kind.String(), func(t *testing.T) {
			t.Parallel()
			for name, tc := range kindTestCases {
				tc := tc
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					mt := NewMockT(t)
					defer mt.Verify(tc.verifier)
					With(mt).Verify(tc.actual).Will(BeGreaterThan(tc.min)).OrFail()
				})
			}
		})
	}
	t.Run("MinTypeMismatches", func(t *testing.T) {
		t.Parallel()
		mt := NewMockT(t)
		defer mt.Verify(FailureVerifier(`Expected actual value to be of type 'int64', but it is of type 'int'`))
		With(mt).Verify(1).Will(BeGreaterThan(int64(0))).OrFail()
	})
}
