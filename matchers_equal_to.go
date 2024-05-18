package justest

import (
	"strings"

	"github.com/google/go-cmp/cmp"
)

type Comparator func(t T, expected any, actual any)

type EqualToMatcher interface {
	Matcher
	Using(comparator Comparator) EqualToMatcher
}

type equalTo struct {
	expected   []any
	comparator Comparator
}

//go:noinline
func (m *equalTo) Assert(t T, actuals ...any) {
	GetHelper(t).Helper()
	if len(m.expected) != len(actuals) {
		t.Fatalf("Unexpected difference: received %d actual values and %d expected values", len(actuals), len(m.expected))
	} else {
		for i := 0; i < len(m.expected); i++ {
			expected := m.expected[i]
			actual := actuals[i]
			m.comparator(t, expected, actual)
		}
	}
}

func (m *equalTo) Using(comparator Comparator) EqualToMatcher {
	m.comparator = comparator
	return m
}

//go:noinline
func EqualTo(expected ...any) EqualToMatcher {
	var opts []cmp.Option
	var expectedWithoutOpts []any
	for _, e := range expected {
		if _, ok := e.(cmp.Option); ok {
			opts = append(opts, e.(cmp.Option))
		} else {
			expectedWithoutOpts = append(expectedWithoutOpts, e)
		}
	}
	return &equalTo{
		expected: expectedWithoutOpts,
		comparator: func(t T, expected, actual any) {
			GetHelper(t).Helper()
			if !cmp.Equal(expected, actual, opts...) {
				t.Fatalf("Unexpected difference (\"-\" lines are expected values; \"+\" lines are actual values):\n%s", strings.TrimSpace(cmp.Diff(expected, actual, opts...)))
			}
		},
	}
}
