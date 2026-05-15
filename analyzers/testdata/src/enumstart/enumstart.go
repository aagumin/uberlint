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

type BadOffset int

const (
	BadOffsetUnknown BadOffset = iota + 0 // want "enum should start at a non-zero value"
	BadOffsetReady
)

type LogOutput int

//nolint:enumstart // OK: zero value is intentional for this enum family
const (
	LogToStdout LogOutput = iota
	LogToFile
	LogToRemote
)

type WireStatus int

const (
	WireUnknown WireStatus = iota //nolint:enumstart // OK: wire-compatible zero value
	WireReady
)

const (
	NoIotaA = 0 // OK: not an iota enum
	NoIotaB = 1
)

const single = 42 // OK: not a group
