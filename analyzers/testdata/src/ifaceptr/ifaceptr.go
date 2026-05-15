package ifaceptr

import "io"

type Reader interface {
	Read()
}

func bad(r *Reader) {} // want "avoid pointer to interface"

type BadStruct struct {
	R *Reader // want "avoid pointer to interface"
}

type BadAlias *Reader // want "avoid pointer to interface"

var badVar *Reader // want "avoid pointer to interface"

func good(r Reader) {}

type GoodStruct struct {
	R Reader
}

func badReturn() *io.Reader { // want "avoid pointer to interface"
	return nil
}

func badLocal() {
	var r *Reader // want "avoid pointer to interface"
	_ = r
}

func goodConcretePointer(r *concreteReader) {}

type concreteReader struct{}

func (*concreteReader) Read() {}
