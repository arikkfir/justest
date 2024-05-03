package justest

type Matcher interface {
	Assert(t T, actuals ...any)
}

type MatcherFunc func(t T, actuals ...any)

//go:noinline
func (f MatcherFunc) Assert(t T, actuals ...any) {
	GetHelper(t).Helper()
	f(t, actuals...)
}
