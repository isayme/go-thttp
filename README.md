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

	app.Get("/version", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]string{
			"version": "v1.0.0",
		})
	})

	app.Get("/tasks/:tid", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]string{
			"tid": ctx.PathParam("tid"),
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