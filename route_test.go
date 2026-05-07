package thttp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePath(t *testing.T) {
	require := require.New(t)

	t.Run("blank", func(t *testing.T) {
		segs := ParsePath("  ")
		require.Equal(0, len(segs))
	})

	t.Run("empty", func(t *testing.T) {
		segs := ParsePath("")
		require.Equal(0, len(segs))
	})

	t.Run("root", func(t *testing.T) {
		segs := ParsePath("/")
		require.Equal(0, len(segs))
	})

	t.Run("multi parts: /abc/123", func(t *testing.T) {
		segs := ParsePath("/abc/123")
		require.Equal(2, len(segs))
		require.Equal(Static, segs[0].Type)
		require.Equal("abc", segs[0].Raw)
		require.Equal(Static, segs[1].Type)
		require.Equal("123", segs[1].Raw)
	})

	t.Run("double slash: /abc//123", func(t *testing.T) {
		segs := ParsePath("/abc//123")
		require.Equal(2, len(segs))
		require.Equal(Static, segs[0].Type)
		require.Equal("abc", segs[0].Raw)
		require.Equal(Static, segs[1].Type)
		require.Equal("123", segs[1].Raw)
	})

	testCases := []struct {
		desc string

		typ  SegmentType
		raw  string
		name string
	}{
		{
			desc: "static: /abc",

			raw:  "abc",
			typ:  Static,
			name: "abc",
		},
		{
			desc: "brace path param: /{abc}",

			raw:  "{abc}",
			typ:  Param,
			name: "abc",
		},
		{
			desc: "colon path param: /:abc",

			raw:  ":abc",
			typ:  Param,
			name: "abc",
		},

		// {
		// 	desc: "catch-all: *",

		// 	raw:  "*",
		// 	typ:  CatchAll,
		// 	name: "*",
		// },

		{
			desc: "catch-all: *abc",

			raw:  "*abc",
			typ:  CatchAll,
			name: "abc",
		},

		{
			desc: "catch-all: {abc...}",

			raw:  "{abc...}",
			typ:  CatchAll,
			name: "abc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			segs := ParsePath(tc.raw)
			require.Equal(1, len(segs))
			require.Equal(tc.typ, segs[0].Type)
			require.Equal(tc.raw, segs[0].Raw)
			require.Equal(tc.name, segs[0].Name)
		})
	}
}
