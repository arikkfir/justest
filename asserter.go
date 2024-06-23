package justest

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	"github.com/arikkfir/justest/internal"
)

const SlowFactorEnvVarName = "JUSTEST_SLOW_FACTOR"

//go:noinline
func With(t T) VerifyOrEnsure {
	if t == nil {
		panic("given T instance must not be nil")
	}
	GetHelper(t).Helper()
	return &verifier{t: t}
}

type VerifyOrEnsure interface {
	// EnsureThat adds a description to the upcoming assertion, which will be printed in case it fails.
	EnsureThat(string, ...any) Ensurer

	// Deprecated: Ensure is a synonym for EnsureThat.
	Ensure(string, ...any) Ensurer

	// VerifyThat starts an assertion without a description.
	VerifyThat(actuals ...any) Asserter

	// Deprecated: Verify is a synonym for VerifyThat.
	Verify(actuals ...any) Asserter
}

type Ensurer interface {
	ByVerifying(actuals ...any) Asserter
}

type verifier struct {
	t    T
	desc string
}

//go:noinline
func (v *verifier) EnsureThat(format string, args ...any) Ensurer {
	GetHelper(v.t).Helper()
	v.desc = fmt.Sprintf(format, args...)
	return v
}

//go:noinline
func (v *verifier) Ensure(format string, args ...any) Ensurer {
	GetHelper(v.t).Helper()
	v.desc = fmt.Sprintf(format, args...)
	return v
}

//go:noinline
func (v *verifier) ByVerifying(actuals ...any) Asserter {
	GetHelper(v.t).Helper()
	return &asserter{t: v.t, desc: v.desc, actuals: actuals}
}

//go:noinline
func (v *verifier) VerifyThat(actuals ...any) Asserter {
	GetHelper(v.t).Helper()
	return &asserter{t: v.t, desc: v.desc, actuals: actuals}
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
	Now()
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
func (a *assertion) Now() {
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
		format = fmt.Sprintf("Assertion that %s failed: %s", a.desc, format)
	}

	if a.contain {
		panic(internal.FormatAndArgs{Format: &format, Args: args})
	} else {
		caller := internal.CallerAt(1)
		callerFunction, callerFile, callerLine := caller.Location()

		// Check if direct caller is from within the "justest" package; if NOT (application test code) print the caller
		if internalCall, err := regexp.MatchString(`.*/arikkfir/justest\.`, callerFunction); err != nil {
			panic(fmt.Errorf("illegal regexp matching: %+v", err))
		} else if !internalCall {
			// Direct caller is NOT from the "justest" package; thus we also print the caller, in addition to the
			// location of the actual assertion (which is always printed)
			format = format + "\n%s:%d --> %s"
			args = append(args, filepath.Base(callerFile), callerLine, indentIfMultiLine(readSourceAt(callerFile, callerLine)))
		}

		// Always print the assertion location
		format = format + "\n%s:%d --> %s"
		args = append(args, filepath.Base(a.location.File), a.location.Line, indentIfMultiLine(a.location.Source))

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
