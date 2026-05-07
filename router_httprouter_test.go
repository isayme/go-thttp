package thttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHttprouterDuplicaPattern conflicts pattern not allowed by httprouter
func TestHttprouterDuplicaPattern(t *testing.T) {
	assert := assert.New(t)

	t.Run("allow duplicate route 1", func(t *testing.T) {
		app := New(WithRouterType(RouterTypeHttprouter))

		app.Get("/task/{tid}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("tid"))
		})

		assert.Panics(func() {
			app.Get("/task/VERSION", func(ctx Context) error {
				return ctx.String(http.StatusOK, "-VERSION-")
			})
		})
	})

	t.Run("allow duplicate route 2", func(t *testing.T) {
		app := New(WithRouterType(RouterTypeHttprouter))

		app.Get("/task/VERSION", func(ctx Context) error {
			return ctx.String(http.StatusOK, "-VERSION-")
		})

		assert.Panics(func() {
			app.Get("/task/{tid}", func(ctx Context) error {
				return ctx.String(http.StatusOK, ctx.PathParam("tid"))
			})
		})
	})
}
