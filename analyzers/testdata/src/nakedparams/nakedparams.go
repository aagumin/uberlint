package nakedparams

func connect(host string, secure bool, timeout int) {}

func bad() {
	connect("localhost", true, 30) // want "avoid naked literal parameters"
}

func good() {
	const host = "localhost"
	const timeout = 30
	connect(host, true, timeout)
}
