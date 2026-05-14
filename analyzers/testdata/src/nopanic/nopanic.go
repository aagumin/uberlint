package nopanic

import "text/template"

func badPanic() {
	panic("something went wrong") // want "avoid panic in production code"
}

func goodError() error {
	return nil
}

func init() {
	panic("setup failed") // OK: init is allowed
}

func MustParse(s string) *template.Template {
	panic("invalid template") // OK: Must* pattern
}

type Foo struct{}

func (f *Foo) bad() {
	panic("method panic") // want "avoid panic in production code"
}
