package zerofields

type User struct {
	Name    string
	Age     int
	Enabled bool
	Tags    []string
	Score   float64
	Offset  complex64
}

func bad() {
	_ = User{
		Name:    "",    // want "omit zero-value field"
		Age:     0,     // want "omit zero-value field"
		Enabled: false, // want "omit zero-value field"
		Tags:    nil,   // want "omit zero-value field"
		Score:   0.0,   // want "omit zero-value field"
		Offset:  0i,    // want "omit zero-value field"
	}
}

func good() {
	_ = User{Name: "alice", Age: 42, Enabled: true, Score: 1.5, Offset: 1i}
}
