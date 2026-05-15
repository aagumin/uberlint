package zerovar

type User struct {
	Name string
	Age  int
}

type Alias User
type MyMap map[string]int

func bad() {
	u := User{} // want "use var declaration instead"
	_ = u

	a := Alias{} // want "use var declaration instead"
	_ = a
}

func good() {
	var u User
	_ = u

	initialized := User{Name: "Alice"} // OK: has non-zero fields
	_ = initialized

	var a Alias
	_ = a
}

func goodMapNotStruct() {
	m := MyMap{} // OK: map, not struct — nil map != empty map
	_ = m
}

func goodSliceNotStruct() {
	s := []int{} // OK: slice, not struct — handled by nilslice
	_ = s
}
