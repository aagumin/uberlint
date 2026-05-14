package constprintf

import "fmt"

const msgFormat = "hello %s"

func bad(name string) {
	_ = fmt.Sprintf("hello %s", name) // want "move printf format string to a const"
}

func good(name string) {
	_ = fmt.Sprintf(msgFormat, name)
}
