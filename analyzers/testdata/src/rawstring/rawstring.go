package rawstring

func bad() {
	_ = "unknown name:\"test\""     // want "use raw string literal"
	_ = "path:\"/foo\" bar:\"baz\"" // want "use raw string literal"
}

func good() {
	_ = `unknown name:"test"`
	_ = "simple string"
	_ = "has `backtick` inside" // OK: can't use raw literal
	_ = "ab"                    // OK: too short, not worth it
	_ = "line1\n\"quoted\""     // OK: raw string would change the newline semantics
}
