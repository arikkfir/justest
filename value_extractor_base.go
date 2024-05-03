package justest

import (
	. "github.com/arikkfir/justest/internal"
	"reflect"
)

var (
	tTypePkgPath string
	tTypeName    string
)

//go:noinline
func init() {
	rt := reflect.TypeOf((*T)(nil)).Elem()
	tTypePkgPath = rt.PkgPath()
	tTypeName = rt.Name()
}

type Extractor func(T, any) (any, bool)

type ValueExtractor map[reflect.Kind]Extractor

//go:noinline
func NewValueExtractor(defaultExtractor Extractor) ValueExtractor {
	ve := make(map[reflect.Kind]Extractor)
	ve[reflect.Invalid] = defaultExtractor
	return ve
}

//go:noinline
func NewValueExtractorWithMap(defaultExtractor Extractor, extractors map[reflect.Kind]Extractor) ValueExtractor {
	ve := NewValueExtractor(defaultExtractor)
	for k, v := range extractors {
		ve[k] = v
	}
	return ve
}

//go:noinline
func (ve ValueExtractor) ExtractValue(t T, actual any) (any, bool) {
	GetHelper(t).Helper()

	if actual == nil {
		return nil, true
	}

	v := reflect.ValueOf(actual)
	extractor, ok := ve[v.Kind()]
	if !ok {
		extractor, ok = ve[reflect.Invalid]
		if !ok {
			t.Fatalf("Default extractor is missing")
		}
	}

	return extractor(t, actual)
}

//go:noinline
func (ve ValueExtractor) MustExtractValue(t T, actual any) any {
	GetHelper(t).Helper()
	value, found := ve.ExtractValue(t, actual)
	if !found {
		t.Fatalf("Value could not be extracted from an actual of type '%T': %+v", actual, actual)
	}
	return value
}

//go:noinline
func ExtractSameValue(t T, v any) (any, bool) {
	GetHelper(t).Helper()
	return v, true
}

//go:noinline
func ExtractorUnsupported(t T, v any) (any, bool) {
	GetHelper(t).Helper()
	t.Fatalf("Unsupported actual value: %+v", v)
	panic("unreachable")
}

//go:noinline
func NewChannelExtractor(ve ValueExtractor, recurse bool) Extractor {
	return func(t T, v any) (any, bool) {
		GetHelper(t).Helper()

		msg, ok := reflect.ValueOf(v).TryRecv()
		if ok {
			if recurse {
				return ve.ExtractValue(t, msg.Interface())
			} else {
				return msg.Interface(), true
			}
		} else {
			return nil, false
		}
	}
}

//go:noinline
func NewPointerExtractor(ve ValueExtractor, recurse bool) Extractor {
	return func(t T, v any) (any, bool) {
		GetHelper(t).Helper()

		underlyingValue := reflect.ValueOf(v).Elem()
		if recurse {
			return ve.ExtractValue(t, underlyingValue.Interface())
		} else {
			return underlyingValue.Interface(), true
		}
	}
}

//go:noinline
func NewFuncExtractor(ve ValueExtractor, recurse bool) Extractor {
	return func(t T, v any) (any, bool) {
		GetHelper(t).Helper()

		funcValue := reflect.ValueOf(v)
		funcType := funcValue.Type()

		// Prepare input parameters
		var in []reflect.Value
		switch funcType.NumIn() {
		case 0:
			in = nil
		case 1:
			arg0Type := funcType.In(0)
			if arg0Type.PkgPath() == tTypePkgPath && arg0Type.Name() == tTypeName {
				in = append(in, reflect.ValueOf(t))
			} else {
				t.Fatalf("Argument of functions with one argument must be of type TT, found: %+v", arg0Type.Name())
				panic("unreachable")
			}
		default:
			t.Fatalf("Functions with more than 1 input parameter are not supported in this context: %+v", v)
			panic("unreachable")
		}

		// Prepare call & output extraction
		switch funcType.NumOut() {
		case 0:
			funcValue.Call(in)
			return nil, false
		case 1:
			returnValues := funcValue.Call(in)
			if IsErrorType(funcType.Out(0)) {
				if returnValues[0].IsNil() {
					return nil, true
				} else {
					t.Fatalf("Function failed: %+v", returnValues[0].Interface())
					panic("unreachable")
				}
			} else if recurse {
				return ve.ExtractValue(t, returnValues[0].Interface())
			} else {
				return returnValues[0].Interface(), true
			}
		case 2:
			if IsErrorType(funcType.Out(1)) {
				returnValues := funcValue.Call(in)
				if returnValues[1].IsNil() {
					if recurse {
						return ve.ExtractValue(t, returnValues[0].Interface())
					} else {
						return returnValues[0].Interface(), true
					}
				} else {
					t.Fatalf("Function failed: %+v", returnValues[1].Interface())
					panic("unreachable")
				}
			} else {
				t.Fatalf("Functions with 2 return values must return 'error' as the 2nd return value: %+v", v)
				panic("unreachable")
			}
		default:
			t.Fatalf("Functions with %d return values are not supported: %+v", funcType.NumOut(), v)
			panic("unreachable")
		}
	}
}
