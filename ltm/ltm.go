package ltm

import (
	bigip "github.com/YaleUniversity/go-bigip"
	log "github.com/sirupsen/logrus"
)

// LTMIface is abstraction for testing
type LTMIface interface {
	ListClientSSLProfiles() ([]string, error)
	GetClientSSLProfile(string) (*bigip.ClientSSLProfile, error)
	UploadFile(string, string) error
	ImportKey(string, string) error
	ImportCertificate(string, string) error
	ModifyClientSSLProfile(string, string, string, string, string, string) error
	CreateClientSSLProfile(string, string, string, string, string, string) error
	RemoveClientSSLProfile(string) error
	RemoveKey(string) error
	RemoveCertificate(string) error
}

// LTM is struct containing login info
type LTM struct {
	Service    *bigip.BigIP
	UploadPath string
	Host       string
}

// NewSession creates a new LTM session
func NewSession(host, user, pass, uploadPath string) *LTM {
	log.Infof("creating a new LTM session with host %s with username %s", host, user)

	return &LTM{
		Service:    bigip.NewSession(host, user, pass, nil),
		UploadPath: uploadPath,
		Host:       host,
	}
}
