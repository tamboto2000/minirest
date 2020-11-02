package minirest

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Minirest struct {
	services    map[string]Service
	controllers map[string]Controller
	router      *httprouter.Router
	port        string
	ip          string
}

type keyVal struct {
	key string
	val string
}

//New initiate new Minirest
func New() *Minirest {
	return &Minirest{
		services:    make(map[string]Service),
		controllers: make(map[string]Controller),
		router:      httprouter.New(),
	}
}

//RunServer run http server
func (mn *Minirest) RunServer() {
	var addr string
	if mn.ip != "" {
		addr += mn.ip
	}

	if mn.port != "" {
		addr += ":" + mn.port
	}

	log.Fatal(http.ListenAndServe(addr, mn.router))
}

//ServeIP set http server IP
func (mn *Minirest) ServeIP(ip string) {
	mn.ip = ip
}

//ServePort set http server port
func (mn *Minirest) ServePort(port string) {
	mn.port = port
}

//AddService add service
//service must be pointer to struct
func (mn *Minirest) AddService(service Service) {
	val := reflect.ValueOf(service)
	service.Init()
	servName := strings.Split(val.Type().String(), ".")
	mn.services[servName[len(servName)-1]] = service
}

//AddController add controller
//you can link services srv into controller, see Controller and AddService for more information
//if service is not registered, it will automatically register it
func (mn *Minirest) AddController(controller Controller, srv ...Service) {
	val := reflect.ValueOf(controller)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		val = reflect.New(val.Type()).Elem()
	}

	//if srv not nil, link services to controller
	if srv != nil || len(srv) != 0 {
		for _, s := range srv {
			sname := strings.Split(reflect.ValueOf(s).Type().String(), ".")
			snamestr := sname[len(sname)-1]
			f := val.FieldByName(snamestr)
			if !f.IsValid() {
				panic(snamestr + " cannot be linked to controller: no match field exist")
			}

			if regsrv, ok := mn.services[snamestr]; ok {
				f.Set(reflect.ValueOf(regsrv))
			} else {
				mn.AddService(s)
			}
		}
	}

	//call controller.Endpoints and register all endpoints
	endpoints := controller.Endpoints()
	basePath := val.FieldByName("BasePath").Interface().(string)
	for _, endpoint := range endpoints.endpoints {
		method := strings.ToLower(endpoint.method)
		if method == "get" {
			mn.router.GET(basePath+endpoint.path, handleWithoutPostBody(endpoint.callback))
		}
	}

	ctrlName := strings.Split(val.Type().String(), ".")
	mn.controllers[ctrlName[len(ctrlName)-1]] = controller
}

// //HandleFunc assign callback to an endpoint/path
// func (mn *Minirest) HandleFunc(method string, path string, callback interface{}) {
// 	mn.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
// 		//right now the supported method without post body or form is GET and DELETE
// 		if strings.ToLower(r.Method) == "get" || strings.ToLower(r.Method) == "delete" {
// 			mn.callbackWithoutPostBody(callback, w, r)
// 			return
// 		}

// 		//right now the supported method with post body or form is POST, PUT, and PATCH
// 		if strings.ToUpper(r.Method) == "post" || strings.ToUpper(r.Method) == "put" || strings.ToUpper(r.Method) == "patch" {
// 			mn.callbackWithPostBody(callback, w, r)
// 			return
// 		}
// 	}).Methods(method)
// }

// //wrapper for callback without post body
// func (mn *Minirest) callbackWithoutPostBody(callback interface{}, w http.ResponseWriter, r *http.Request) {
// 	writer := new(ResponseBuilder)
// 	decoder := schema.NewDecoder()
// 	var params []reflect.Value
// 	//get all parameters in callback
// 	m := reflect.ValueOf(callback)
// 	for i := 0; i < m.Type().NumIn(); i++ {
// 		params = append(params, reflect.New(m.Type().In(i)))
// 	}

// 	//extract path variables
// 	var pathVars []keyVal

// 	for key, val := range mux.Vars(r) {
// 		pathVars = append(pathVars, keyVal{key, val})
// 	}

// 	paramsFromContext := httprouter.ParamsFromContext(r.Context())
// 	fmt.Println(paramsFromContext)

// 	//match all path variables with callback parameters index and assign callback param
// 	//right now the supported type are only int and float64
// 	for i, pair := range pathVars {
// 		param := params[i].Elem()

// 		//exclude type struct as it's only for filtering (url query)
// 		if param.Kind() == reflect.Struct {

// 			continue
// 		}

// 		//handle pointer param
// 		if param.Kind() == reflect.Ptr {
// 			if param.Type().Elem().Kind() == reflect.Struct {

