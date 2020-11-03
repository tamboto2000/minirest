package minirest

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
)

//Minirest is singleton for Minirest framework
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
		if method == "get" || method == "delete" {
			mn.router.GET(endpoints.basePath+endpoint.path, handleWithoutBody(endpoint.callback))
		}

		if method == "post" || method == "put" || method == "patch" {
			mn.router.POST(endpoints.basePath+endpoint.path, handleWithBody(endpoint.callback))
		}
	}

	ctrlName := strings.Split(val.Type().String(), ".")
	mn.controllers[ctrlName[len(ctrlName)-1]] = controller
}

//CORS set CORS
func (mn *Minirest) CORS(origins, creds, exposeHead, maxAge, methods, headers string) {
	mn.router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Add("Access-Control-Allow-Origin", origins)
			header.Add("Access-Control-Allow-Credentials", creds)
			header.Add("Access-Control-Expose-Headers", exposeHead)
			header.Add("Access-Control-Allow-Headers", headers)
			header.Add("Access-Control-Allow-Methods", methods)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
