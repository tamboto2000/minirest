package main

type SimpleService struct {
	Message        string
	Simple2Service *Simple2Service
}

func (sv *SimpleService) Init() {
	sv.Message = "(SimpleService) Hello, world!"
}

type Simple2Service struct {
	Message string
}

func (sv2 *Simple2Service) Init() {
	sv2.Message = "(Simple2Service) Hello, world!"
}
