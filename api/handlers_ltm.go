package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// ListClientSSLProfiles List Client SSL Profiles on LTM
func (s *server) ListClientSSLProfiles(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	host := vars["host"]

	log.Infof("list client ssl profiles %s", host)

	ltmService, ok := s.LTMServices[host]
	if !ok {
		msg := fmt.Sprintf("LTM host service not found for account: %s", host)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}

	out, err := ltmService.ListClientSSLProfiles()
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(out)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// ShowClientSSLProfile Show detail of Client SSL Profile on LTM
func (s *server) ShowClientSSLProfile(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	host := vars["host"]
	name := vars["name"]

	log.Infof("getting details about client ssl profile %s", name)

	ltmService, ok := s.LTMServices[host]
	if !ok {
		msg := fmt.Sprintf("LTM host service not found for account: %s", host)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}

	out, err := ltmService.GetClientSSLProfile(name)
	if err != nil {
		handleError(w, err)
	}

	j, err := json.Marshal(out)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrBadRequest, "failed to marshal json", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// ModifyClientSSLProfile updates a clientssl profile including updating the cert and key if supplied in the body
func (s *server) ModifyClientSSLProfile(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	host := vars["host"]
	name := vars["name"]

	log.Infof("update client-ssl profile %s on host %s", name, host)

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	data := ModifyClientSSLProfileRequest{}
	if err := json.Unmarshal(raw, &data); err != nil {
		handleError(w, err)
	}

	ltmService, ok := s.LTMServices[host]
	if !ok {
		msg := fmt.Sprintf("LTM host service not found for account: %s", host)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}

	orch := &ltmOrchestrator{
		client: ltmService,
	}

	if err := orch.modifyClientSSLProfile(r.Context(), &data); err != nil {
		handleError(w, err)
		return
	}

	out := []byte(fmt.Sprintf("modified client-ssl profile %s on host %s", name, host))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// CreateClientSSLProfile creates SSL Client Profile
func (s *server) CreateClientSSLProfile(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w} vars := mux.Vars(r)
	host := vars["host"]
	name := vars["name"]

	log.Infof("create client-ssl profile %s on host %s", name, host)

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err)
		return
	}
	defer r.Body.Close()

	data := ModifyClientSSLProfileRequest{}
	if err := json.Unmarshal(raw, &data); err != nil {
		handleError(w, err)
	}

	ltmService, ok := s.LTMServices[host]
	if !ok {
		msg := fmt.Sprintf("LTM host service not found for account: %s", host)
		handleError(w, apierror.New(apierror.ErrNotFound, msg, nil))
		return
	}

	orch := &ltmOrchestrator{
		client: ltmService,
	}

	if err := orch.createClientSSLProfile(r.Context(), &data); err != nil {
		handleError(w, err)
		return
	}

	out := []byte(fmt.Sprintf("created client-ssl profile %s on host %s", name, host))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
