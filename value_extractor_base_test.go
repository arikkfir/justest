package justest_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/arikkfir/justest"
	. "github.com/arikkfir/justest/internal"
)

var (
	StringExtractorAddingFooPrefix = func(t T, v any) (any, bool) { return "foo: " + v.(string), true }
)

func TestValueExtractor(t *testing.T) {
	t.Parallel()
	alwaysBarExtractor := func(t T, v any) (any, bool) { return "bar", true }
	type testCase struct {
		valueExtractorFactory func() ValueExtractor
		verifier              func(t T, ve ValueExtractor) []any
		outcomeVerifier       TestOutcomeVerifier
		expectedResults       []any
	}
	testCases := map[string]testCase{
		"Default extractor is used when no extractors have been defined": {
			valueExtractorFactory: func() ValueExtractor { return NewValueExtractor(alwaysBarExtractor) },
			verifier:              func(t T, ve ValueExtractor) []any { return []any{ve.MustExtractValue(t, "foo")} },
			outcomeVerifier:       SuccessVerifier(),
			expectedResults:       []any{"bar"},
		},
		"Nil actual finds nil result": {
			valueExtractorFactory: func() ValueExtractor { return NewValueExtractor(alwaysBarExtractor) },
			verifier: func(t T, ve ValueExtractor) []any {
				v, found := ve.ExtractValue(t, nil)
				return []any{v, found}
			},
			outcomeVerifier: SuccessVerifier(),
			expectedResults: []any{nil, true},
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
			outcomeVerifier: SuccessVerifier(),
			expectedResults: []any{"foo: bar", true},
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
			outcomeVerifier: SuccessVerifier(),
			expectedResults: []any{1, true},
		},
		"Failure occurs when value is required and not found": {
			valueExtractorFactory: func() ValueExtractor { return NewValueExtractor(func(t T, v any) (any, bool) { return nil, false }) },
			verifier:              func(t T, ve ValueExtractor) []any { return []any{ve.MustExtractValue(t, 1)} },
			outcomeVerifier:       FailureVerifier(`Value could not be extracted from an actual of type 'int': 1`),
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.outcomeVerifier)
			ve := tc.valueExtractorFactory()
			if tc.verifier != nil {
				verifierResults := tc.verifier(mt, ve)
				if !cmp.Equal(tc.expectedResults, verifierResults) {
					t.Fatalf("Unexpected verifier results:\n%s", cmp.Diff(tc.expectedResults, verifierResults))
				}
			} else if tc.expectedResults != nil {
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
		expectedOutcome          TestOutcomeVerifier
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
			expectedOutcome:          SuccessVerifier(),
			expectedExtractorResults: []any{nil, false},
		},
		"Empty & open channel returns nil & not-found": {
			defaultExtractor:         ExtractorUnsupported,
			chanProvider:             func() chan any { return make(chan any, 1) },
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome: FailureVerifier(`^Extractor fails$`),
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
			expectedOutcome:          SuccessVerifier(),
			expectedExtractorResults: []any{"bar", true}, // Will show that the string extractor above wasn't called
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.expectedOutcome)
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
		expectedOutcome          TestOutcomeVerifier
		expectedExtractorResults []any
	}
	testCases := map[string]testCase{
		"Recurse properly extracts found non-nil result": {
			defaultExtractor:         ExtractorUnsupported,
			extractorsMap:            map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actual:                   Ptr[string]("bar"),
			recurse:                  true,
			expectedOutcome:          SuccessVerifier(),
			expectedExtractorResults: []any{"foo: bar", true},
		},
		"Recurse properly extracts nil & not-found result": {
			defaultExtractor:         ExtractorUnsupported,
			extractorsMap:            map[reflect.Kind]Extractor{reflect.String: func(t T, v any) (any, bool) { return nil, false }},
			actual:                   Ptr[string]("bar"),
			recurse:                  true,
			expectedOutcome:          SuccessVerifier(),
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
			actual:          Ptr[string]("bar"),
			recurse:         true,
			expectedOutcome: FailureVerifier(`^Extractor failed$`),
		},
		"No recurse returns raw result": {
			defaultExtractor:         ExtractorUnsupported,
			actual:                   Ptr[string]("bar"),
			recurse:                  false,
			expectedOutcome:          SuccessVerifier(),
			expectedExtractorResults: []any{"bar", true},
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mt := NewMockT(t)
			defer mt.Verify(tc.expectedOutcome)
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
		expectedOutcome          TestOutcomeVerifier
		expectedExtractorResults []any
	}
	testCases := map[string]testCase{
		"func() called & returns nil, false": {
			defaultExtractor:         ExtractorUnsupported,
			extractorsMap:            map[reflect.Kind]Extractor{reflect.String: StringExtractorAddingFooPrefix},
			actualProvider:           func(tc *testCase) any { return func() { tc.called = true } },
			wantCalled:               true,
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome:          SuccessVerifier(),
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
			wantCalled:      true,
			expectedOutcome: FailureVerifier(`^Function failed: foobar$`),
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
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome:          SuccessVerifier(),
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
			expectedOutcome:          SuccessVerifier(),
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
			wantCalled:      true,
			expectedOutcome: FailureVerifier(`^Function failed: expected failure$`),
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
			recurse:         true,
			wantCalled:      false,
			expectedOutcome: FailureVerifier(`^Functions with 2 return values must return 'error' as the 2nd return value: .+$`),
		},
	}
	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				GetHelper(t).Helper()
				if !t.Failed() && tc.wantCalled && !tc.called {
					t.Fatalf("Expected function to be called, but it was not")
				}
			}()
			mt := NewMockT(t)
			defer mt.Verify(tc.expectedOutcome)
			extractor := NewFuncExtractor(NewValueExtractorWithMap(tc.defaultExtractor, tc.extractorsMap), tc.recurse)
			actual, found := extractor(mt, tc.actualProvider(&tc))
			if !cmp.Equal(tc.expectedExtractorResults, []any{actual, found}) {
				t.Fatalf("Incorrect extractor results:\n%s", cmp.Diff(tc.expectedExtractorResults, []any{actual, found}))
			}
		})
	}
}
