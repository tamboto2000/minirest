package minirest

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/tamboto2000/common/http/response"
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

func (mn *Minirest) HandleFunc(method string, path string, callback interface{}) {
	mn.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		//right now the supported method without post body or form is GET and DELETE
		if strings.ToLower(r.Method) == "get" || strings.ToLower(r.Method) == "delete" {
			mn.callbackWithoutPostBody(callback, w, r)
			return
		}
	}).Methods(method)
}

func (mn *Minirest) callbackWithoutPostBody(callback interface{}, w http.ResponseWriter, r *http.Request) {
	writer := response.NewWriter(w)
	decoder := schema.NewDecoder()
	var params []reflect.Value
	//get all parameters in callback
	m := reflect.ValueOf(callback)
	for i := 0; i < m.Type().NumIn(); i++ {
		params = append(params, reflect.New(m.Type().In(i)))
	}

	//extract path variables
	pathVars := mux.Vars(r)
	pathVarsIdx := 0
	//match all path variables with callback parameters index and assign callback param
	//right now the supported type are only int and float64
	for key, pathVar := range pathVars {
		//exclude type struct as it's only for filtering (url query)
		if v := params[pathVarsIdx].Elem(); v.Kind() == reflect.Struct {
			continue
		}

		if v := params[pathVarsIdx].Elem(); v.Kind() == reflect.String {
			v.SetString(pathVar)
			params[pathVarsIdx] = v
			continue
		}

		if v := params[pathVarsIdx].Elem(); v.Kind() == reflect.Int {
			if intVal, err := stringToInt(pathVar); err != nil {
				writer.BadRequest(key + " is not type int")
				return
			} else {
				v.SetInt(int64(intVal))
				params[pathVarsIdx] = v
			}
			continue
		}

		if v := params[pathVarsIdx].Elem(); v.Kind() == reflect.Float64 {
			if flt64Val, err := stringToFloat64(pathVar); err != nil {
				writer.BadRequest(key + " is not type float64")
				return
			} else {
				v.SetFloat(flt64Val)
				params[pathVarsIdx] = v
			}
			continue
		}

		pathVarsIdx++
	}

	//extract url queries
	queriesVars := r.URL.Query()
	//iterate callback params, find the one with type struct. Note that once param
	//with type struct is found, iteration will stop and parse the url queries to it
	//so there will be only one param with type struct is allowed as queries
	for i, param := range params {
		//if param can set, then it must be path variable, skip!
		if param.CanSet() {
			continue
		}

		if param.Elem().Kind() == reflect.Struct {
			if err := decoder.Decode(param.Interface(), queriesVars); err != nil {
				writer.BadRequest(err.Error())
				return
			}

			params[i] = param.Elem()

			break
		}
	}

	//call callback
	//note that callback only can have one return value, and it must be *ResponseBuilder
	returns := reflect.ValueOf(callback).Call(params)
	respBuilder := returns[0].Interface().(*ResponseBuilder)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respBuilder.statusCode)
	json.NewEncoder(w).Encode(respBuilder.data)
	return
}

func stringToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func stringToFloat64(data string) (float64, error) {
	return strconv.ParseFloat(data, 64)
}
