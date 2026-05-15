package localvar

var pkg = "ok"

func bad() {
	var s = "value" // want "use short variable declaration"
	_ = s

	var a, b = 1, 2 // want "use short variable declaration"
	_, _ = a, b
}

func good() {
	s := "value"
	_ = s
}

func typed() {
	var n int = 1
	_ = n
}

func noInitializer() {
	var s string
	_ = s
}
