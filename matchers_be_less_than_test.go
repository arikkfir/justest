package justest_test

import (
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"reflect"
	"testing"
)

func TestBeLessThan(t *testing.T) {
	type testCase struct {
		expectFailurePattern *string
		actual, max          any
	}
	//goland:noinspection GoRedundantConversion
	testCases := map[reflect.Kind]map[string]testCase{
		reflect.Float32: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5.1 to be less than 5.1`), actual: float32(5.1), max: float32(5.1)},
			"BelowMax succeeds": {actual: float32(5.1), max: float32(9.1)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5.1 to be less than 0.1`), actual: float32(5.1), max: float32(0.1)},
		},
		reflect.Float64: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5.1 to be less than 5.1`), actual: float64(5.1), max: float64(5.1)},
			"BelowMax succeeds": {actual: float64(5.1), max: float64(9.1)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5.1 to be less than 0.1`), actual: float64(5.1), max: float64(0.1)},
		},
		reflect.Int: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: 5, max: 5},
			"BelowMax succeeds": {actual: 5, max: 9},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: 5, max: 0},
		},
		reflect.Int8: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: int8(5), max: int8(5)},
			"BelowMax succeeds": {actual: int8(5), max: int8(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: int8(5), max: int8(0)},
		},
		reflect.Int16: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: int16(5), max: int16(5)},
			"BelowMax succeeds": {actual: int16(5), max: int16(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: int16(5), max: int16(0)},
		},
		reflect.Int32: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: int32(5), max: int32(5)},
			"BelowMax succeeds": {actual: int32(5), max: int32(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: int32(5), max: int32(0)},
		},
		reflect.Int64: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: int64(5), max: int64(5)},
			"BelowMax succeeds": {actual: int64(5), max: int64(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: int64(5), max: int64(0)},
		},
		reflect.Uint: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: uint(5), max: uint(5)},
			"BelowMax succeeds": {actual: uint(5), max: uint(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: uint(5), max: uint(0)},
		},
		reflect.Uint8: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: uint8(5), max: uint8(5)},
			"BelowMax succeeds": {actual: uint8(5), max: uint8(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: uint8(5), max: uint8(0)},
		},
		reflect.Uint16: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: uint16(5), max: uint16(5)},
			"BelowMax succeeds": {actual: uint16(5), max: uint16(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: uint16(5), max: uint16(0)},
		},
		reflect.Uint32: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: uint32(5), max: uint32(5)},
			"BelowMax succeeds": {actual: uint32(5), max: uint32(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: uint32(5), max: uint32(0)},
		},
		reflect.Uint64: {
			"EqualMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 5`), actual: uint64(5), max: uint64(5)},
			"BelowMax succeeds": {actual: uint64(5), max: uint64(9)},
			"AboveMax fails":    {expectFailurePattern: Ptr(`Expected actual value 5 to be less than 0`), actual: uint64(5), max: uint64(0)},
		},
	}
	for kind, kindTestCases := range testCases {
		t.Run(kind.String(), func(t *testing.T) {
			t.Parallel()
			for name, tc := range kindTestCases {
				t.Run(name, func(t *testing.T) {
					t.Parallel()
					if tc.expectFailurePattern != nil {
						defer VerifyTestOutcome(t, ExpectFailure, *tc.expectFailurePattern)
					} else {
						defer VerifyTestOutcome(t, ExpectSuccess, "")
					}
					With(NewMockT(t)).Verify(tc.actual).Will(BeLessThan(tc.max)).OrFail()
				})
			}
		})
	}
	t.Run("MinTypeMismatches", func(t *testing.T) {
		t.Parallel()
		defer VerifyTestOutcome(t, ExpectFailure, `Expected actual value to be of type 'int64', but it is of type 'int'`)
		With(NewMockT(t)).Verify(1).Will(BeLessThan(int64(0))).OrFail()
	})
}
