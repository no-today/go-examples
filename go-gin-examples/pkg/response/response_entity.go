package response

import (
	"net/http"
)

type ResponseEntity struct {
	Status  int
	Headers map[string]string
	Body    interface{}
	Msg     string
}

type responseEntityBuilder struct {
	status  int
	headers map[string]string
	body    interface{}
	msg     string
}

func ResponseEntityBuilder() *responseEntityBuilder {
	return &responseEntityBuilder{}
}

func (r *responseEntityBuilder) Headers(headers map[string]string) *responseEntityBuilder {
	r.headers = headers
	return r
}

func (r *responseEntityBuilder) Header(header, value string) *responseEntityBuilder {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[header] = value
	return r
}

func (r *responseEntityBuilder) Body(body interface{}) *responseEntityBuilder {
	r.body = body
	return r
}

func (r *responseEntityBuilder) Build() *ResponseEntity {
	return &ResponseEntity{
		Status:  r.status,
		Headers: r.headers,
		Body:    r.body,
		Msg:     r.msg,
	}
}

func (r *responseEntityBuilder) Ok(msg string) *responseEntityBuilder {
	r.status = http.StatusOK
	r.msg = msg
	return r
}

func (r *responseEntityBuilder) BadRequest() *responseEntityBuilder {
	r.status = http.StatusBadRequest
	return r
}

func (r *responseEntityBuilder) NotFound() *responseEntityBuilder {
	r.status = http.StatusNotFound
	return r
}

func (r *responseEntityBuilder) InternalServerError() *responseEntityBuilder {
	r.status = http.StatusInternalServerError
	return r
}
