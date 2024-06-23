package justest

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/arikkfir/justest/internal"
)

const SlowFactorEnvVarName = "JUSTEST_SLOW_FACTOR"

//go:noinline
func With(t T) VerifierAndEnsurer {
	if t == nil {
		panic("given T instance must not be nil")
	}
	GetHelper(t).Helper()
	return &verifier{t: t}
}

type VerifierAndEnsurer interface {
	Verifier
	Ensure(string, ...any) Verifier
}

type Verifier interface {
	Verify(actuals ...any) Asserter
}

type verifier struct {
	t    T
	desc string
}

//go:noinline
func (v *verifier) Ensure(format string, args ...any) Verifier {
	GetHelper(v.t).Helper()
	v.desc = fmt.Sprintf(format, args...)
	return v
}

//go:noinline
func (v *verifier) Verify(actuals ...any) Asserter {
	GetHelper(v.t).Helper()
	return &asserter{t: v.t, desc: v.desc, actuals: actuals}
}

type Asserter interface {
	Will(m Matcher) Assertion
}

type asserter struct {
	t       T
	actuals []any
	desc    string
}

//go:noinline
func (a *asserter) Will(m Matcher) Assertion {
	GetHelper(a.t).Helper()

	aa := &assertion{
		t:        a.t,
		desc:     a.desc,
		location: nearestLocation(),
		actuals:  a.actuals,
		matcher:  m,
	}

	location := nearestLocation()
	a.t.Cleanup(func() {
		if !a.t.Failed() && !aa.evaluated {
			a.t.Fatalf("An assertion was not evaluated!\n%s:%d: --> %s", filepath.Base(location.File), location.Line, location.Source)
		}
	})

	return aa
}

type Assertion interface {
	OrFail()
	For(duration time.Duration, interval time.Duration)
	Within(duration time.Duration, interval time.Duration)
}

type assertion struct {
	t         T
	location  Location
	actuals   []any
	matcher   Matcher
	contain   bool
	cleanup   []func()
	evaluated bool
	desc      string
}

//go:noinline
func (a *assertion) OrFail() {
	GetHelper(a.t).Helper()
	if a.evaluated {
		panic("assertion already evaluated")
	} else {
		a.evaluated = true
	}
	a.matcher.Assert(a, a.actuals...)
}

