package zerofields

type User struct {
	Name    string
	Age     int
	Enabled bool
	Tags    []string
}

func bad() {
	_ = User{
		Name:    "",    // want "omit zero-value field"
		Age:     0,     // want "omit zero-value field"
		Enabled: false, // want "omit zero-value field"
		Tags:    nil,   // want "omit zero-value field"
	}
}

func good() {
	_ = User{Name: "alice", Age: 42, Enabled: true}
}
