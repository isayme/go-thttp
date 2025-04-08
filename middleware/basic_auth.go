package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/isayme/go-thttp"
)

func BasicAuth(realm string, creds map[string]string) thttp.MiddlewareFunc {
	return func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			r := ctx.Request()
			w := ctx.Response()

			user, pass, ok := r.BasicAuth()
			if !ok {
				basicAuthFailed(w, realm)
				return nil
			}

			credPass, credUserOk := creds[user]
			if !credUserOk || subtle.ConstantTimeCompare([]byte(pass), []byte(credPass)) != 1 {
				basicAuthFailed(w, realm)
				return nil
			}

			return next(ctx)
		}
	}
}

func basicAuthFailed(w http.ResponseWriter, realm string) {
	w.Header().Add("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, realm))
	w.WriteHeader(http.StatusUnauthorized)
}
