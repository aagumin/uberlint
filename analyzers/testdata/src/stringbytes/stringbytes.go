package stringbytes

import "io"

func bad(w io.Writer) {
	for i := 0; i < 10; i++ {
		w.Write([]byte("Hello world")) // want "avoid repeated string-to-byte conversion"
	}
}

func badRange(w io.Writer, values []int) {
	for range values {
		w.Write([]byte("Hello range")) // want "avoid repeated string-to-byte conversion"
	}
}

func good(w io.Writer) {
	data := []byte("Hello world")
	for i := 0; i < 10; i++ {
		w.Write(data)
	}
}

func outsideLoop(w io.Writer) {
	w.Write([]byte("one time")) // OK: not in a loop
}

func dynamicVar(w io.Writer, s string) {
	for i := 0; i < 10; i++ {
		w.Write([]byte(s)) // OK: dynamic, not a literal
	}
}
