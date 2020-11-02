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
				s.Init()
				f.Set(reflect.ValueOf(s))
				mn.AddService(s)
			}
		}
	}

	//call controller.Endpoints and register all endpoints
	endpoints := controller.Endpoints()
	for _, endpoint := range endpoints.endpoints {
		method := strings.ToLower(endpoint.method)
		if method == "get" {
			mn.router.GET(endpoints.basePath+endpoint.path, handleWithoutBody(endpoint.callback))
		}

		if method == "delete" {
			mn.router.DELETE(endpoints.basePath+endpoint.path, handleWithoutBody(endpoint.callback))
		}
	}

	ctrlName := strings.Split(val.Type().String(), ".")
	mn.controllers[ctrlName[len(ctrlName)-1]] = controller
}

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
