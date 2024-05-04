package justest_test

import (
	"fmt"
	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

var (
	StringExtractorAddingFooPrefix = func(t T, v any) (any, bool) { return "foo: " + v.(string), true }
)

func TestValueExtractor(t *testing.T) {
	t.Parallel()
	alwaysBarExtractor := func(t T, v any) (any, bool) { return "bar", true }
	type testCase struct {
		valueExtractorFactory   func() ValueExtractor
		verifier                func(t T, ve ValueExtractor) []any
		expectedOutcome         TestOutcomeExpectation
		expectedOutcomePattern  string
		expectedVerifierResults []any
	}
	testCases := map[string]testCase{
		"Default extractor is used when no extractors have been defined": {
			valueExtractorFactory:   func() ValueExtractor { return NewValueExtractor(alwaysBarExtractor) },
			verifier:                func(t T, ve ValueExtractor) []any { return []any{ve.MustExtractValue(t, "foo")} },
			expectedOutcome:         ExpectSuccess,
			expectedVerifierResults: []any{"bar"},
		},
		"Nil actual finds nil result": {
			valueExtractorFactory: func() ValueExtractor { return NewValueExtractor(alwaysBarExtractor) },
			verifier: func(t T, ve ValueExtractor) []any {
				v, found := ve.ExtractValue(t, nil)
				return []any{v, found}
			},
			expectedOutcome:         ExpectSuccess,
			expectedVerifierResults: []any{nil, true},
		},
		"Invokes correct extractor when kind found": {
			valueExtractorFactory: func() ValueExtractor {
				return NewValueExtractorWithMap(ExtractorUnsupported, map[reflect.Kind]Extractor{
					reflect.String: StringExtractorAddingFooPrefix,
				})
			},
			verifier: func(t T, ve ValueExtractor) []any {
				v, found := ve.ExtractValue(t, "bar")
				return []any{v, found}
			},
			expectedOutcome:         ExpectSuccess,
			expectedVerifierResults: []any{"foo: bar", true},
		},
		"Default extractor when kind not found": {
			valueExtractorFactory: func() ValueExtractor {
				return NewValueExtractorWithMap(ExtractSameValue, map[reflect.Kind]Extractor{
					reflect.String: StringExtractorAddingFooPrefix,
				})
			},
			verifier: func(t T, ve ValueExtractor) []any {
				v, found := ve.ExtractValue(t, 1)
				return []any{v, found}
			},
			expectedOutcome:         ExpectSuccess,
			expectedVerifierResults: []any{1, true},
		},
		"Failure occurs when value is required and not found": {
			valueExtractorFactory:  func() ValueExtractor { return NewValueExtractor(func(t T, v any) (any, bool) { return nil, false }) },
			verifier:               func(t T, ve ValueExtractor) []any { return []any{ve.MustExtractValue(t, 1)} },
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `Value could not be extracted from an actual of type 'int': 1`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			ve := tc.valueExtractorFactory()
			if tc.verifier != nil {
				verifierResults := tc.verifier(mt, ve)
				if !cmp.Equal(tc.expectedVerifierResults, verifierResults) {
					t.Fatalf("Unexpected verifier results:\n%s", cmp.Diff(tc.expectedVerifierResults, verifierResults))
				}
			} else if tc.expectedVerifierResults != nil {
				t.Fatalf("Illegal test definition - verifier is nil, but expected verifier results are not")
			}
		})
	}
}

func TestNewChannelExtractor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		defaultExtractor         Extractor
		extractorsMap            map[reflect.Kind]Extractor
		chanProvider             func() chan any
		recurse                  bool
		expectedOutcome          TestOutcomeExpectation
		expectedOutcomePattern   string
		expectedExtractorResults []any
	}
	testCases := map[string]testCase{
		"Empty & closed channel returns nil & not-found": {
			defaultExtractor: ExtractorUnsupported,
			chanProvider: func() chan any {
				ch := make(chan any, 1)
				close(ch)
				return ch
			},
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{nil, false},
		},
		"Empty & open channel returns nil & not-found": {
			defaultExtractor:         ExtractorUnsupported,
			chanProvider:             func() chan any { return make(chan any, 1) },
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{nil, false},
		},
		"Recurse properly returns a found value result": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			chanProvider: func() chan any {
				ch := make(chan any, 1)
				ch <- "bar"
				return ch
			},
			recurse:                  true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"foo: bar", true},
		},
		"Recurse properly returns a nil & not-found result": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: func(t T, v any) (any, bool) { return nil, false }},
			chanProvider: func() chan any {
				ch := make(chan any, 1)
				ch <- "bar"
				return ch
			},
			recurse:                  true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{nil, false},
		},
		"Recurse properly propagates extraction failure": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap: map[reflect.Kind]Extractor{
				reflect.String: func(t T, v any) (any, bool) {
					t.Fatalf("Extractor fails")
					panic("unreachable")
				},
			},
			chanProvider: func() chan any {
				ch := make(chan any, 1)
				ch <- "bar"
				return ch
			},
			recurse:         true,
			expectedOutcome: ExpectFailure,
		},
		"No recurse returns raw value": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			chanProvider: func() chan any {
				ch := make(chan any, 1)
				ch <- "bar"
				return ch
			},
			recurse:                  false,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"bar", true}, // Will show that the string extractor above wasn't called
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			extractor := NewChannelExtractor(NewValueExtractorWithMap(tc.defaultExtractor, tc.extractorsMap), tc.recurse)
			actual, found := extractor(mt, tc.chanProvider())
			if !cmp.Equal(tc.expectedExtractorResults, []any{actual, found}) {
				t.Fatalf("Incorrect extractor results:\n%s", cmp.Diff(tc.expectedExtractorResults, []any{actual, found}))
			}
		})
	}
}

