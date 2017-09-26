package main

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/ldap.v2"
)

type LDAPController struct {
	Connection *ldap.Conn
	Server     string
	BaseDN     string
	Rouser     string
	Ropass     string
	Port       int
	Portsecure int
	Execsecure bool
}

type LDAPEmailUserResult struct {
	CN       string
	DN       string
	SAM      string
	Name     string
	Email    string
	Password string
}

func (l *LDAPController) Init() error {
	var lerr error

	if l.Execsecure {
		l.Connection, lerr = ldap.DialTLS("tcp", fmt.Sprintf("%s:%d", l.Server, l.Portsecure), &tls.Config{InsecureSkipVerify: true})
	} else {
		l.Connection, lerr = ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.Server, l.Port))
	}

	if lerr == nil {
		lerr = l.Connection.Bind(l.Rouser, l.Ropass)
	}

	if lerr != nil && l.Connection != nil {
		l.Connection.Close()
		l.Connection = nil
	}
	return lerr
}

func (l *LDAPController) SearchEmailUser(email string) (error, *LDAPEmailUserResult) {
	var lerr error
	var u *LDAPEmailUserResult

	searchRequest := ldap.NewSearchRequest(
		l.BaseDN,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf("(mail=%s)", email),
		[]string{"name", "mail", "cn", "distinguishedName", "samaccountname"},
		nil)

	var sr *ldap.SearchResult
	sr, lerr = l.Connection.Search(searchRequest)

	if lerr == nil && len(sr.Entries) > 0 {
		entry := sr.Entries[0]
		var cn, ds, name, sam, mail string

		for _, attr := range entry.Attributes {
			switch attr.Name {
			case "cn":
				{
					cn = attr.Values[0]
				}
			case "distinguishedName":
				{
					ds = attr.Values[0]
				}
			case "name":
				{
					name = attr.Values[0]
				}
			case "sAMAccountName":
				{
					sam = attr.Values[0]
				}
			case "mail":
				{
					mail = attr.Values[0]
				}
			}
		}
		u = &LDAPEmailUserResult{cn, ds, sam, name, mail, ""}
	}
	return lerr, u
}
func (l *LDAPController) ValidateLogin(user *LDAPEmailUserResult, password string) bool {
	err := l.Connection.Bind(user.DN, password)
	success := false
	if err == nil {
		success = true
	}
	return success
}
func (l *LDAPController) Close() {
	if l.Connection != nil {
		l.Connection.Close()
	}
}
