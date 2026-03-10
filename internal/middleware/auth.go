package middleware

import (
	"net/http"

	conf "github.com/PrasitSahu/proxy/internal"
	"github.com/PrasitSahu/proxy/internal/api"
)

func VerifySignature(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		signature := req.Header.Get("Signature")

		if signature != conf.Config.Signature {
			http.Error(res, api.ErrAuth.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(res, req)
	})
}
