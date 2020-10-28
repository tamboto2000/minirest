package main

import "github.com/tamboto2000/minirest"

func main() {
	mns := minirest.New()
	mns.ServePort("8081")
	mns.HandleFunc("GET", "/hello/{name}", Get)
	mns.RunServer()
}

func Get(param struct{ Name string }) *minirest.ResponseBuilder {
	responseBuilder := new(minirest.ResponseBuilder)
	return responseBuilder.Ok(param.Name)
}
