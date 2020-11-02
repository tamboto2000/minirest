package main

import (
	"github.com/tamboto2000/minirest"
)

func main() {
	mns := minirest.New()

	mns.AddService(new(Simple2Service))
	mns.AddController(new(SimpleController), new(SimpleService), new(Simple2Service))

	mns.ServePort("8081")
	mns.RunServer()
}
