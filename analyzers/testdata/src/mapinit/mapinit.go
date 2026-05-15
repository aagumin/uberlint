package mapinit

type Headers map[string]string

func bad() {
	_ = map[string]int{} // want "use make for empty map initialization"
	_ = Headers{}        // want "use make for empty map initialization"
}

func good() {
	_ = make(map[string]int)
	_ = make(Headers)
	_ = map[string]int{"a": 1}
	_ = Headers{"a": "b"}
}
