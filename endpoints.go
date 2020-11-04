package minirest

type endpoint struct {
	method   string
	path     string
	callback interface{}
}

//Endpoints register handlers its path and method
type Endpoints struct {
	basePath   string
	endpoints  []endpoint
	middleware *handleChain
}

//BasePath set base path for endpoints
func (ep *Endpoints) BasePath(path string) {
	ep.basePath = path
}

//Add add endpoint with custom method
func (ep *Endpoints) Add(method, path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{method, path, callback})
}

//GET add endpoint with method GET
func (ep *Endpoints) GET(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"GET", path, callback})
}

//DELETE add endpoint with method DELETE
func (ep *Endpoints) DELETE(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"DELETE", path, callback})
}

//POST add method endpoint with method POST
func (ep *Endpoints) POST(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"POST", path, callback})
}

//PUT add method endpoint with method PUT
func (ep *Endpoints) PUT(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"PUT", path, callback})
}

//PATCH add method endpoint with method PATCH
func (ep *Endpoints) PATCH(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"PATCH", path, callback})
}

//Middlewares resgiter middleware chain
func (ep *Endpoints) Middlewares(mds ...handleToHandle) {
	if ep.middleware == nil {
		ep.middleware = new(handleChain)
	}

	ep.middleware.handles = append(ep.middleware.handles, mds...)
}
