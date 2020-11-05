package minirest

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Minirest is singleton for Minirest framework
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

// New initiate new Minirest
func New() *Minirest {
	return &Minirest{
		services:    make(map[string]Service),
		controllers: make(map[string]Controller),
		router:      httprouter.New(),
	}
}

// RunServer run http server
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

// ServeIP set http server IP
func (mn *Minirest) ServeIP(ip string) {
	mn.ip = ip
}

// ServePort set http server port
func (mn *Minirest) ServePort(port string) {
	mn.port = port
}

// AddService add service.
// Service must be pointer to struct
func (mn *Minirest) AddService(service Service) {
	val := reflect.ValueOf(service)
	service.Init()
	servName := strings.Split(val.Type().String(), ".")
	servNameStr := servName[len(servName)-1]
	if _, ok := mn.services[servNameStr]; !ok {
		mn.services[servNameStr] = service
	}
}

// LinkService link service dest with services svc.
// If service not registered, it will automatically registered.
// Service must have fields with same name as services that want to be linked.
// example:
//  type Service struct {
//  	UserService *UserService
//  	ItemService *ItemService
//  }
func (mn *Minirest) LinkService(dest Service, svcs ...Service) {
	var val reflect.Value
	sname := strings.Split(reflect.ValueOf(dest).Type().String(), ".")
	snamestr := sname[len(sname)-1]
	if svc, ok := mn.services[snamestr]; ok {
		val = reflect.ValueOf(svc).Elem()
	} else {
		dest.Init()
		mn.AddService(dest)
		val = reflect.ValueOf(dest).Elem()
	}

	for _, svc := range svcs {
		sname := strings.Split(reflect.ValueOf(svc).Type().String(), ".")
		snamestr := sname[len(sname)-1]
		f := val.FieldByName(snamestr)
		if !f.IsValid() {
			panic(snamestr + " cannot be linked to service: no match field exist")
		}

		if regsrv, ok := mn.services[snamestr]; ok {
			f.Set(reflect.ValueOf(regsrv))
		} else {
			svc.Init()
			f.Set(reflect.ValueOf(svc))
			mn.AddService(svc)
		}
	}
}

// AddController add controller.
// You can link services srv into controller, see Controller and AddService for more information.
// If service is not registered, it will automatically register it
func (mn *Minirest) AddController(controller Controller, srv ...Service) {
	val := reflect.ValueOf(controller)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	} else {
		val = reflect.New(val.Type()).Elem()
	}

	// if srv not nil, link services to controller
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

	// call controller.Endpoints and register all endpoints
	endpoints := controller.Endpoints()
	for _, endpoint := range endpoints.endpoints {
		method := strings.ToLower(endpoint.method)
		var handle httprouter.Handle
		if method == "get" || method == "delete" {
			if endpoints.middleware != nil {
				handle = endpoints.middleware.handleChain(handleWithoutBody(endpoint.callback))
			} else {
				handle = handleWithoutBody(endpoint.callback)
			}

			mn.router.Handle(endpoint.method, endpoints.basePath+endpoint.path, handle)
		}

		if method == "post" || method == "put" || method == "patch" {
			if endpoints.middleware != nil {
				handle = endpoints.middleware.handleChain(handleWithBody(endpoint.callback))
			} else {
				handle = handleWithBody(endpoint.callback)
			}

			mn.router.Handle(endpoint.method, endpoints.basePath+endpoint.path, handle)
		}
	}

	ctrlName := strings.Split(val.Type().String(), ".")
	mn.controllers[ctrlName[len(ctrlName)-1]] = controller
}

// CORS set CORS
func (mn *Minirest) CORS(opt CORSOption) {
	mn.router.HandleMethodNotAllowed = true
	mn.router.HandleOPTIONS = true
	mn.router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			if opt.AllowOrigin != "" {
				header.Add("Access-Control-Allow-Origin", opt.AllowOrigin)
			}

			if opt.AllowCredentials != "" {
				header.Add("Access-Control-Allow-Credentials", opt.AllowCredentials)
			}

			if opt.ExposeHeaders != "" {
				header.Add("Access-Control-Expose-Headers", opt.ExposeHeaders)
			}

			if opt.AllowHeaders != "" {
				header.Add("Access-Control-Allow-Headers", opt.AllowHeaders)
			}

			if opt.AllowMethods != "" {
				header.Add("Access-Control-Allow-Methods", opt.AllowMethods)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

// CORSOption set options for CORS headers
type CORSOption struct {
	AllowOrigin      string
	AllowCredentials string
	ExposeHeaders    string
	AllowHeaders     string
	AllowMethods     string
}
