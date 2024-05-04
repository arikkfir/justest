package justest_test

import (
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"reflect"
	"testing"
)

func TestBeBetween(t *testing.T) {
	type testCase struct {
		expectFailurePattern *string
		actual, min, max     any
	}
	//goland:noinspection GoRedundantConversion
	testCases := map[reflect.Kind]map[string]testCase{
		reflect.Float32: {
			"EqualMin succeeds":    {actual: float32(5.1), min: float32(5.1), max: float32(9.1)},
			"WithinRange succeeds": {actual: float32(5.1), min: float32(0.1), max: float32(9.1)},
			"EqualMax succeeds":    {actual: float32(5.1), min: float32(0.1), max: float32(5.1)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5.1 to be between 6.1 and 9.1`), actual: float32(5.1), min: float32(6.1), max: float32(9.1)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10.1 to be between 0.1 and 9.1`), actual: float32(10.1), min: float32(0.1), max: float32(9.1)},
		},
		reflect.Float64: {
			"EqualMin succeeds":    {actual: float64(5.1), min: float64(5.1), max: float64(9.1)},
			"WithinRange succeeds": {actual: float64(5.1), min: float64(0.1), max: float64(9.1)},
			"EqualMax succeeds":    {actual: float64(5.1), min: float64(0.1), max: float64(5.1)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5.1 to be between 6.1 and 9.1`), actual: float64(5.1), min: float64(6.1), max: float64(9.1)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10.1 to be between 0.1 and 9.1`), actual: float64(10.1), min: float64(0.1), max: float64(9.1)},
		},
		reflect.Int: {
			"EqualMin succeeds":    {actual: 5, min: 5, max: 9},
			"WithinRange succeeds": {actual: 5, min: 0, max: 9},
			"EqualMax succeeds":    {actual: 5, min: 0, max: 5},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: 5, min: 6, max: 9},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: 10, min: 0, max: 9},
		},
		reflect.Int8: {
			"EqualMin succeeds":    {actual: int8(5), min: int8(5), max: int8(9)},
			"WithinRange succeeds": {actual: int8(5), min: int8(0), max: int8(9)},
			"EqualMax succeeds":    {actual: int8(5), min: int8(0), max: int8(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: int8(5), min: int8(6), max: int8(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: int8(10), min: int8(0), max: int8(9)},
		},
		reflect.Int16: {
			"EqualMin succeeds":    {actual: int16(5), min: int16(5), max: int16(9)},
			"WithinRange succeeds": {actual: int16(5), min: int16(0), max: int16(9)},
			"EqualMax succeeds":    {actual: int16(5), min: int16(0), max: int16(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: int16(5), min: int16(6), max: int16(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: int16(10), min: int16(0), max: int16(9)},
		},
		reflect.Int32: {
			"EqualMin succeeds":    {actual: int32(5), min: int32(5), max: int32(9)},
			"WithinRange succeeds": {actual: int32(5), min: int32(0), max: int32(9)},
			"EqualMax succeeds":    {actual: int32(5), min: int32(0), max: int32(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: int32(5), min: int32(6), max: int32(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: int32(10), min: int32(0), max: int32(9)},
		},
		reflect.Int64: {
			"EqualMin succeeds":    {actual: int64(5), min: int64(5), max: int64(9)},
			"WithinRange succeeds": {actual: int64(5), min: int64(0), max: int64(9)},
			"EqualMax succeeds":    {actual: int64(5), min: int64(0), max: int64(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: int64(5), min: int64(6), max: int64(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: int64(10), min: int64(0), max: int64(9)},
		},
		reflect.Uint: {
			"EqualMin succeeds":    {actual: uint(5), min: uint(5), max: uint(9)},
			"WithinRange succeeds": {actual: uint(5), min: uint(0), max: uint(9)},
			"EqualMax succeeds":    {actual: uint(5), min: uint(0), max: uint(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: uint(5), min: uint(6), max: uint(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: uint(10), min: uint(0), max: uint(9)},
		},
		reflect.Uint8: {
			"EqualMin succeeds":    {actual: uint8(5), min: uint8(5), max: uint8(9)},
			"WithinRange succeeds": {actual: uint8(5), min: uint8(0), max: uint8(9)},
			"EqualMax succeeds":    {actual: uint8(5), min: uint8(0), max: uint8(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: uint8(5), min: uint8(6), max: uint8(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: uint8(10), min: uint8(0), max: uint8(9)},
		},
		reflect.Uint16: {
			"EqualMin succeeds":    {actual: uint16(5), min: uint16(5), max: uint16(9)},
			"WithinRange succeeds": {actual: uint16(5), min: uint16(0), max: uint16(9)},
			"EqualMax succeeds":    {actual: uint16(5), min: uint16(0), max: uint16(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: uint16(5), min: uint16(6), max: uint16(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: uint16(10), min: uint16(0), max: uint16(9)},
		},
		reflect.Uint32: {
			"EqualMin succeeds":    {actual: uint32(5), min: uint32(5), max: uint32(9)},
			"WithinRange succeeds": {actual: uint32(5), min: uint32(0), max: uint32(9)},
			"EqualMax succeeds":    {actual: uint32(5), min: uint32(0), max: uint32(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: uint32(5), min: uint32(6), max: uint32(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: uint32(10), min: uint32(0), max: uint32(9)},
		},
		reflect.Uint64: {
			"EqualMin succeeds":    {actual: uint64(5), min: uint64(5), max: uint64(9)},
			"WithinRange succeeds": {actual: uint64(5), min: uint64(0), max: uint64(9)},
			"EqualMax succeeds":    {actual: uint64(5), min: uint64(0), max: uint64(5)},
			"BelowMin fails":       {expectFailurePattern: Ptr(`Expected actual value 5 to be between 6 and 9`), actual: uint64(5), min: uint64(6), max: uint64(9)},
			"AboveMax fails":       {expectFailurePattern: Ptr(`Expected actual value 10 to be between 0 and 9`), actual: uint64(10), min: uint64(0), max: uint64(9)},
		},
	}
	for kind, kindTestCases := range testCases {
		kind := kind
		kindTestCases := kindTestCases
		t.Run(kind.String(), func(t *testing.T) {
			for name, tc := range kindTestCases {
				tc := tc
				t.Run(name, func(t *testing.T) {
					if tc.expectFailurePattern != nil {
						defer VerifyTestOutcome(t, ExpectFailure, *tc.expectFailurePattern)
					} else {
						defer VerifyTestOutcome(t, ExpectSuccess, "")
					}
					With(NewMockT(t)).Verify(tc.actual).Will(BeBetween(tc.min, tc.max)).OrFail()
				})
			}
		})
	}
	t.Run("MaxTypeMismatches", func(t *testing.T) {
		t.Parallel()
		defer VerifyTestOutcome(t, ExpectFailure, `Expected actual value to be of type 'int64', but it is of type 'int'`)
		With(NewMockT(t)).Verify(1).Will(BeBetween(0, int64(9))).OrFail()
	})
	t.Run("MinTypeMismatches", func(t *testing.T) {
		t.Parallel()
		defer VerifyTestOutcome(t, ExpectFailure, `Expected actual value to be of type 'int64', but it is of type 'int'`)
		With(NewMockT(t)).Verify(1).Will(BeBetween(int64(0), 9)).OrFail()
	})
}
