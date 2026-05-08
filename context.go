package thttp

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
)

var (
	_ Context = &thttpContext{}
)

const (
	defaultMemory = 32 << 20 // 32 MB
)

// Context represents the HTTP request context.
// It provides methods to access request data and write responses.
type Context interface {
	// Context returns the underlying context.Context.
	Context() context.Context

	// Request returns the current http.Request.
	Request() *http.Request
	// SetRequest sets the http.Request.
	SetRequest(r *http.Request)

	// Response returns the http.ResponseWriter.
	Response() http.ResponseWriter
	// SetResponse sets the http.ResponseWriter.
	SetResponse(r http.ResponseWriter)

	// Get retrieves a value from the context storage.
	Get(key interface{}) interface{}
	// Set stores a value in the context storage.
	Set(key, value interface{})

	// Method returns the HTTP method (GET, POST, etc.).
	Method() string
	// SetPathParam sets the path parameters.
	SetPathParam(fn PathParams)
	// PathParam returns the path parameter value by name.
	PathParam(name string) string
	// QueryParam returns the query parameter value by name.
	QueryParam(name string) string
	// QueryParams returns all query parameters as url.Values.
	QueryParams() url.Values
	// QueryString returns the raw query string.
	QueryString() string
	// FormParam returns the form parameter value by name.
	FormParam(name string) string
	// FormFile returns the uploaded file from the form.
	FormFile(name string) (*multipart.FileHeader, error)
	// MultipartForm returns the multipart form data.
	MultipartForm() (*multipart.Form, error)
	// Cookie returns the cookie by name.
	Cookie(name string) (*http.Cookie, error)
	// Cookies returns all cookies.
	Cookies() []*http.Cookie
	// SetCookie adds a cookie to the response.
	SetCookie(cookie *http.Cookie)
	// Header returns the header value by key.
	Header(key string) string
	// SetHeader sets a response header.
	SetHeader(key string, value string)
	// AddHeader adds a response header (allows multiple values).
	AddHeader(key string, value string)
	// DelHeader removes a response header.
	DelHeader(key string)

	// Blob writes binary data with the given content type.
	Blob(code int, contentType string, b []byte) error
	// JSON writes JSON data with the given status code.
	JSON(code int, i interface{}) error
	// String writes a plain text response.
	String(code int, s string) error
	// Stream streams data from an io.Reader.
	Stream(code int, contentType string, r io.Reader) error
	// Redirect redirects to the given URL with the given status code.
	Redirect(code int, url string) error

	// Logger returns the logger for this context.
	Logger() *slog.Logger
	// Reset resets the context for a new request.
	Reset(r *http.Request, w http.ResponseWriter, logger *slog.Logger)
}

type thttpContext struct {
	r *http.Request
	w http.ResponseWriter

	params PathParams

	query url.Values

	logger *slog.Logger

	lock  sync.RWMutex
	store map[interface{}]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) Context {
	return &thttpContext{
		r:     r,
		w:     w,
		store: make(map[interface{}]interface{}),
	}
}

func (ctx *thttpContext) Context() context.Context {
	return ctx.r.Context()
}

func (ctx *thttpContext) Request() *http.Request {
	return ctx.r
}

func (ctx *thttpContext) SetRequest(r *http.Request) {
	ctx.r = r
}

func (ctx *thttpContext) Response() http.ResponseWriter {
	return ctx.w
}

func (ctx *thttpContext) SetResponse(w http.ResponseWriter) {
	ctx.w = w
}

func (ctx *thttpContext) Get(key interface{}) interface{} {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	return ctx.store[key]
}

func (ctx *thttpContext) Set(key, value interface{}) {
	ctx.lock.Lock()
	defer ctx.lock.Unlock()

	ctx.store[key] = value
}

func (ctx *thttpContext) Method() string {
	return ctx.r.Method
}

func (ctx *thttpContext) SetPathParam(fn PathParams) {
	ctx.Set(PathParamsCtxKey, fn)
	ctx.params = fn
}

func (ctx *thttpContext) PathParam(name string) string {
	if ctx.params == nil {
		return ""
	}
	return ctx.params.Get(name)
}

func (ctx *thttpContext) QueryParam(name string) string {
	return ctx.QueryParams().Get(name)
}

func (ctx *thttpContext) QueryParams() url.Values {
	if ctx.query == nil {
		ctx.query = ctx.r.URL.Query()
	}
	return ctx.query
}

func (ctx *thttpContext) QueryString() string {
	return ctx.r.URL.RawQuery
}

func (ctx *thttpContext) FormParam(name string) string {
	return ctx.r.FormValue(name)
}

func (ctx *thttpContext) FormFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := ctx.r.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

func (ctx *thttpContext) MultipartForm() (*multipart.Form, error) {
	err := ctx.r.ParseMultipartForm(defaultMemory)
	return ctx.r.MultipartForm, err
}

func (ctx *thttpContext) Cookie(name string) (*http.Cookie, error) {
	return ctx.r.Cookie(name)
}

func (ctx *thttpContext) Cookies() []*http.Cookie {
	return ctx.r.Cookies()
}

func (ctx *thttpContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.Response(), cookie)
}

func (ctx *thttpContext) Header(key string) string {
	return ctx.r.Header.Get(key)
}

func (ctx *thttpContext) SetHeader(key string, value string) {
	ctx.w.Header().Set(key, value)
}

func (ctx *thttpContext) AddHeader(key string, value string) {
	ctx.w.Header().Add(key, value)
}

func (ctx *thttpContext) DelHeader(key string) {
	ctx.w.Header().Del(key)
}

func (ctx *thttpContext) Blob(code int, contentType string, b []byte) (err error) {
	ctx.SetHeader(HeaderContentType, contentType)
	ctx.w.WriteHeader(code)
	_, err = ctx.w.Write(b)
	return
}

func (ctx *thttpContext) JSON(code int, body interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx.w.WriteHeader(code)
	ctx.w.Write(data)

	return nil
}

func (ctx *thttpContext) String(statusCode int, body string) error {
	ctx.w.WriteHeader(statusCode)
	ctx.w.Write([]byte(body))
	return nil
}

func (ctx *thttpContext) Stream(code int, contentType string, r io.Reader) (err error) {
	ctx.SetHeader(HeaderContentType, contentType)
	ctx.w.WriteHeader(code)
	_, err = io.Copy(ctx.w, r)
	return
}

func (ctx *thttpContext) Redirect(code int, url string) error {
	ctx.w.Header().Set(HeaderLocation, url)
	ctx.w.WriteHeader(code)
	return nil
}

func (ctx *thttpContext) Logger() *slog.Logger {
	return ctx.logger
}

func (ctx *thttpContext) Reset(r *http.Request, w http.ResponseWriter, logger *slog.Logger) {
	ctx.r = r
	ctx.w = w
	ctx.query = nil
	ctx.params = nil
	ctx.logger = logger
	ctx.store = make(map[interface{}]interface{})
}

func MustGetContextFromRequest(r *http.Request) Context {
	ctx := r.Context().Value(ContextKey)
	if ctx == nil {
		panic("thttp: no context found in request")
	}
	return ctx.(Context)
}

func SetHandlerInCtx(ctx Context, h HandlerFunc) {
	ctx.Set(HandlerKey, h)
}

func MustGetHandlerFromCtx(ctx Context) HandlerFunc {
	h := ctx.Get(HandlerKey)
	if h == nil {
		panic("thttp: no handler found in context")
	}
	return h.(HandlerFunc)
}
