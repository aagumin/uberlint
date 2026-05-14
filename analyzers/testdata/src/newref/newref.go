package newref

type MyStruct struct {
	Name string
}

func bad() {
	_ = new(MyStruct) // want "instead of new"
}

func good() {
	_ = &MyStruct{Name: "foo"}
	_ = new(int) // OK: scalar types do not have an &T{} equivalent
	s := MyStruct{}
	_ = &s
}
