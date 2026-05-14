package enumstart

type Operation int

const (
	Add Operation = iota // want "enum should start at a non-zero value"
	Subtract
	Multiply
)

type Status int

const (
	StatusOK Status = iota + 1 // OK: starts at 1
	StatusError
)

type LogOutput int

//nolint:enumstart // OK: zero value is intentional for this enum family
const (
	LogToStdout LogOutput = iota
	LogToFile
	LogToRemote
)

const single = 42 // OK: not a group
