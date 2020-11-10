package minirest

import (
	"encoding/json"
	"log"
	"net/http"
)

// HTTP status codes
const (
	CodeOk               = 200
	CodeNoContent        = 204
	CodeBadRequest       = 400
	CodeNotFound         = 404
	CodeMethodNotAllowed = 405
	CodeTooManyRequest   = 429
	CodeInternalError    = 500
	CodeOverload         = 503
)

// HTTP status message
const (
	MsgOk               = "ok"
	MsgNoContent        = "no_content"
	MsgBadRequest       = "bad_request"
	MsgNotFound         = "not_found"
	MsgMethodNotAllowed = "method_not_allowed"
	MsgTooManyRequest   = "too_many_request"
	MsgInternalError    = "internal_error"
	MsgOverloadError    = "server_overload"
)

// Response is body for HTTP response
type Response struct {
	StatusCode  int         `json:"statusCode"`
	Status      string      `json:"status"`
	Description string      `json:"description,omitempty"`
	Body        interface{} `json:"body,omitempty"`
}

// ResponseBuilder is a response builder
type ResponseBuilder struct {
	// Set to true for returning gzip encoded response
	Gzip       bool
	statusCode int
	headers    [][2]string
	body       interface{}
}

// Status set status code
func (resp *ResponseBuilder) Status(code int) *ResponseBuilder {
	resp.statusCode = code
	return resp
}

// Headers add headers
func (resp *ResponseBuilder) Headers(headers [][2]string) *ResponseBuilder {
	resp.headers = headers

	return resp
}

// Body set body
func (resp *ResponseBuilder) Body(body interface{}) *ResponseBuilder {
	resp.body = body

	return resp
}

// Ok build response with HTTP Status 200
func (resp *ResponseBuilder) Ok(data interface{}) *ResponseBuilder {
	resp.statusCode = CodeOk
	resp.body = Response{
		StatusCode: CodeOk,
		Status:     MsgOk,
		Body:       data,
	}

	return resp
}

// NoContent build response with HTTP Status 204
func (resp *ResponseBuilder) NoContent(desc string) *ResponseBuilder {
	resp.statusCode = CodeNoContent
	resp.body = Response{
		StatusCode:  CodeNoContent,
		Status:      MsgNoContent,
		Description: desc,
	}

	return resp
}

// BadRequest build response with HTTP Status 400
func (resp *ResponseBuilder) BadRequest(desc string) *ResponseBuilder {
	resp.statusCode = CodeBadRequest
	resp.body = Response{
		StatusCode:  CodeBadRequest,
		Status:      MsgBadRequest,
		Description: desc,
	}

	return resp
}

// NotFound build response with HTTP Status 404
func (resp *ResponseBuilder) NotFound(desc string) *ResponseBuilder {
	resp.statusCode = CodeNotFound
	resp.body = Response{
		StatusCode:  CodeNotFound,
		Status:      MsgNotFound,
		Description: desc,
	}

	return resp
}

// MethodNotAllowed build response with HTTP Status 405
func (resp *ResponseBuilder) MethodNotAllowed(desc string) *ResponseBuilder {
	resp.statusCode = CodeMethodNotAllowed
	resp.body = Response{
		StatusCode:  CodeMethodNotAllowed,
		Status:      MsgMethodNotAllowed,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) TooManyRequest(desc string) *ResponseBuilder {
	resp.statusCode = CodeTooManyRequest
	resp.body = Response{
		StatusCode:  CodeTooManyRequest,
		Status:      MsgTooManyRequest,
		Description: desc,
	}

	return resp
}

// InternalError build response with HTTP Status 500
func (resp *ResponseBuilder) InternalError(desc string) *ResponseBuilder {
	resp.statusCode = CodeInternalError
	resp.body = Response{
		StatusCode:  CodeInternalError,
		Status:      MsgInternalError,
		Description: desc,
	}

	return resp
}

// ServerOverload build response with HTTP Status 503
func (resp *ResponseBuilder) ServerOverload(desc string) *ResponseBuilder {
	resp.statusCode = CodeOverload
	resp.body = Response{
		StatusCode:  CodeOverload,
		Status:      MsgOverloadError,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) write(w http.ResponseWriter) {
	for _, header := range resp.headers {
		w.Header().Add(header[0], header[1])
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Gzip {
		data, err := json.Marshal(resp.body)
		if err != nil {
			log.Println(err.Error())
			return
		}

		writeGzipResp(w, data, resp.statusCode)
		return
	}

	w.WriteHeader(resp.statusCode)
	if err := json.NewEncoder(w).Encode(resp.body); err != nil {
		log.Println(err.Error())
	}
}
