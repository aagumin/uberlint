package chansize

var n = 10

const buffered = 2

func bad() {
	_ = make(chan int, 64)        // want "channel size should be 0 or 1"
	_ = make(chan int, 100)       // want "channel size should be 0 or 1"
	_ = make(chan int, 2)         // want "channel size should be 0 or 1"
	_ = make(chan string, 1+1)    // want "channel size should be 0 or 1"
	_ = make(chan bool, buffered) // want "channel size should be 0 or 1"
}

func good() {
	_ = make(chan int)    // OK: unbuffered
	_ = make(chan int, 0) // OK: unbuffered explicit
	_ = make(chan int, 1) // OK: size 1
	_ = make([]int, 2)    // OK: not a channel
}

func dynamic() {
	_ = make(chan int, n) // OK: dynamic size, not lintable
}
