package constprintf

import "fmt"

const msgFormat = "hello %s"

func bad(name string) {
	_ = fmt.Sprintf("hello %s", name) // want "move printf format string to a const"
	fmt.Printf("hello %s", name)      // want "move printf format string to a const"
	printf("hello %s", name)          // want "move printf format string to a const"
}

func good(name string) {
	_ = fmt.Sprintf(msgFormat, name)
	fmt.Println("hello", name)
}

func printf(format string, args ...any) {}
