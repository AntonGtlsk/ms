package login

import (
	"logging"
	"net/http"
)

func New(logging logging.LoggerInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL := ""
		logging.Infof("Successful redirect to URL: %s", targetURL)

		http.Redirect(w, r, targetURL, http.StatusFound)
	}
}
