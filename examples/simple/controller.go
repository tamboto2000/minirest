package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
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

func (smp *SimpleController) Post(person *Person) *minirest.ResponseBuilder {
	responseBuilder := new(minirest.ResponseBuilder)
	return responseBuilder.
		Status(200).
		Headers([][2]string{{"Hello", "World"}}).
		Body(person)
}

func (smp *SimpleController) Endpoints() *minirest.Endpoints {
	endpoints := new(minirest.Endpoints)
	endpoints.BasePath("/simple")
	endpoints.Middlewares(Middleware1, Middleware2)
	endpoints.GET("/:id/:name/:uuid", smp.Get)
	endpoints.POST("/", smp.Post)

	return endpoints
}

func Middleware1(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Println("middleware 1")
		next(w, r, p)
	})
}

func Middleware2(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmt.Println("middleware 2")
		next(w, r, p)
	})
}
