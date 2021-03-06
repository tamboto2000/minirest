package minirest

type endpoint struct {
	method   string
	path     string
	callback interface{}
}

// Endpoints register handlers its path and method
type Endpoints struct {
	// Set to true for returning gzip encoded response on all endpoints
	Gzip       bool
	basePath   string
	endpoints  []endpoint
	middleware *handleChain
}

// BasePath set base path for endpoints
func (ep *Endpoints) BasePath(path string) {
	ep.basePath = path
}

// Add add endpoint with custom method
func (ep *Endpoints) Add(method, path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{method, path, callback})
}

// GET add endpoint with method GET
func (ep *Endpoints) GET(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"GET", path, callback})
}

// DELETE add endpoint with method DELETE
func (ep *Endpoints) DELETE(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"DELETE", path, callback})
}

// POST add method endpoint with method POST
func (ep *Endpoints) POST(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"POST", path, callback})
}

// PUT add method endpoint with method PUT
func (ep *Endpoints) PUT(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"PUT", path, callback})
}

// PATCH add method endpoint with method PATCH
func (ep *Endpoints) PATCH(path string, callback interface{}) {
	ep.endpoints = append(ep.endpoints, endpoint{"PATCH", path, callback})
}

// Middlewares register middleware chain.
// miniREST is using julienschmidt/httprouter for implementing router,
// so the middleware will use httprouter.Handle as its handle
func (ep *Endpoints) Middlewares(mds ...handleToHandle) {
	if ep.middleware == nil {
		ep.middleware = new(handleChain)
	}

	ep.middleware.handles = append(ep.middleware.handles, mds...)
}
