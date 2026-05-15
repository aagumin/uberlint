package main

import (
	"github.com/aagumin/uberlint"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(uberlint.NewPlugins()...)
}
