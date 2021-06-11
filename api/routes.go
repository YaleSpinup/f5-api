package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *server) routes() {

	// ltm subrouter - /v1/f5
	api := s.router.PathPrefix("/v1/f5").Subrouter()
	api.HandleFunc("/ping", s.PingHandler).Methods(http.MethodGet)
	api.HandleFunc("/version", s.VersionHandler).Methods(http.MethodGet)
	api.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	api.HandleFunc("/{host}/clientssl", s.ListClientSSLProfiles).Methods(http.MethodGet)
	api.HandleFunc("/{host}/clientssl/{name}", s.ShowClientSSLProfile).Methods(http.MethodGet)
	api.HandleFunc("/{host}/clientssl/{name}", s.DeleteClientSSLProfile).Methods(http.MethodDelete)
	api.HandleFunc("/{host}/createclientssl/{name}", s.CreateClientSSLProfile).Methods(http.MethodPut)
	api.HandleFunc("/{host}/updateclientssl/{name}", s.ModifyClientSSLProfile).Methods(http.MethodPut)
}
