package vartype

func F() string { return "A" }

var _s string = F() // want "omit redundant type in var declaration"

// OK: type differs from expression type
var _e error = myError{}

type myError struct{}

func (myError) Error() string { return "error" }
