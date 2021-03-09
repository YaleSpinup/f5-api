package ltm

import (
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleUniversity/go-bigip"
	log "github.com/sirupsen/logrus"
)

// MyClientSSLProfile gets json from a incoming post to send to ClientSSLProfile requests
type MyClientSSLProfile struct {
	Name         string `json:"clientssl-profile"`
	Cert         string `json:"cert"`
	Key          string `json:"key"`
	Chain        string `json:"chain"`
	DefaultsFrom string `json:"defaultsfrom"`
	CipherGroup  string `json:"ciphergroup"`
	Ciphers      string `json:"ciphers"`
}

// ListClientSSLProfiles lists the client ssl profiles
func (l *LTM) ListClientSSLProfiles() ([]string, error) {
	out, err := l.Service.ClientSSLProfiles()
	if err != nil {
		// GOOGLEY FIXME
		//return nil, ErrCode("failed to list client ssl profiles", err)
		return nil, err
	}

	//log.Debugf("got output from list client ssl profiles: %+v", out)

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
		// GOOGLEY FIXME
		//return nil, ErrCode("failed to get ssl profile", err)
		return nil, err
	}

	// do something better with this
	//log.Debugf("output from get client ssl profiles: %+v", out)

	return out, nil
}

// UploadFile uploads a file to an ltm
func (l *LTM) UploadFile(file, name string) error {
	if file == "" || name == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	uploads, err := l.Service.UploadBytes([]byte(file), name)
	if err != nil {
		//return ErrCode("failed to upload file on ", err)
		msg := fmt.Sprintf("failed to upload file %s on %s", name, l.Host)
		return apierror.New(apierror.ErrBadRequest, msg, err)
	}

	// do something better with this
	// what is in the uploads reference
	log.Debugf("uploadfile content: %v\n", uploads)

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
		log.Infof("Addcert error on %s, proceeding...", l.Host)
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
		// See AddCertificate comment
		log.Infof("Addkey error on %s, proceeding...", l.Host)
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

	fmt.Printf("modifycert: %v\n", modifycert)

	err := l.Service.ModifyClientSSLProfile(ClientSSLProfileName, modifycert)
	if err != nil {
		fmt.Printf("modify ClientSSL profile error on %s: %s\n", ClientSSLProfileName, err)
	} else {
		fmt.Printf("modified ClientSSL profile: %s\n", ClientSSLProfileName)
	}
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

	fmt.Printf("addcert: %v\n", addcert)

	err := l.Service.AddClientSSLProfile(addcert)
	if err != nil {
		fmt.Printf("create ClientSSL profile error on %s: %s\n", ClientSSLProfileName, err)
	} else {
		fmt.Printf("created ClientSSL profile: %s\n", ClientSSLProfileName)
	}
	return nil

}
