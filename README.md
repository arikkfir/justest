# justest

![Maintainer](https://img.shields.io/badge/maintainer-arikkfir-blue)
![GoVersion](https://img.shields.io/github/go-mod/go-version/arikkfir/justest.svg)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/arikkfir/justest)
[![GoReportCard](https://goreportcard.com/badge/github.com/arikkfir/justest)](https://goreportcard.com/report/github.com/arikkfir/justest)

> Go testing framework with extra sugar

This Go testing framework has the following goals:

* Play nice with `go test`
* Provide a fluent API for making assertions
* Provide a succinct yet informative error information on failures
* Make testing easier to read and write

## Attribution

This library is **heavily** inspired by [Gomega](https://github.com/onsi/gomega)
and [Ginkgo](https://github.com/onsi/ginkgo),
two **excellent** Go testing libraries that I've used extensively in the past. If you need a full **standalone** (see
below)
BDD testing framework, those are highly recommended.

The reason Justest was created is because Gomega and Ginkgo do not play great with `go test` (though Gomega does have a
`go test` integration, it's still primarily meant to be used with Ginkgo) whereas Justest is meant all along to be used
with `go test`.

## Usage

```go
package my_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	. "github.com/arikkfir/justest"
)

func TestSomething(t *testing.T) {
	
	// Simple assertions
	With(t).Verify(1).Will(BeBetween(0, 2)).OrFail()
	With(t).Verify("").Will(BeEmpty()).OrFail()
	With(t).Verify([]int{1,2,3}).Will(BeEmpty()).OrFail() // <-- This will fail!
	With(t).Verify(1).Will(BeGreaterThan(0)).OrFail()
	With(t).Verify(1).Will(BeLessThan(2)).OrFail()
	With(t).Verify("abc").Will(BeNil()).OrFail() // <-- This will fail!
	With(t).Verify(1).Will(EqualTo(1)).OrFail()
	With(t).Verify("abc").Will(EqualTo("def")).OrFail() // <-- This will fail!

	// Assert success or failure of a function (functions can have any set of return values or none at all)
	succeedingFunc := func() (string, error) { return "abc", nil }
	With(t).Verify(succeedingFunc).Will(Succeed()).OrFail() // <-- Will succeed since error return value is nil
	With(t).Verify(succeedingFunc).Will(Fail()).OrFail() // <-- Will fail since it expects error return value to be non-nil
	failingFunc := func() (string, error) { return "", fmt.Errorf("error") }
	With(t).Verify(failingFunc).Will(Succeed()).OrFail() // <-- Will fail since error return value is not nil
	With(t).Verify(failingFunc).Will(Fail()).OrFail() // <-- Will succeed since it expects error return value to be non-nil

	// Assert negation of another assertion
	With(t).Verify(1).Will(Not(EqualTo(2))).OrFail()
	
	// Assert something will **eventually** match
	// It will stop when the function succeeds (no assertion failure) or when time runs out
	With(t).Verify(func(t T) {

		// Will be invoked every 100ms until either it no longer fails or until time runs out (10s)
		With(t).Verify(2).Will(EqualTo(2)).OrFail()

	}).Will(Succeed()).Within(10*time.Second, 100*time.Millisecond)

	// Assert something will **repeatedly** match for a certain amount of time
	// It will stop on the first time the function fails
	With(t).Verify(func(t T) {

		// Will be invoked every 100ms until either it fails or until time runs out (10s)
		With(t).Verify(2).Will(EqualTo(2)).OrFail()

	}).Will(Succeed()).For(10*time.Second, 100*time.Millisecond)

	// Assert on text patterns
	With(t).Verify("abc").Will(Say("^a*c$")).OrFail()
	With(t).Verify("abc").Will(Say(regexp.MustCompile("^a*c$"))).OrFail()
	With(t).Verify([]byte("abc")).Will(Say("^a*c$")).OrFail()
}
```

## Custom matchers

You can easily create your own matchers by implementing the `Matcher` interface:

```go
package my_test

import (
	"reflect"

	. "github.com/arikkfir/justest"
)

var (
	myValueExtractor = NewValueExtractor(ExtractSameValue)
)

// BeSuperDuper returns a matcher that will ensure that each actual value passed to "With(t).Verify(...)" will be either
// "super duper" or "extra super duper", depending on the value of the `extraDuper` parameter.
func BeSuperDuper(extraDuper bool) Matcher {
	return MatcherFunc(func(t T, actuals ...any) {
		GetHelper(t).Helper()
		for _, actual := range actuals {
			v := myValueExtractor.MustExtractValue(t, actual) // This is optional, but recommended, see value extraction below
			if extraDuper {
				// Fail if it's not EXTRA super-duper
				if v.(string) != "extra super duper" {
					t.Fatalf("Value '%s' is not extra super-duper!", v)
				}
			} else {
				// Fail if it's not super-duper
				if v.(string) != "super duper" {
					t.Fatalf("Value '%s' is not super-duper!", v)
				}
			}
		}
	})
}
```

## Builtin matchers

| Matcher Name          | Description                                                                  |
|-----------------------|------------------------------------------------------------------------------|
| `BeBetween(min, max)` | Checks that all given values are between a minimum and maximum value         |
| `BeEmpty()`           | Checks that all given values are empty                                       |
| `BeGreaterThan(min)`  | Checks that all given values are greater than a minimum value                |
| `BeLessThan(max)`     | Checks that all given values are less than a maximum value                   |
| `BeNil()`             | Checks that all given values are nil                                         |
| `EqualTo(expected)`   | Checks that all given values are equal to their corresponding expected value |
| `Fail()`              | Checks that the last given value is a non-nil `error` instance               |
| `Not()`               | Checks that the given matcher fails                                          |
| `Say()`               | Checks that all given values match the given regular expression              |
| `Succeed()`           | Checks that the last given value is either nil or not an `error` instance    |

## Contributing

Please do :ok_hand: :muscle: !

See [CONTRIBUTING.md](CONTRIBUTING.md) for more information :pray:
