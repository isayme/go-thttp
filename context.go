package thttp

import (
	"context"
	"encoding/json"
	"io"
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

type Context interface {
	Context() context.Context

	Request() *http.Request
	SetRequest(r *http.Request)

	Response() http.ResponseWriter
	SetResponse(r http.ResponseWriter)

	Get(key interface{}) interface{}
	Set(key, value interface{})

	Method() string
	SetPathParam(fn PathParams)
	PathParam(name string) string
	QueryParam(name string) string
	QueryParams() url.Values
	QueryString() string
	FormParam(name string) string
	FormFile(name string) (*multipart.FileHeader, error)
	MultipartForm() (*multipart.Form, error)
	Cookie(name string) (*http.Cookie, error)
	Cookies() []*http.Cookie
	SetCookie(cookie *http.Cookie)
	Header(key string) string
	SetHeader(key string, value string)
	AddHeader(key string, value string)
	DelHeader(key string)

	Blob(code int, contentType string, b []byte) error
	JSON(code int, i interface{}) error
	String(code int, s string) error
	Stream(code int, contentType string, r io.Reader) error
	Redirect(code int, url string) error

	Reset(r *http.Request, w http.ResponseWriter)
}

type thttpContext struct {
	r *http.Request
	w http.ResponseWriter

	params PathParams

	query url.Values

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

func (ctx *thttpContext) Reset(r *http.Request, w http.ResponseWriter) {
	ctx.r = r
	ctx.w = w
	ctx.query = nil
	ctx.params = nil
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
