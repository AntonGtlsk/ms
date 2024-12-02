package validate

import (
	"auth-ms/internal/service/auth"
	"fmt"
	"logging"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

func New(logging logging.LoggerInterface, jwtParser *auth.JwtParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			logging.Errorf("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if len(strings.Split(tokenString, " ")) != 2 {
			logging.Errorf("Unauthorized")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, _, _, _, err := jwtParser.ParseToken(tokenString)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		render.JSON(w, r, http.StatusOK)
	}
}
