package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/YaleSpinup/f5-api/common"
	"github.com/YaleSpinup/f5-api/ltm"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

var (
	// Org will carry throughout the api and get tagged on resources
	Org string
)

type server struct {
	router      *mux.Router
	version     common.Version
	context     context.Context
	LTMServices map[string]ltm.LTMIface
}

// NewServer creates a new server and starts it
func NewServer(config common.Config) error {
	// setup server context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := server{
		router:      mux.NewRouter(),
		version:     config.Version,
		context:     ctx,
		LTMServices: make(map[string]ltm.LTMIface),
	}

	if config.Org == "" {
		return errors.New("'org' cannot be empty in the configuration")
	}
	Org = config.Org

	// Create shared F5 BigIP sessions
	for name, c := range config.Accounts {
		s.LTMServices[name] = ltm.NewSession(c.LTMHost, c.Username, c.Password, c.UploadPath)
	}

	publicURLs := map[string]string{
		"/v1/f5/ping":    "public",
		"/v1/f5/version": "public",
		"/v1/f5/metrics": "public",
	}

	// load routes
	s.routes()

	if config.ListenAddress == "" {
		config.ListenAddress = ":8080"
	}
	handler := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, TokenMiddleware(config.Token, publicURLs, s.router)))
	srv := &http.Server{
		Handler:      handler,
		Addr:         config.ListenAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("Starting listener on %s", config.ListenAddress)
	if err := srv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// LogWriter is an http.ResponseWriter
type LogWriter struct {
	http.ResponseWriter
}

// Write log message if http response writer returns an error
func (w LogWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	if err != nil {
		log.Errorf("Write failed: %v", err)
	}
	return
}