// 				continue
// 			}

// 			param.Set(reflect.New(param.Type().Elem()))

// 			if param.Type().Elem().Kind() == reflect.String {
// 				param.Elem().SetString(pair.val)
// 				params[i] = param

// 				continue
// 			}

// 			if param.Type().Elem().Kind() == reflect.Int {
// 				if intVal, err := stringToInt(pair.val); err == nil {
// 					param.Elem().SetInt(int64(intVal))
// 					params[i] = param

// 					continue
// 				}

// 				writer.BadRequest(pair.key + " is not type int")
// 				writer.write(w)
// 				return
// 			}

// 			if param.Type().Elem().Kind() == reflect.Float64 {
// 				if flt64Val, err := stringToFloat64(pair.val); err == nil {
// 					param.Elem().SetFloat(flt64Val)
// 					params[i] = param

// 					continue
// 				}

// 				writer.BadRequest(pair.key + " is not type float64")
// 				writer.write(w)
// 				return
// 			}
// 		}

// 		if param.Kind() == reflect.String {
// 			param.SetString(pair.val)
// 			params[i] = param

// 			continue
// 		}

// 		if param.Kind() == reflect.Int {
// 			if intVal, err := stringToInt(pair.val); err == nil {
// 				param.SetInt(int64(intVal))
// 				params[i] = param

// 				continue
// 			}

// 			fmt.Println("index:", i)
// 			fmt.Println("val:", pair.val)
// 			fmt.Println("param type:", param.Kind())
// 			fmt.Println("variables:", pathVars)
// 			writer.BadRequest(pair.key + " is not type int")
// 			writer.write(w)
// 			return
// 		}

// 		if param.Kind() == reflect.Float64 {
// 			if flt64Val, err := stringToFloat64(pair.val); err == nil {
// 				param.SetFloat(flt64Val)
// 				params[i] = param

// 				continue
// 			}

// 			writer.BadRequest(pair.key + " is not type float64")
// 			writer.write(w)
// 			return
// 		}
// 	}

// 	//extract url queries
// 	queriesVars := r.URL.Query()
// 	//iterate callback params, find the one with type struct. Note that once param
// 	//with type struct is found, iteration will stop and parse the url queries to it
// 	//so there will be only one param with type struct is allowed as queries
// 	for i, param := range params {
// 		//if param can set, then it must be path variable, skip!
// 		if param.CanSet() {
// 			continue
// 		}

// 		if v := param.Elem(); v.Kind() == reflect.Ptr {
// 			if v.Type().Elem().Kind() != reflect.Struct {
// 				continue
// 			}

// 			//v now a pointer to a struct
// 			//initialize the struct
// 			v.Set(reflect.New(v.Type().Elem()))

// 			if err := decoder.Decode(v.Interface(), queriesVars); err != nil {
// 				writer.BadRequest(err.Error())
// 				writer.write(w)
// 				return
// 			}

// 			params[i] = param.Elem()

// 			break
// 		}

// 		if param.Elem().Kind() == reflect.Struct {
// 			if err := decoder.Decode(param.Interface(), queriesVars); err != nil {
// 				writer.BadRequest(err.Error())
// 				writer.write(w)
// 				return
// 			}

// 			params[i] = param.Elem()

// 			break
// 		}
// 	}

// 	//call callback
// 	//note that callback only can have one return value, and it must be *ResponseBuilder
// 	returns := reflect.ValueOf(callback).Call(params)
// 	respBuilder := returns[0].Interface().(*ResponseBuilder)
// 	respBuilder.write(w)
// }

// //wrapper for callback with post body, such as POST, PUT, or PATCH
// func (mn *Minirest) callbackWithPostBody(callback interface{}, w http.ResponseWriter, r *http.Request) {
// 	writer := new(ResponseBuilder)
// 	//get parameter in callback
// 	//only the first parameter are considered the real parameter,
// 	//no matter how much params you have
// 	param := reflect.New(reflect.ValueOf(callback).Type().In(0))
// 	if err := json.NewDecoder(r.Body).Decode(param.Interface()); err != nil {
// 		writer.BadRequest(err.Error())
// 		writer.write(w)
// 		return
// 	}

// 	returns := reflect.ValueOf(callback).Call([]reflect.Value{param.Elem()})
// 	respBody := returns[0].Interface().(*ResponseBuilder)
// 	respBody.write(w)

// 	r.Body.Close()
// }

// func stringToInt(data string) (int, error) {
// 	return strconv.Atoi(data)
// }

// func stringToFloat64(data string) (float64, error) {
// 	return strconv.ParseFloat(data, 64)
// }
