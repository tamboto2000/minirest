package main

import "github.com/tamboto2000/minirest"

func main() {
	mns := minirest.New()
	mns.ServePort("8081")
	mns.HandleFunc("GET", "/hello/{name}", Get)
	mns.RunServer()
}

type Person struct {
	Name     string
	Birthday string
	Gender   string
}

func Get(name string, filter Person) *minirest.ResponseBuilder {
	responseBuilder := new(minirest.ResponseBuilder)
	return responseBuilder.Ok(Person{
		Name:     name,
		Birthday: filter.Birthday,
		Gender:   filter.Gender,
	})
}
