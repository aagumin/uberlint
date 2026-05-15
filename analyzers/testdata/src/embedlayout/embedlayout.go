package embedlayout

type Base struct{}

type Bad struct {
	Name string
	Base // want "embedded fields should be grouped before named fields"
}

type BadMultiple struct {
	Name  string
	Age   int
	*Base // want "embedded fields should be grouped before named fields"
}

type Good struct {
	Base

	Name string
}

type GoodMultiple struct {
	Base
	*Other

	Name string
}

type Other struct{}
