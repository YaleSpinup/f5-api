package api

import (
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// TokenMiddleware checks the tokens for non-public URLs
func TokenMiddleware(psk []byte, public map[string]string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Processing token middleware for protected URLs")

		// Handle CORS preflight checks
		if r.Method == "OPTIONS" {
			log.Info("Setting CORS preflight options and returning")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "X-Auth-Token")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
			return
		}

		uri, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			log.Error("Unable to parse request URI ", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if _, ok := public[uri.Path]; ok {
			log.Debugf("Not authenticating for '%s'", uri.Path)
		} else {
			log.Debugf("Authenticating token for protected URL '%s'", r.URL)

			htoken := r.Header.Get("X-Auth-Token")
			if err := bcrypt.CompareHashAndPassword([]byte(htoken), psk); err != nil {
				log.Warnf("Unable to authenticate session for URL '%s': '%s'", r.URL, err)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			log.Infof("Successfully authenticated token for URL '%s'", r.URL)
		}

		h.ServeHTTP(w, r)
	})
}
