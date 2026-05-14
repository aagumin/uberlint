package embedlayout

type Base struct{}

type Bad struct {
	Name string
	Base // want "embedded fields should be grouped before named fields"
}

type Good struct {
	Base

	Name string
}
