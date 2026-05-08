package thttp

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpMethod(t *testing.T) {
	require := require.New(t)
	for _, routerType := range allRouterTypes {

		t.Run(string(routerType)+": Any", func(t *testing.T) {
			for _, method := range allowedHttpMethods {
				t.Run(method, func(t *testing.T) {
					req := httptest.NewRequest(method, "/method", nil)

					w := httptest.NewRecorder()

					app := New(WithRouterType(routerType))
					app.Any("/method", func(ctx Context) error {
						return ctx.String(http.StatusOK, ctx.Method())
					})

					app.ServeHTTP(w, req)

					require.Equal(http.StatusOK, w.Code)
					require.Equal(method, w.Body.String())
				})
			}
		})

		t.Run(string(routerType)+": Get", func(t *testing.T) {
			method := http.MethodGet

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Get("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})

		t.Run(string(routerType)+": Post", func(t *testing.T) {
			method := http.MethodPost

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Post("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})

		t.Run(string(routerType)+": Put", func(t *testing.T) {
			method := http.MethodPut

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Put("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})

		t.Run(string(routerType)+": Patch", func(t *testing.T) {
			method := http.MethodPatch

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Patch("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})

		t.Run(string(routerType)+": Delete", func(t *testing.T) {
			method := http.MethodDelete

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Delete("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})

		t.Run(string(routerType)+": Options", func(t *testing.T) {
			method := http.MethodOptions

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Options("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})

		t.Run(string(routerType)+": Head", func(t *testing.T) {
			method := http.MethodHead

			req := httptest.NewRequest(method, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))
			app.Head("/method", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.Method())
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(method, w.Body.String())
		})
	}
}

func TestStatusCode(t *testing.T) {
	require := require.New(t)

	for _, routerType := range allRouterTypes {
		for _, code := range []int{200, 201, 301, 302, 400, 401, 403, 404, 500, 501, 502} {
			t.Run(string(routerType)+": 200", func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/method", nil)

				w := httptest.NewRecorder()

				app := New(WithRouterType(routerType))

				app.Get("/method", func(ctx Context) error {
					return ctx.String(code, fmt.Sprintf("%d", code))
				})

				app.ServeHTTP(w, req)

				require.Equal(code, w.Code)
				require.Equal(fmt.Sprintf("%d", code), w.Body.String())
			})
		}
	}
}

func TestNotFound(t *testing.T) {
	require := require.New(t)

	for _, routerType := range allRouterTypes {
		t.Run(string(routerType)+": not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))

			app.ServeHTTP(w, req)

			require.Equal(http.StatusNotFound, w.Code)
			require.Equal("404 page not found", w.Body.String())
		})

		t.Run(string(routerType)+": custom not found", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/method", nil)

			w := httptest.NewRecorder()

			msg := randomString()

			app := New(WithRouterType(routerType), WithNotFoundHandler(func(ctx Context) error {
				return ctx.String(http.StatusNotImplemented, msg)
			}))

			app.ServeHTTP(w, req)

			require.Equal(http.StatusNotImplemented, w.Code)
			require.Equal(msg, w.Body.String())
		})
	}
}

func TestErrorHandler(t *testing.T) {
	require := require.New(t)

	for _, routerType := range allRouterTypes {
		t.Run(string(routerType)+": error handler", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/method", nil)

			w := httptest.NewRecorder()

			app := New(WithRouterType(routerType))

			errMsg := randomString()
			app.Get("/method", func(ctx Context) error {
				return errors.New(errMsg)
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusInternalServerError, w.Code)
			require.Equal(errMsg, w.Body.String())
		})

		t.Run(string(routerType)+": custom error handler", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/method", nil)

			w := httptest.NewRecorder()

			errMsgPrefix := randomString()
			app := New(WithRouterType(routerType), WithErrorHandler(func(ctx Context, err error) error {
				return ctx.String(http.StatusBadGateway, errMsgPrefix+err.Error())
			}))

			errMsg := randomString()
			app.Get("/method", func(ctx Context) error {
				return errors.New(errMsg)
			})

			app.ServeHTTP(w, req)

			require.Equal(http.StatusBadGateway, w.Code)
			require.Equal(errMsgPrefix+errMsg, w.Body.String())
		})
	}
}

func TestContextPathParam(t *testing.T) {
	require := require.New(t)

	t.Run("PathParam", func(t *testing.T) {
		key := randomString()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%s", key), nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/user/{id}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("id"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(key, w.Body.String())
	})
}

func TestContextQueryParam(t *testing.T) {
	require := require.New(t)

	t.Run("QueryParam", func(t *testing.T) {
		key := randomString()
		value := randomString()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/search?%s=%s", key, value), nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/search", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.QueryParam(key))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(value, w.Body.String())
	})
}

func TestContextQueryParams(t *testing.T) {
	require := require.New(t)

	t.Run("QueryParams", func(t *testing.T) {
		key1 := randomString()
		value1 := randomString()
		key2 := randomString()
		value2 := randomString()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/search?%s=%s&%s=%s", key1, value1, key2, value2), nil)
		w := httptest.NewRecorder()

		app := New()

		flag := false
		app.Get("/search", func(ctx Context) error {
			params := ctx.QueryParams()
			require.Equal(value1, params.Get(key1))
			require.Equal(value2, params.Get(key2))
			require.Equal(2, len(params))
			flag = true
			return ctx.String(http.StatusOK, "ok")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.True(flag)
	})
}

func TestContextQueryString(t *testing.T) {
	require := require.New(t)

	t.Run("QueryString", func(t *testing.T) {
		key1 := randomString()
		value1 := randomString()
		key2 := randomString()
		value2 := randomString()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/search?%s=%s&%s=%s", key1, value1, key2, value2), nil)
		w := httptest.NewRecorder()

		app := New()
		flag := false
		app.Get("/search", func(ctx Context) error {
			qs := ctx.QueryString()
			require.Equal(qs, fmt.Sprintf("%s=%s&%s=%s", key1, value1, key2, value2))
			flag = true
			return ctx.String(http.StatusOK, "ok")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.True(flag)
	})
}

func TestContextFormParam(t *testing.T) {
	require := require.New(t)

	t.Run("FormParam", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader("username=john&age=30"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()

		app := New()
		app.Post("/submit", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.FormParam("username")+","+ctx.FormParam("age"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("john,30", w.Body.String())
	})
}

func TestContextFormFile(t *testing.T) {
	require := require.New(t)

	t.Run("FormFile", func(t *testing.T) {
		var b bytes.Buffer
		wr := multipart.NewWriter(&b)
		writer, _ := wr.CreateFormFile("file", "test.txt")
		writer.Write([]byte("test content"))
		wr.Close()

		req := httptest.NewRequest(http.MethodPost, "/upload", &b)
		req.Header.Set("Content-Type", wr.FormDataContentType())

		w := httptest.NewRecorder()

		app := New()
		app.Post("/upload", func(ctx Context) error {
			fh, err := ctx.FormFile("file")
			if err != nil {
				return err
			}

			return ctx.String(http.StatusOK, fh.Filename)
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("test.txt", w.Body.String())
	})
}

func TestContextMultipartForm(t *testing.T) {
	require := require.New(t)

	t.Run("MultipartForm", func(t *testing.T) {
		var b bytes.Buffer
		wr := multipart.NewWriter(&b)
		wr.WriteField("username", "john")
		writer, _ := wr.CreateFormFile("file", "test.txt")
		writer.Write([]byte("test content"))
		wr.Close()

		req := httptest.NewRequest(http.MethodPost, "/submit", &b)
		req.Header.Set("Content-Type", wr.FormDataContentType())

		w := httptest.NewRecorder()

		app := New()
		app.Post("/submit", func(ctx Context) error {
			form, err := ctx.MultipartForm()
			if err != nil {
				return err
			}
			return ctx.String(http.StatusOK, form.Value["username"][0])
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("john", w.Body.String())
	})
}

func TestContextCookie(t *testing.T) {
	require := require.New(t)

	t.Run("Cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/cookie", nil)
		req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})

		w := httptest.NewRecorder()

		app := New()
		app.Get("/cookie", func(ctx Context) error {
			cookie, err := ctx.Cookie("session")
			if err != nil {
				return err
			}
			return ctx.String(http.StatusOK, cookie.Value)
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("abc123", w.Body.String())
	})
}

func TestContextCookies(t *testing.T) {
	require := require.New(t)

	t.Run("Cookies", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/cookies", nil)
		req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
		req.AddCookie(&http.Cookie{Name: "theme", Value: "dark"})

		w := httptest.NewRecorder()

		app := New()
		app.Get("/cookies", func(ctx Context) error {
			cookies := ctx.Cookies()
			return ctx.String(http.StatusOK, cookies[0].Name+":"+cookies[0].Value+","+cookies[1].Name+":"+cookies[1].Value)
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("session:abc123,theme:dark", w.Body.String())
	})
}

func TestContextSetCookie(t *testing.T) {
	require := require.New(t)

	t.Run("SetCookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/setcookie", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/setcookie", func(ctx Context) error {
			ctx.SetCookie(&http.Cookie{Name: "token", Value: "xyz789", Path: "/"})
			return ctx.String(http.StatusOK, "ok")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Contains(w.Result().Cookies()[0].Name, "token")
		require.Equal("xyz789", w.Result().Cookies()[0].Value)
	})
}

func TestContextHeader(t *testing.T) {
	require := require.New(t)

	t.Run("Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/header", nil)
		req.Header.Set("X-Custom-Header", "customvalue")

		w := httptest.NewRecorder()

		app := New()
		app.Get("/header", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Header("X-Custom-Header"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("customvalue", w.Body.String())
	})
}

func TestContextSetHeader(t *testing.T) {
	require := require.New(t)

	t.Run("SetHeader", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/setheader", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/setheader", func(ctx Context) error {
			ctx.SetHeader("X-Response-Header", "responsevalue")
			return ctx.String(http.StatusOK, "ok")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("responsevalue", w.Header().Get("X-Response-Header"))
	})
}

func TestContextAddHeader(t *testing.T) {
	require := require.New(t)

	t.Run("AddHeader", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/addheader", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/addheader", func(ctx Context) error {
			ctx.AddHeader("X-Custom-Header", "value1")
			ctx.AddHeader("X-Custom-Header", "value2")
			return ctx.String(http.StatusOK, "ok")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal([]string{"value1", "value2"}, w.Header().Values("X-Custom-Header"))
	})
}

func TestContextDelHeader(t *testing.T) {
	require := require.New(t)

	t.Run("DelHeader", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/delheader", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/delheader", func(ctx Context) error {
			ctx.SetHeader("X-Custom-Header", "value")
			ctx.DelHeader("X-Custom-Header")
			return ctx.String(http.StatusOK, w.Header().Get("X-Custom-Header"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("", w.Body.String())
	})
}

func TestContextBlob(t *testing.T) {
	require := require.New(t)

	t.Run("Blob", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/blob", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/blob", func(ctx Context) error {
			return ctx.Blob(http.StatusOK, "application/octet-stream", []byte("binary data"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("application/octet-stream", w.Header().Get("Content-Type"))
		require.Equal("binary data", w.Body.String())
	})
}

func TestContextJSON(t *testing.T) {
	require := require.New(t)

	t.Run("JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/json", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/json", func(ctx Context) error {
			return ctx.JSON(http.StatusOK, map[string]string{"name": "john", "age": "30"})
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Contains(w.Body.String(), `"name":"john"`)
		require.Contains(w.Body.String(), `"age":"30"`)
	})
}

func TestContextString(t *testing.T) {
	require := require.New(t)

	t.Run("String", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/string", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/string", func(ctx Context) error {
			return ctx.String(http.StatusOK, "plain text response")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("plain text response", w.Body.String())
	})
}

func TestContextStream(t *testing.T) {
	require := require.New(t)

	t.Run("Stream", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/stream", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/stream", func(ctx Context) error {
			return ctx.Stream(http.StatusOK, "text/plain", strings.NewReader("streamed data"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("text/plain", w.Header().Get("Content-Type"))
		require.Equal("streamed data", w.Body.String())
	})
}

func TestContextRedirect(t *testing.T) {
	require := require.New(t)

	t.Run("Redirect", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/redirect", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/redirect", func(ctx Context) error {
			return ctx.Redirect(http.StatusFound, "/newlocation")
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusFound, w.Code)
		require.Equal("/newlocation", w.Header().Get("Location"))
	})
}

func TestPool(t *testing.T) {
	require := require.New(t)

	t.Run("reuse ctx", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/hi", nil)

		app := New()

		expected := ""

		app.Get("/hi", func(ctx Context) error {
			expected = randomString()
			return ctx.String(http.StatusOK, expected)
		})

		for i := 0; i < 100; i++ {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
			require.Equal(http.StatusOK, w.Code)
			require.Equal(expected, w.Body.String())
		}
	})
}
