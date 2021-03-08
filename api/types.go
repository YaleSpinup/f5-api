package api

// ClientSSLProfile is an ltm clientSSL Profile
type ClientSSLProfile struct {
	Cert                 string `json:"cert"`
	Key                  string `json:"key"`
	Chain                string `json:"chain"`
	DefaultsFrom         string `json:"defaultsfrom"`
	CipherGroup          string `json:"ciphergroup"`
	Ciphers              string `json:"ciphers"`
	ClientSSLProfileName string `json:"clientssl-profile"`
}

// ModifyClientSSLProfileRequest defines the key and cert data uploaded from a client
type ModifyClientSSLProfileRequest struct {
	CertificateFile      string `json:"cert"`
	KeyFile              string `json:"key"`
	ClientSSLProfileName string `json:"clientssl-profile"`
	Chain                string `json:"chain"`
	DefaultsFrom         string `json:"defaultsfrom"`
	Ciphers              string `json:"ciphers"`
	CipherGroup          string `json:"ciphergroup"`
	ClientSSLProfile     *ClientSSLProfile
}
