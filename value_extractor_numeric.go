package justest

import (
	"cmp"
	"reflect"
)

var (
	NumericValueExtractor = NewNumericValueExtractor()
)

//go:noinline
func NewNumericValueExtractor() ValueExtractor {
	sve := NewValueExtractor(ExtractorUnsupported)
	sve[reflect.Chan] = NewChannelExtractor(sve, true)
	sve[reflect.Float32] = ExtractSameValue
	sve[reflect.Float64] = ExtractSameValue
	sve[reflect.Func] = NewFuncExtractor(sve, true)
	sve[reflect.Int] = ExtractSameValue
	sve[reflect.Int8] = ExtractSameValue
	sve[reflect.Int16] = ExtractSameValue
	sve[reflect.Int32] = ExtractSameValue
	sve[reflect.Int64] = ExtractSameValue
	sve[reflect.Pointer] = NewPointerExtractor(sve, true)
	sve[reflect.Uint] = ExtractSameValue
	sve[reflect.Uint8] = ExtractSameValue
	sve[reflect.Uint16] = ExtractSameValue
	sve[reflect.Uint32] = ExtractSameValue
	sve[reflect.Uint64] = ExtractSameValue
	return sve
}

// getNumericCompareFuncFor returns a reflection wrapper for the [cmp.Compare] function with the correct generic type
// for the numeric type provided.
//
//go:noinline
func getNumericCompareFuncFor(t T, v any) reflect.Value {
	GetHelper(t).Helper()
	switch v.(type) {
	case int:
		return reflect.ValueOf(cmp.Compare[int])
	case int8:
		return reflect.ValueOf(cmp.Compare[int8])
	case int16:
		return reflect.ValueOf(cmp.Compare[int16])
	case int32:
		return reflect.ValueOf(cmp.Compare[int32])
	case int64:
		return reflect.ValueOf(cmp.Compare[int64])
	case uint:
		return reflect.ValueOf(cmp.Compare[uint])
	case uint8:
		return reflect.ValueOf(cmp.Compare[uint8])
	case uint16:
		return reflect.ValueOf(cmp.Compare[uint16])
	case uint32:
		return reflect.ValueOf(cmp.Compare[uint32])
	case uint64:
		return reflect.ValueOf(cmp.Compare[uint64])
	case float32:
		return reflect.ValueOf(cmp.Compare[float32])
	case float64:
		return reflect.ValueOf(cmp.Compare[float64])
	default:
		t.Fatalf("Type '%T' of actual '%+v' does not have a defined comparison function", v, v)
		panic("unreachable")
	}
}
