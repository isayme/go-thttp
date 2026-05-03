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
	allowedHttpMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
)
