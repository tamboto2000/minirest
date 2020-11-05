package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
