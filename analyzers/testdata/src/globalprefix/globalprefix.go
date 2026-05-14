package globalprefix

import "errors"

const (
	defaultPort = 8080   // want "prefix unexported global with _"
	defaultUser = "user" // want "prefix unexported global with _"
)

const (
	_defaultPort = 8080 // OK: prefixed
	_defaultUser = "user"
)

var ErrNotFound = errors.New("not found") // OK: Err prefix for errors
var errInternal = errors.New("internal")  // OK: err prefix for errors

var someGlobal = 42 // want "prefix unexported global with _"

var _okGlobal = 42 // OK: prefixed

const ExportedConst = 100 // OK: exported
