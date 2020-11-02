package minirest

import (
	"encoding/json"
	"log"
	"net/http"
)

//HTTP status codes
const (
	CodeOk               = 200
	CodeNoContent        = 204
	CodeBadRequest       = 400
	CodeNotFound         = 404
	CodeMethodNotAllowed = 405
	CodeInternalError    = 500
	CodeOverload         = 503
)

//HTTP status message
const (
	MsgOk               = "ok"
	MsgNoContent        = "no_content"
	MsgBadRequest       = "bad_request"
	MsgNotFound         = "not_found"
	MsgMethodNotAllowed = "method_not_allowed"
	MsgInternalError    = "internal_error"
	MsgOverloadError    = "server_overload"
)

//Response is body for HTTP response
type Response struct {
	StatusCode  int         `json:"statusCode"`
	Status      string      `json:"status"`
	Description string      `json:"description,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

//ResponseBuilder is a response builder
type ResponseBuilder struct {
	statusCode int
	data       interface{}
}

//Ok build response with HTTP Status 200
func (resp *ResponseBuilder) Ok(data interface{}) *ResponseBuilder {
	resp.statusCode = CodeOk
	resp.data = Response{
		StatusCode: CodeOk,
		Status:     MsgOk,
		Data:       data,
	}

	return resp
}

//NoContent build response with HTTP Status 204
func (resp *ResponseBuilder) NoContent(desc string) *ResponseBuilder {
	resp.statusCode = CodeNoContent
	resp.data = Response{
		StatusCode:  CodeNoContent,
		Status:      MsgNoContent,
		Description: desc,
	}

	return resp
}

//BadRequest build response with HTTP Status 400
func (resp *ResponseBuilder) BadRequest(desc string) *ResponseBuilder {
	resp.statusCode = CodeBadRequest
	resp.data = Response{
		StatusCode:  CodeBadRequest,
		Status:      MsgBadRequest,
		Description: desc,
	}

	return resp
}

//NotFound build response with HTTP Status 404
func (resp *ResponseBuilder) NotFound(desc string) *ResponseBuilder {
	resp.statusCode = CodeNotFound
	resp.data = Response{
		StatusCode:  CodeNotFound,
		Status:      MsgNotFound,
		Description: desc,
	}

	return resp
}

//MethodNotAllowed build response with HTTP Status 405
func (resp *ResponseBuilder) MethodNotAllowed(desc string) *ResponseBuilder {
	resp.statusCode = CodeMethodNotAllowed
	resp.data = Response{
		StatusCode:  CodeMethodNotAllowed,
		Status:      MsgMethodNotAllowed,
		Description: desc,
	}

	return resp
}

//InternalError build response with HTTP Status 500
func (resp *ResponseBuilder) InternalError(desc string) *ResponseBuilder {
	resp.statusCode = CodeInternalError
	resp.data = Response{
		StatusCode:  CodeInternalError,
		Status:      MsgInternalError,
		Description: desc,
	}

	return resp
}

//ServerOverload build response with HTTP Status 503
func (resp *ResponseBuilder) ServerOverload(desc string) *ResponseBuilder {
	resp.statusCode = CodeOverload
	resp.data = Response{
		StatusCode:  CodeOverload,
		Status:      MsgOverloadError,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.statusCode)
	if err := json.NewEncoder(w).Encode(resp.data); err != nil {
		log.Println(err.Error())
	}
}
