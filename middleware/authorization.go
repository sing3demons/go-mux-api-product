package middleware

import (
	"github/sing3demons/go_mux_api/config"
	"github/sing3demons/go_mux_api/models"
	"net/http"

	"github.com/casbin/casbin"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("id")
		var user models.User
		config.GetDB().First(&user, id)

		enforcer := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		ok := enforcer.Enforce(user, r.URL.Path, r.Method)
		if !ok {
			JSON(w, http.StatusForbidden)(Map{"error": "ou are not allowed to access this resource"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
