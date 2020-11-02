package main

type SimpleService struct {
	Message string
}

func (sv *SimpleService) Init() {
	sv.Message = "Hello, world!"
}

type Simple2Service struct {
	Message string
}

func (sv2 *Simple2Service) Init() {
	sv2.Message = "Hello, world!"
}
