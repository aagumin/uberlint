package nilslice

func badEmptyLiteral() []int {
	return []int{} // want "instead of empty slice literal"
}

func badMakeEmpty() []string {
	return make([]string, 0) // want "instead of make"
}

func badMakeEmptyWithCapacity() []byte {
	return make([]byte, 0, 16) // want "instead of make"
}

func good() []int {
	return nil
}

func goodNonEmpty() []int {
	return []int{1, 2, 3}
}

type response struct {
	items []int
}

func goodStructLiteral() response {
	return response{items: []int{}}
}

func goodMakeForField() response {
	return response{items: make([]int, 0)}
}

func goodAllocatedBuffer() []byte {
	return make([]byte, 16)
}
