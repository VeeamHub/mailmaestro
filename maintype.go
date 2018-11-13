package main

type MailMaestroConfig struct {
	LDAPServer   string `json:"LDAPServer"`
	LDAPUser     string `json:"LDAPUser"`
	LDAPPassword string `json:"LDAPPassword"`
	LDAPPort     int    `json:"LDAPPort"`
	LDAPPortSec  int    `json:"LDAPPortSec"`
	LDAPSecure   bool   `json:"LDAPSecure"`
	LDAPBase     string `json:"LDAPBase"`
	VBOServer    string `json:"VBOServer"`
	VBOPort      int    `json:"VBOPort"`
	VBOUsername  string `json:"VBOUsername"`
	VBOPassword  string `json:"VBOPassword"`
	VBOVersion   string `json:"VBOVersion"`
	VBOOrg       string `json:"VBOOrg"`
	VBOMailBox   string `json:"VBOMailBox"`
	LocalStop    bool   `json:"LocalStop"`
	LocalPort    int    `json:"LocalPort"`
	LocalCert    string `json:"LocalCert"`
	LocalKey     string `json:"LocalKey"`
	ManualMap    string `json:"ManualMap"`
}
