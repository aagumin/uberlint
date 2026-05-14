package mapinit

func bad() {
	_ = map[string]int{} // want "use make for empty map initialization"
}

func good() {
	_ = make(map[string]int)
	_ = map[string]int{"a": 1}
}
