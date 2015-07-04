package middleware

import (
	"net/http"
	"strings"
	//"fmt"

	"github.com/zenazn/goji/web"
	"github.com/ckpt/backend-services/players"
)

func TokenHandler(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/login") {
			h.ServeHTTP(w,r)
			return
		}
		authzHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authzHeader, "CKPT ")
		p, err := players.PlayerByUserToken(token)
		if err != nil {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}
		c.Env["authPlayer"] = p.UUID
		c.Env["authUser"] = p.User.Username
		c.Env["authIsAdmin"] = p.User.Admin
		h.ServeHTTP(w,r)
	}
	return http.HandlerFunc(fn)
}
