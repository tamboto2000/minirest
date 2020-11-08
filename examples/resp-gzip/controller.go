package main

import (
	"github.com/tamboto2000/minirest"
)

type Controller struct{}

func (ctrl *Controller) Get() *minirest.ResponseBuilder {
	resp := &minirest.ResponseBuilder{Gzip: true}

	return resp.Ok(User{
		Name:   "Franklin Collin Tamboto",
		Age:    20,
		Gender: "Apache Attack Helicopter",
	})
}

func (ctrl *Controller) Endpoints() *minirest.Endpoints {
	endp := new(minirest.Endpoints)
	endp.BasePath("/user")
	endp.GET("/", ctrl.Get)

	return endp
}

type User struct {
	Name   string
	Age    int
	Gender string
}
