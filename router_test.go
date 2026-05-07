package thttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatSegment(t *testing.T) {
	require := require.New(t)

	testCases := []struct {
		desc string

		seg      Segment
		expected map[RouterType]string
	}{
		{
			desc: "static",

			seg: Segment{
				Name: "static",
				Raw:  "static",
				Type: Static,
			},
			expected: map[RouterType]string{
				RouterTypeStd:        "static",
				RouterTypeGorillaMux: "static",
				RouterTypeHttprouter: "static",
			},
		},

		{
			desc: "param",

			seg: Segment{
				Name: "param",
				Raw:  ":param",
				Type: Param,
			},
			expected: map[RouterType]string{
				RouterTypeStd:        "{param}",
				RouterTypeGorillaMux: "{param}",
				RouterTypeHttprouter: ":param",
			},
		},

		{
			desc: "catch-all",

			seg: Segment{
				Name: "param",
				Type: CatchAll,
			},
			expected: map[RouterType]string{
				RouterTypeStd:        "{param...}",
				RouterTypeGorillaMux: "{param:.*}",
				RouterTypeHttprouter: "*param",
			},
		},
	}

	for _, tc := range testCases {
		for routerType, expected := range tc.expected {
			t.Run(string(routerType)+": "+tc.desc, func(t *testing.T) {
				app := New(WithRouterType(routerType))
				router := app.router
				require.Equal(expected, router.FormatSegment(tc.seg))
			})
		}
	}
}

func TestPattern(t *testing.T) {
	require := require.New(t)

	testCases := []struct {
		desc string

		pattern  string
		path     string
		expected string
	}{
		{
			desc: "static",

			pattern:  "/static",
			path:     "/static",
			expected: "",
		},

		{
			desc: "param 1",

			pattern:  "/tasks/{key}",
			path:     "/tasks/abc",
			expected: "abc",
		},
		{
			desc: "param 2",

			pattern:  "/tasks/:key",
			path:     "/tasks/abc",
			expected: "abc",
		},

		{
			desc: "catch-all 1",

			pattern:  "/tasks/*key",
			path:     "/tasks/abc",
			expected: "abc",
		},

		{
			desc: "catch-all 1",

			pattern:  "/tasks/*key",
			path:     "/tasks/abc/123",
			expected: "abc/123",
		},
	}

	for _, tc := range testCases {
		for _, routerType := range allRouterTypes {
			t.Run(string(routerType)+": "+tc.desc, func(t *testing.T) {
				w := httptest.NewRecorder()

				app := New(WithRouterType(routerType))

				prefix := randomString()

				app.Get(tc.pattern, func(ctx Context) error {
					return ctx.String(http.StatusOK, prefix+"|"+ctx.PathParam("key"))
				})
				req := httptest.NewRequest(http.MethodGet, tc.path, nil)
				app.ServeHTTP(w, req)

				require.Equal(http.StatusOK, w.Code)
				require.Equal(prefix+"|"+tc.expected, w.Body.String())
			})
		}
	}
}