//go:noinline
func (a *assertion) For(duration time.Duration, interval time.Duration) {
	GetHelper(a.t).Helper()
	duration = transformDurationIfNecessary(a.t, duration)

	if a.evaluated {
		panic("assertion already evaluated")
	} else {
		a.evaluated = true
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ticking := false
	cleaningUp := false
	var failure *internal.FormatAndArgs
	succeeded := false
	tick := func() {
		GetHelper(a).Helper()

		// Notify we're no longer in a "tick"
		defer func() { ticking = false }()

		// Contain the potential "Fatal" calls from this tick as failures
		defer func() {
			if r := recover(); r != nil {
				if fa, ok := r.(internal.FormatAndArgs); ok {
					failure = &fa
				} else if !a.Failed() {
					panic(r)
				}
			} else {
				succeeded = true
			}
		}()

		// Perform cleanups for this tick
		a.cleanup = nil
		defer func() {
			cleaningUp = true
			defer func() { cleaningUp = false }()

			// TODO: decide what to do with failures during cleanups
			for i := len(a.cleanup) - 1; i >= 0; i-- {
				a.cleanup[i]()
			}
		}()

		a.matcher.Assert(a, a.actuals...)
	}

	a.contain = true
	started := time.Now()
	for {
		select {
		case <-timer.C:
			for cleaningUp {
				time.Sleep(50 * time.Millisecond)
			}
			a.contain = false
			if failure != nil {
				a.Fatalf("%s\nAssertion failed while waiting for %s", failure, duration)
			} else if !succeeded {
				a.Fatalf("Timed out after %s waiting for assertion to pass (tick never finished once)", duration)
			} else {
				return
			}
		case <-ticker.C:
			verifyNotInterrupted(a.t)
			if failure != nil {
				for cleaningUp {
					time.Sleep(50 * time.Millisecond)
				}
				a.contain = false
				a.Fatalf("%s\nAssertion failed after %s and did not pass repeatedly for %s", failure, time.Since(started), duration)
			} else if !ticking {
				ticking = true
				go tick()
			}
		}
	}
}

//go:noinline
func (a *assertion) Within(duration time.Duration, interval time.Duration) {
	GetHelper(a.t).Helper()
	duration = transformDurationIfNecessary(a.t, duration)

	if a.evaluated {
		panic("assertion already evaluated")
	} else {
		a.evaluated = true
	}

	timer := time.NewTimer(duration)
	defer timer.Stop()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ticking := false
	cleaningUp := false
	var failure *internal.FormatAndArgs
	succeeded := false
	tick := func() {
		GetHelper(a).Helper()

		// Notify we're no longer in a "tick"
		defer func() { ticking = false }()

		// Contain the potential "Fatal" calls from this tick as failures
		defer func() {
			if r := recover(); r != nil {
				if fa, ok := r.(internal.FormatAndArgs); ok {
					failure = &fa
				} else if !a.Failed() {
					panic(r)
				}
			} else {
				succeeded = true
			}
		}()

		// Perform cleanups for this tick
		a.cleanup = nil
		defer func() {
			cleaningUp = true
			defer func() { cleaningUp = false }()

			// TODO: decide what to do with failures during cleanups
			for i := len(a.cleanup) - 1; i >= 0; i-- {
				a.cleanup[i]()
			}
		}()

		a.matcher.Assert(a, a.actuals...)
	}

	a.contain = true
	started := time.Now()
	for {
		select {
		case <-timer.C:
			for cleaningUp {
				time.Sleep(50 * time.Millisecond)
			}
			if succeeded {
				return
			}

			a.contain = false
			if failure != nil {
				a.Fatalf("%s\nTimed out after %s waiting for assertion to pass", failure, time.Since(started))
			} else {
				a.Fatalf("Timed out after %s waiting for assertion to pass (tick never finished once)", duration)
			}
		case <-ticker.C:
			verifyNotInterrupted(a.t)
			if succeeded {
				for cleaningUp {
					time.Sleep(50 * time.Millisecond)
				}
				return
			} else if !ticking {
				ticking = true
				go tick()
			}
		}
	}
}

//go:noinline
func (a *assertion) Name() string {
	return a.t.Name()
}

//go:noinline
func (a *assertion) Cleanup(f func()) {
	GetHelper(a).Helper()
	if a.contain {
		a.cleanup = append(a.cleanup, f)
	} else {
		a.t.Cleanup(f)
	}
}

//go:noinline
func (a *assertion) Failed() bool {
	GetHelper(a).Helper()
	return a.t.Failed()
}

//go:noinline
func (a *assertion) Fatalf(format string, args ...any) {
	GetHelper(a).Helper()

	if a.desc != "" {
		f := strings.ToLower(string(format[0])) + format[1:]
		format = a.desc + " failed: " + f
	}

	if a.contain {
		panic(internal.FormatAndArgs{Format: &format, Args: args})
	} else {
		caller := internal.CallerAt(1)
		callerFunction, callerFile, callerLine := caller.Location()

		format = format + "\n%s:%d --> %s"
		if matches, err := regexp.MatchString(`.*/arikkfir/justest\.`, callerFunction); err != nil {
			panic(fmt.Errorf("illegal regexp matching: %+v", err))
		} else if matches || (callerFile == a.location.File && callerLine == a.location.Line) {
			// Caller is "justest" internal (e.g. "a.OrFail", "a.For", "a.Within") OR the caller location equals the
			// assertion location (that can happen when assertion location has an inner function call that fails)
			args = append(args, filepath.Base(a.location.File), a.location.Line, a.location.Source)
		} else {
			// Caller is not "justest" internal - add both the assertion and the caller locations
			format = format + "\n%s:%d --> %s"
			args = append(args, filepath.Base(callerFile), callerLine, readSourceAt(callerFile, callerLine))
			args = append(args, filepath.Base(a.location.File), a.location.Line, a.location.Source)
		}
		a.t.Fatalf(format, args...)
	}
}

//go:noinline
func (a *assertion) Log(args ...any) {
	GetHelper(a).Helper()
	a.t.Log(args...)
}

//go:noinline
func (a *assertion) Logf(format string, args ...any) {
	GetHelper(a).Helper()
	a.t.Logf(format, args...)
}

//go:noinline
func (a *assertion) GetParent() T {
	return a.t
}
