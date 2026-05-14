package publicembed

type Logger struct{}

type PublicBad struct {
	Logger // want "avoid embedding types in public structs"
	Name   string
}

type PublicGood struct {
	Log Logger
}

type privateOK struct {
	Logger
}
