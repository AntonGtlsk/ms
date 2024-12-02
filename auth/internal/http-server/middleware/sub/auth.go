package sub

import (
	"auth-ms/internal/service/auth"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
)

type Auth struct {
	jwtParser auth.JwtParser
}

func NewAuth(jwtParser auth.JwtParser) *Auth {
	return &Auth{jwtParser: jwtParser}
}

func (a *Auth) New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			tokenString := r.Header.Get("Authorization")

			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if len(strings.Split(tokenString, " ")) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, _, _, _, err := a.jwtParser.ParseToken(tokenString)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
