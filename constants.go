package thttp

import "net/http"

const (
	HeaderContentType        = "Content-Type"
	HeaderXRequestID         = "X-Request-ID"
	HeaderLocation           = "Location"
	HeaderContentDisposition = "Content-Disposition"
)

const (
	MIMEApplicationJSON = "application/json"
)

var (
	methods = [...]string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
)
