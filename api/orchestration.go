package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/f5-api/ltm"
)

type ltmOrchestrator struct {
	client     ltm.LTMIface
	org        string
	UploadPath string
}

func (o *ltmOrchestrator) modifyClientSSLProfile(ctx context.Context, data *ModifyClientSSLProfileRequest) error {
	var err error
	var ecert []byte
	var ekey []byte

	// decode certificate and key file
	ecert, err = base64.StdEncoding.DecodeString(data.CertificateFile)
	if err != nil {
		return err
	}
	ekey, err = base64.StdEncoding.DecodeString(data.KeyFile)
	if err != nil {
		return err
	}

	// TODO Check cert/key match

	// upload certificate and key file
	err = o.client.UploadFile(string(ecert), fmt.Sprintf("%s.crt", data.ClientSSLProfileName))
	if err != nil {
		return err
	}
	err = o.client.UploadFile(string(ekey), fmt.Sprintf("%s.key", data.ClientSSLProfileName))
	if err != nil {
		return err
	}

	thisYear := fmt.Sprintf("%s", time.Now().Format("2006"))

	err = o.client.ImportCertificate(data.ClientSSLProfileName, thisYear)
	if err != nil {
		return err
	}

	err = o.client.ImportKey(data.ClientSSLProfileName, thisYear)
	if err != nil {
		return err
	}

	// update clientssl profile, i.e., realcert.lab.example.org-2021.(crt|key)
	err = o.client.ModifyClientSSLProfile(data.ClientSSLProfileName, data.DefaultsFrom, data.Chain, data.CipherGroup, data.Ciphers, thisYear)
	if err != nil {
		return err
	}

	return nil
}

func (o *ltmOrchestrator) createClientSSLProfile(ctx context.Context, data *ModifyClientSSLProfileRequest) error {
	var err error
	var ecert []byte
	var ekey []byte

	// decode certificate and key file
	ecert, err = base64.StdEncoding.DecodeString(data.CertificateFile)
	if err != nil {
		return err
	}
	ekey, err = base64.StdEncoding.DecodeString(data.KeyFile)
	if err != nil {
		return err
	}

	// TODO Check cert/key match

	// upload certificate and key file
	err = o.client.UploadFile(string(ecert), fmt.Sprintf("%s.crt", data.ClientSSLProfileName))
	if err != nil {
		return err
	}
	err = o.client.UploadFile(string(ekey), fmt.Sprintf("%s.key", data.ClientSSLProfileName))
	if err != nil {
		return err
	}

	thisYear := fmt.Sprintf("%s", time.Now().Format("2006"))

	err = o.client.ImportCertificate(data.ClientSSLProfileName, thisYear)
	if err != nil {
		return err
	}

	err = o.client.ImportKey(data.ClientSSLProfileName, thisYear)
	if err != nil {
		return err
	}

	// create clientssl profile, i.e., realcert.lab.example.org-2021.(key|crt}
	err = o.client.CreateClientSSLProfile(data.ClientSSLProfileName, data.DefaultsFrom, data.Chain, data.CipherGroup, data.Ciphers, thisYear)
	if err != nil {
		return err
	}

	return nil
}

func (o *ltmOrchestrator) deleteClientSSLProfile(ctx context.Context, name string) error {

	clientSSLProfile, err := o.client.GetClientSSLProfile(name)
	if err != nil {
		return err
	}

	if clientSSLProfile == nil {
		return apierror.New(apierror.ErrNotFound, fmt.Sprintf("%s not found", name), nil)
	}

	err = o.client.RemoveClientSSLProfile(name)
	if err != nil {
		return err
	}

	err = o.client.RemoveCertificate(clientSSLProfile.Cert)
	if err != nil {
		return err
	}

	err = o.client.RemoveKey(clientSSLProfile.Key)
	if err != nil {
		return err
	}

	return nil
}
