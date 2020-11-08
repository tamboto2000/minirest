package main

import (
	"github.com/tamboto2000/minirest"
)

type Controller struct{}

func (ctrl *Controller) Get() *minirest.ResponseBuilder {
	resp := new(minirest.ResponseBuilder)

	return resp.Ok(User{
		Name:   "Franklin Collin Tamboto",
		Age:    20,
		Gender: "Apache Attack Helicopter",
	})
}

func (ctrl *Controller) Post(user *User) *minirest.ResponseBuilder {
	resp := new(minirest.ResponseBuilder)

	return resp.Ok(user)
}

func (ctrl *Controller) Endpoints() *minirest.Endpoints {
	endp := &minirest.Endpoints{Gzip: true}
	endp.BasePath("/user")
	endp.GET("/", ctrl.Get)
	endp.POST("/", ctrl.Post)

	return endp
}

type User struct {
	Name   string
	Age    int
	Gender string
}
