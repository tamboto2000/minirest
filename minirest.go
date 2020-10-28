package minirest

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/tamboto2000/common/http/response"

	"github.com/gorilla/schema"

	"github.com/gorilla/mux"
)

type Minirest struct {
	router *mux.Router
	port   string
	ip     string
}

func New() *Minirest {
	return &Minirest{
		router: new(mux.Router),
	}
}

func (mn *Minirest) RunServer() {
	var ipAndPort string
	if mn.ip != "" {
		ipAndPort += mn.ip
	}

	if mn.port != "" {
		ipAndPort += ":" + mn.port
	}

	log.Fatal(http.ListenAndServe(ipAndPort, mn.router))
}

func (mn *Minirest) ServeIP(ip string) {
	mn.ip = ip
}

func (mn *Minirest) ServePort(port string) {
	mn.port = port
}

func (mn *Minirest) HandleFunc(method string, path string, calback interface{}) {
	mn.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		writer := response.NewWriter(w)
		decoder := schema.NewDecoder()
		param := reflect.New(reflect.ValueOf(calback).Type().In(0))

		//if HTTP Method is GET or DELETE
		method := strings.ToLower(r.Method)
		if method == "get" || method == "delete" {
			//extract path values
			pathVals := mux.Vars(r)
			queryVars := r.URL.Query()
			//append pathVals to queryVars
			for key, val := range pathVals {
				queryVars[key] = []string{val}
			}

			if err := decoder.Decode(param.Interface(), queryVars); err != nil {
				writer.BadRequest(err.Error())
				return
			}
		}

		//if HTTP Method is POST, PUT, or PATCh
		if method == "post" || method == "put" || method == "patch" {
			if err := json.NewDecoder(r.Body).Decode(param.Interface()); err != nil {
				writer.BadRequest(err.Error())
				return
			}
		}

		//return must be a ResponseBuilder
		returns := reflect.ValueOf(calback).Call([]reflect.Value{param.Elem()})
		responseBuilder := returns[0].Interface().(*ResponseBuilder)
		w.WriteHeader(responseBuilder.statusCode)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBuilder.data)

	}).Methods(method)
}
