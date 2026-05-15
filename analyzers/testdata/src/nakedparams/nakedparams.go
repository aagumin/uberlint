package nakedparams

func connect(host string, secure bool, timeout int) {}

func bad() {
	connect("localhost", true, 30) // want "avoid naked literal parameters"
	connect("", false, 0)          // want "avoid naked literal parameters"
}

func good() {
	const host = "localhost"
	const timeout = 30
	connect(host, true, timeout)
	connect(host, true, timeout)
	connect(host, false, timeout)
	connect(host, true, timeout)
	oneLiteral(host, true)
}

func oneLiteral(host string, secure bool) {}
