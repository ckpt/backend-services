package middleware

import (
	"net/http"
	"strings"
	//"fmt"

	"github.com/ckpt/backend-services/players"
	"github.com/zenazn/goji/web"
)

func TokenHandler(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			h.ServeHTTP(w, r)
			return
		}
		authzHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authzHeader, "CKPT ")
		if token == authzHeader || len(token) < 6 {
			w.WriteHeader(403)
			w.Write([]byte("Invalid auth header or token"))
			return
		}
		p, err := players.PlayerByUserToken(token)
		if err != nil {
			w.WriteHeader(403)
			w.Write([]byte("Unauthorized"))
			return
		}
		c.Env["authPlayer"] = p.UUID
		c.Env["authUser"] = p.User.Username
		c.Env["authIsAdmin"] = p.User.Admin
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
