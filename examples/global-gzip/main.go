package main

import (
	"github.com/tamboto2000/minirest"
)

func main() {
	mn := minirest.New()
	mn.Gzip = true
	mn.AddController(new(Controller))
	mn.ServePort("8082")
	mn.RunServer()
}
