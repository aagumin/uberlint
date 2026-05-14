package ifaceptr

import "io"

type Reader interface {
	Read()
}

func bad(r *Reader) {} // want "avoid pointer to interface"

type BadStruct struct {
	R *Reader // want "avoid pointer to interface"
}

func good(r Reader) {}

type GoodStruct struct {
	R Reader
}

func badReturn() *io.Reader { // want "avoid pointer to interface"
	return nil
}
