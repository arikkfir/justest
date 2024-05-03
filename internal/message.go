package internal

import "fmt"

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
