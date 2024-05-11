package internal

import (
	"fmt"
	"regexp"
)

type FormatAndArgs struct {
	Format *string
	Args   []any
}

//go:noinline
func (f FormatAndArgs) String() string {
	if f.Format != nil {
		return fmt.Sprintf(*f.Format, f.Args...)
	} else {
		return fmt.Sprint(f.Args...)
	}
}

func (f FormatAndArgs) MatchesRegexp(re *regexp.Regexp) bool {
	return re.MatchString(f.String())
}

func (f FormatAndArgs) MatchesRegexpString(pattern string) bool {
	return f.MatchesRegexp(regexp.MustCompile(pattern))
}
