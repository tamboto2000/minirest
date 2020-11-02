package main

import (
	"github.com/tamboto2000/minirest"
)

func main() {
	mns := minirest.New()
	mns.ServePort("8081")
	mns.AddService(new(SimpleService))
	mns.AddController(new(SimpleController), new(SimpleService))
	// mns.HandleFunc("GET", "/hello/{id}/{name}/{uuid}", Get)
	mns.RunServer()
}

type SimpleService struct {
	Message string
}

func (sv *SimpleService) Init() {
	sv.Message = "Hello, world!"
}

type SimpleController struct {
	BasePath      string
	SimpleService *SimpleService
}

func (smp *SimpleController) Endpoints() *minirest.Endpoints {
	smp.BasePath = "/simple"
	endpoints := new(minirest.Endpoints)
	endpoints.Add("GET", "/hello/:id/:name/:uuid", smp.Get)

	return endpoints
}

type Person struct {
	ID       int
	UUID     float64
	Name     string
	Birthday string
	Gender   string
	Message  string
}

func (smp *SimpleController) Get(id int, name *string, uuid float64, filter *Person) *minirest.ResponseBuilder {
	responseBuilder := new(minirest.ResponseBuilder)
	return responseBuilder.Ok(Person{
		ID:       id,
		UUID:     uuid,
		Name:     *name,
		Birthday: filter.Birthday,
		Gender:   filter.Gender,
		Message:  smp.SimpleService.Message,
	})
}
