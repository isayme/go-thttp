[![Go Doc](https://godoc.org/github.com/isayme/go-thttp?status.svg)](https://pkg.go.dev/github.com/isayme/go-thttp)


# Qucik Start
## Install
`go get github.com/isayme/go-thttp@latest`

## Usage
```go
package main

import (
	"log"

	"github.com/isayme/go-thttp"
)

func main() {
	app := thttp.New()

	// curl http://127.0.0.1:8080/version
	app.Get("/version", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]string{
			"version": "v1.0.0",
		})
	})

	// curl http://127.0.0.1:8080/tasks/123
	app.Get("/tasks/:tid", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]string{
			"tid": ctx.PathParam("tid"),
		})
	})

	// group route
	g := app.Group("/v1")

	// curl http://127.0.0.1:8080/v1/notes/123
	g.Get("/notes/{nid}", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]string{
			"nid": ctx.PathParam("nid"),
		})
	})

	// curl http://127.0.0.1:8080/static/index.html
	// curl http://127.0.0.1:8080/static/img/favicon.ico
	app.Get("/static/*path", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]string{
			"path": ctx.PathParam("path"),
		})
	})

	addr := ":8080"
	log.Printf("server start, listen on %s", addr)
	log.Fatal(app.Start(addr))
}
```

# Route Pattern
## Static
```
// match /users
app.Get("/users")
```

## Param
```
// match /users/abc
// match /users/123
app.Get("/users/:id")
app.Get("/users/{id}")
```

regex is not supported.

## Catch-All (wildcard)
```
// match /users/abc
// match /users/abc/def/ghi
app.Get("/users/*others")
app.Get("/users/{others...}")

// this is not supported
app.Get("/users/*")
```
