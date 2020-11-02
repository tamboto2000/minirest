package minirest

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
)

//Controller is interface for controller
//controller must have field BasePath string for sub-routing, otherwise root path will be used (/)
//if you want to register service into controller, make sure you have field with the same name as the service
//example:
// type Controller struct {
// 	UserService *UserService
// 	ItemService ItemService
// }
type Controller interface {
	//Endpoints register all endpoints to its handler
	Endpoints() *Endpoints
}

type endpoint struct {
	method   string
	path     string
	callback interface{}
}

//Endpoints register handlers its path and method
type Endpoints struct {
	endpoints []endpoint
}

func (ep *Endpoints) Add(method, path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{method, path, callback})
}

func handleWithoutPostBody(callback interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, pathVars httprouter.Params) {
		writer := new(ResponseBuilder)
		decoder := schema.NewDecoder()
		var params []reflect.Value
		//get all parameters in callback
		m := reflect.ValueOf(callback)
		for i := 0; i < m.Type().NumIn(); i++ {
			params = append(params, reflect.New(m.Type().In(i)))
		}

		//match all path variables with callback parameters index and assign callback param
		//right now the supported type are only int and float64
		for i, pair := range pathVars {
			param := params[i].Elem()

			//exclude type struct as it's only for filtering (url query)
			if param.Kind() == reflect.Struct {

				continue
			}

			//handle pointer param
			if param.Kind() == reflect.Ptr {
				if param.Type().Elem().Kind() == reflect.Struct {

					continue
				}

				param.Set(reflect.New(param.Type().Elem()))

				if param.Type().Elem().Kind() == reflect.String {
					param.Elem().SetString(pair.Value)
					params[i] = param

					continue
				}

				if param.Type().Elem().Kind() == reflect.Int {
					if intVal, err := stringToInt(pair.Value); err == nil {
						param.Elem().SetInt(int64(intVal))
						params[i] = param

						continue
					}

					writer.BadRequest(pair.Key + " is not type int")
					writer.write(w)
					return
				}

				if param.Type().Elem().Kind() == reflect.Float64 {
					if flt64Val, err := stringToFloat64(pair.Value); err == nil {
						param.Elem().SetFloat(flt64Val)
						params[i] = param

						continue
					}

					writer.BadRequest(pair.Key + " is not type float64")
					writer.write(w)
					return
				}
			}

			if param.Kind() == reflect.String {
				param.SetString(pair.Value)
				params[i] = param

				continue
			}

			if param.Kind() == reflect.Int {
				if intVal, err := stringToInt(pair.Value); err == nil {
					param.SetInt(int64(intVal))
					params[i] = param

					continue
				}

				fmt.Println("index:", i)
				fmt.Println("val:", pair.Value)
				fmt.Println("param type:", param.Kind())
				fmt.Println("variables:", pathVars)
				writer.BadRequest(pair.Key + " is not type int")
				writer.write(w)
				return
			}

			if param.Kind() == reflect.Float64 {
				if flt64Val, err := stringToFloat64(pair.Value); err == nil {
					param.SetFloat(flt64Val)
					params[i] = param

					continue
				}

				writer.BadRequest(pair.Key + " is not type float64")
				writer.write(w)
				return
			}
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

			if v := param.Elem(); v.Kind() == reflect.Ptr {
				if v.Type().Elem().Kind() != reflect.Struct {
					continue
				}

				//v now a pointer to a struct
				//initialize the struct
				v.Set(reflect.New(v.Type().Elem()))

				if err := decoder.Decode(v.Interface(), queriesVars); err != nil {
					writer.BadRequest(err.Error())
					writer.write(w)
					return
				}

				params[i] = param.Elem()

				break
			}

			if param.Elem().Kind() == reflect.Struct {
				if err := decoder.Decode(param.Interface(), queriesVars); err != nil {
					writer.BadRequest(err.Error())
					writer.write(w)
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
		respBuilder.write(w)
	}
}

func stringToInt(data string) (int, error) {
	return strconv.Atoi(data)
}

func stringToFloat64(data string) (float64, error) {
	return strconv.ParseFloat(data, 64)
}
