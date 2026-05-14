package localvar

var pkg = "ok"

func bad() {
	var s = "value" // want "use short variable declaration"
	_ = s
}

func good() {
	s := "value"
	_ = s
}

func typed() {
	var n int = 1
	_ = n
}
