package ltm

import (
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleUniversity/go-bigip"
	log "github.com/sirupsen/logrus"
)

// ListClientSSLProfiles lists the client ssl profiles
func (l *LTM) ListClientSSLProfiles() ([]string, error) {
	out, err := l.Service.ClientSSLProfiles()
	if err != nil {
		msg := fmt.Sprintf("failed to list client ssl profiles")
		return nil, apierror.New(apierror.ErrInternalError, msg, err)
	}

	profiles := make([]string, 0, len(out.ClientSSLProfiles))
	for _, p := range out.ClientSSLProfiles {
		profiles = append(profiles, p.Name)
	}

	return profiles, nil
}

// GetClientSSLProfile gets a client ssl profile from ltm
func (l *LTM) GetClientSSLProfile(name string) (*bigip.ClientSSLProfile, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	out, err := l.Service.GetClientSSLProfile(name)
	if err != nil {
		msg := fmt.Sprintf("failed to get ssl profile %s on %s", name, l.Host)
		return nil, apierror.New(apierror.ErrInternalError, msg, err)
	}

	return out, nil
}

// UploadFile uploads a file to an ltm
func (l *LTM) UploadFile(file, name string) error {
	if file == "" || name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if _, err := l.Service.UploadBytes([]byte(file), name); err != nil {
		msg := fmt.Sprintf("failed to upload file %s on %s", name, l.Host)
		return apierror.New(apierror.ErrInternalError, msg, err)
	}

	log.Infof("uploaded file %s on host %s", name, l.Host)

	return nil
}

// ImportCertificate imports certificate to System SSL
func (l *LTM) ImportCertificate(name, thisYear string) error {
	if name == "" || thisYear == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	addcert := &bigip.Certificate{
		Name:       fmt.Sprintf("%s-%s.crt", name, thisYear),
		SourcePath: fmt.Sprintf("file:%s/%s.crt", l.UploadPath, name),
	}

	err := l.Service.AddCertificate(addcert)
	if err != nil {
		// todo: button-up with more logic perhaps, make a call to cert/key
		// api, and look to see if it exists first, but failing to 'add' isn't
		// terrible.
		// Mostly, this soft fail is to support cert renewals and profile
		// option changes to LTM client-ssl profiles - we are ok if the
		// cert/key is already on the system SSL don't return, just log the
		// condition and move on...
		log.Infof("add cert error on host %s: %s, proceeding...", l.Host, err)
	} else {
		log.Infof("added cert %s on host %s", name, l.Host)
	}

	return nil
}

// ImportKey Imports Key to System SSL
func (l *LTM) ImportKey(name, thisYear string) error {
	if name == "" || thisYear == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	addkey := &bigip.Key{
		Name:       fmt.Sprintf("%s-%s.key", name, thisYear),
		SourcePath: fmt.Sprintf("file:%s/%s.key", l.UploadPath, name),
	}

	err := l.Service.AddKey(addkey)
	if err != nil {
		// See AddCertificate comment above
		log.Infof("add key error on host %s: %s, proceeding...", l.Host, err)
	} else {
		log.Infof("added key %s on host %s", name, l.Host)
	}

	return nil
}

// ModifyClientSSLProfile update cert and key on a client-ssl profile
func (l *LTM) ModifyClientSSLProfile(ClientSSLProfileName, DefaultsFrom, Chain, CipherGroup, Ciphers, thisYear string) error {
	if ClientSSLProfileName == "" || DefaultsFrom == "" || Chain == "" || CipherGroup == "" || Ciphers == "" || thisYear == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	modifycert := &bigip.ClientSSLProfile{
		Name:         ClientSSLProfileName,
		Cert:         fmt.Sprintf("%s-%s.crt", ClientSSLProfileName, thisYear),
		Key:          fmt.Sprintf("%s-%s.key", ClientSSLProfileName, thisYear),
		Chain:        Chain,
		DefaultsFrom: DefaultsFrom,
		CipherGroup:  CipherGroup,
		Ciphers:      Ciphers,
	}

	if err := l.Service.ModifyClientSSLProfile(ClientSSLProfileName, modifycert); err != nil {
		msg := fmt.Sprintf("failed to modify client-ssl profile %s on %s", ClientSSLProfileName, l.Host)
		return apierror.New(apierror.ErrInternalError, msg, err)
	}

	log.Infof("modified client-ssl profile %s on %s\n", ClientSSLProfileName, l.Host)

	return nil

}

// CreateClientSSLProfile creates cert and key on a client-ssl profile
func (l *LTM) CreateClientSSLProfile(ClientSSLProfileName, DefaultsFrom, Chain, CipherGroup, Ciphers, thisYear string) error {
	if ClientSSLProfileName == "" || DefaultsFrom == "" || Chain == "" || CipherGroup == "" || Ciphers == "" || thisYear == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	addcert := &bigip.ClientSSLProfile{
		Name:         ClientSSLProfileName,
		Cert:         fmt.Sprintf("%s-%s.crt", ClientSSLProfileName, thisYear),
		Key:          fmt.Sprintf("%s-%s.key", ClientSSLProfileName, thisYear),
		Chain:        Chain,
		DefaultsFrom: DefaultsFrom,
		CipherGroup:  CipherGroup,
		Ciphers:      Ciphers,
	}

	if err := l.Service.AddClientSSLProfile(addcert); err != nil {
		msg := fmt.Sprintf("error creating client-ssl profile %s on %s", ClientSSLProfileName, l.Host)
		return apierror.New(apierror.ErrBadRequest, msg, err)
	}

	log.Infof("created client-ssl profile %s on host %s\n", ClientSSLProfileName, l.Host)

	return nil

}
