package constprintf

import "fmt"

const msgFormat = "hello %s"

func bad(name string) {
	msg := "hello %s"
	_ = fmt.Sprintf(msg, name) // want "make printf format string a const"

	var other = "other %s"
	fmt.Printf(other, name) // want "make printf format string a const"

	localFormat := "local %s"
	printf(localFormat, name) // want "make printf format string a const"
}

func good(name string) {
	_ = fmt.Sprintf("inline %s", name)
	_ = fmt.Sprintf(msgFormat, name)
	fmt.Println("hello", name)

	msg := "not a format string"
	fmt.Print(msg)
}

func printf(format string, args ...any) {}
