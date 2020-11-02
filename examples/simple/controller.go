package main

import (
	"fmt"

	"github.com/tamboto2000/minirest"
)

type Person struct {
	ID       int
	UUID     float64
	Name     string
	Birthday string
	Gender   string
	Message  string
}

type SimpleController struct {
	SimpleService  *SimpleService
	Simple2Service *Simple2Service
}

func (smp *SimpleController) Get(id int, name *string, uuid float64, filter *Person) *minirest.ResponseBuilder {
	responseBuilder := new(minirest.ResponseBuilder)
	fmt.Println(smp.Simple2Service.Message)
	return responseBuilder.Ok(Person{
		ID:       id,
		UUID:     uuid,
		Name:     *name,
		Birthday: filter.Birthday,
		Gender:   filter.Gender,
		Message:  smp.SimpleService.Message,
	})
}

func (smp *SimpleController) Endpoints() *minirest.Endpoints {
	endpoints := new(minirest.Endpoints)
	endpoints.BasePath("/simple")
	endpoints.GET("/hello/:id/:name/:uuid", smp.Get)

	return endpoints
}