func TestNewPointerExtractor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		defaultExtractor         Extractor
		extractorsMap            map[reflect.Kind]Extractor
		actual                   any
		recurse                  bool
		expectedOutcome          TestOutcomeExpectation
		expectedOutcomePattern   string
		expectedExtractorResults []any
	}
	testCases := map[string]testCase{
		"Recurse properly extracts found non-nil result": {
			defaultExtractor:         ExtractorUnsupported,
			extractorsMap:            map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actual:                   Ptr[string]("bar"),
			recurse:                  true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"foo: bar", true},
		},
		"Recurse properly extracts nil & not-found result": {
			defaultExtractor:         ExtractorUnsupported,
			extractorsMap:            map[reflect.Kind]Extractor{reflect.String: func(t T, v any) (any, bool) { return nil, false }},
			actual:                   Ptr[string]("bar"),
			recurse:                  true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{nil, false},
		},
		"Recurse properly propagates extraction failure": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap: map[reflect.Kind]Extractor{
				reflect.String: func(t T, v any) (any, bool) {
					t.Fatalf("Extractor failed")
					panic("unreachable")
				},
			},
			actual:                 Ptr[string]("bar"),
			recurse:                true,
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^Extractor failed$`,
		},
		"No recurse returns raw result": {
			defaultExtractor:         ExtractorUnsupported,
			actual:                   Ptr[string]("bar"),
			recurse:                  false,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"bar", true},
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			extractor := NewPointerExtractor(NewValueExtractorWithMap(tc.defaultExtractor, tc.extractorsMap), tc.recurse)
			actual, found := extractor(mt, tc.actual)
			if !cmp.Equal(tc.expectedExtractorResults, []any{actual, found}) {
				t.Fatalf("Incorrect extractor results:\n%s", cmp.Diff(tc.expectedExtractorResults, []any{actual, found}))
			}
		})
	}
}

func TestNewFuncExtractor(t *testing.T) {
	t.Parallel()
	type testCase struct {
		defaultExtractor         Extractor
		extractorsMap            map[reflect.Kind]Extractor
		actualProvider           func(*testCase) any
		recurse                  bool
		expectedValue            any
		expectedFound            bool
		called                   bool
		wantCalled               bool
		wantErr                  bool
		expectedOutcome          TestOutcomeExpectation
		expectedOutcomePattern   string
		expectedExtractorResults []any
	}
	testCases := map[string]testCase{
		"func() called & returns nil, false": {
			defaultExtractor:         ExtractorUnsupported,
			extractorsMap:            map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider:           func(tc *testCase) any { return func() { tc.called = true } },
			wantCalled:               true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{nil, false},
		},
		"func() string called & returns result": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() string {
					tc.called = true
					return "bar"
				}
			},
			wantCalled:               true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"bar", true},
		},
		"func() string returns recursed result": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() string {
					tc.called = true
					return "bar"
				}
			},
			recurse:                  true,
			wantCalled:               true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"foo: bar", true},
		},
		"func() error propagates returned error": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() error {
					tc.called = true
					return fmt.Errorf("foobar")
				}
			},
			wantCalled:             true,
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: "^Function failed: foobar$",
		},
		"func() error returns nil, true": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() error {
					tc.called = true
					return nil
				}
			},
			wantCalled:               true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{nil, true},
		},
		"func() (string, error) returns result": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() (string, error) {
					tc.called = true
					return "bar", nil
				}
			},
			wantCalled:               true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"bar", true},
		},
		"func() (string, error) returns recursed result": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() (string, error) {
					tc.called = true
					return "bar", nil
				}
			},
			recurse:                  true,
			wantCalled:               true,
			expectedOutcome:          ExpectSuccess,
			expectedExtractorResults: []any{"foo: bar", true},
		},
		"func() (string, error) propagates returned error": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() (string, error) {
					tc.called = true
					return "bar", fmt.Errorf("expected failure")
				}
			},
			wantCalled:             true,
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: "^Function failed: expected failure$",
		},
		"func() (string, int) fails because it returns more than one value": {
			defaultExtractor: ExtractorUnsupported,
			extractorsMap:    map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider: func(tc *testCase) any {
				return func() (string, int) {
					tc.called = true
					return "bar", 2
				}
			},
			recurse:                true,
			wantCalled:             false,
			expectedOutcome:        ExpectFailure,
			expectedOutcomePattern: `^Functions with 2 return values must return 'error' as the 2nd return value: .+$`,
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer func() {
				GetHelper(t).Helper()
				if !t.Failed() && tc.wantCalled && !tc.called {
					t.Fatalf("Expected function to be called, but it was not")
				}
			}()
			defer VerifyTestOutcome(t, tc.expectedOutcome, tc.expectedOutcomePattern)
			extractor := NewFuncExtractor(NewValueExtractorWithMap(tc.defaultExtractor, tc.extractorsMap), tc.recurse)
			actual, found := extractor(mt, tc.actualProvider(&tc))
			if !cmp.Equal(tc.expectedExtractorResults, []any{actual, found}) {
				t.Fatalf("Incorrect extractor results:\n%s", cmp.Diff(tc.expectedExtractorResults, []any{actual, found}))
			}
		})
	}
}
