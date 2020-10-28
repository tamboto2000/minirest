package minirest

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
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

type ResponseBuilder struct {
	statusCode int
	data       interface{}
}

func (resp *ResponseBuilder) Ok(data interface{}) *ResponseBuilder {
	resp.statusCode = CodeOk
	resp.data = Response{
		Code:    CodeOk,
		Message: MsgOk,
		Data:    data,
	}

	return resp
}

func (resp *ResponseBuilder) NoContent(desc string) *ResponseBuilder {
	resp.statusCode = CodeNoContent
	resp.data = Response{
		Code:        CodeNoContent,
		Message:     MsgNoContent,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) BadRequest(desc string) *ResponseBuilder {
	resp.statusCode = CodeBadRequest
	resp.data = Response{
		Code:        CodeBadRequest,
		Message:     MsgBadRequest,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) NotFound(desc string) *ResponseBuilder {
	resp.statusCode = CodeNotFound
	resp.data = Response{
		Code:        CodeNotFound,
		Message:     MsgNotFound,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) MethodNotAllowed(desc string) *ResponseBuilder {
	resp.statusCode = CodeMethodNotAllowed
	resp.data = Response{
		Code:        CodeMethodNotAllowed,
		Message:     MsgMethodNotAllowed,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) InternalError(desc string) *ResponseBuilder {
	resp.statusCode = CodeInternalError
	resp.data = Response{
		Code:        CodeInternalError,
		Message:     MsgInternalError,
		Description: desc,
	}

	return resp
}

func (resp *ResponseBuilder) ServerOverload(desc string) *ResponseBuilder {
	resp.statusCode = CodeOverload
	resp.data = Response{
		Code:        CodeOverload,
		Message:     MsgOverloadError,
		Description: desc,
	}

	return resp
}
