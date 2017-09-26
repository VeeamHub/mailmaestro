package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
)

func checkvalnotempty(check string, iferr string) bool {
	retval := false
	if check == "" {
		flag.PrintDefaults()
		fmt.Printf("\n\nPlease supply: %s\n", iferr)
		retval = true
	}
	return retval
}

func lookForFolder(mbxFolders RestoreMailboxFolderPage, id string) *RestoreMailboxFolder {
	for _, folder := range mbxFolders.Folders {
		if folder.Id == id {
			return folder
		}
	}
	return nil
}
func main() {
	var pport = flag.Int("port", 4123, "Mail maestro web port")
	/*
		// For windows : http://gnuwin32.sourceforge.net/packages/openssl.htm
		set OPENSSL_CONF=C:\d\openssl\share\openssl.cnf
		openssl.exe genrsa 2048 > key.pem
		openssl.exe req -new -x509 -key key.pem -out cert.pem -days 3650
	*/
	var pcert = flag.String("cert", "", "Cert file for https, if none is provided, it will start in http modus")
	var pkey = flag.String("key", "", "Key file for https, if none is provided, it will start in http modus")
	var plocalstop = flag.Bool("lstop", false, "Allow typing stop to stop the service")

	var pldapserver = flag.String("ldapserver", "", "ldap server")
	var pldapreadonlyuser = flag.String("ldapreadonlyuser", "", "ldap username")
	var pldapreadonlypassword = flag.String("ldapreadonlypassword", "", "ldap password")
	var pldapbaseDN = flag.String("ldapbaseDN", "", "baseDN")

	var pldapport = flag.Int("ldapport", 389, "ldapport override")
	var pldapportsec = flag.Int("ldapsecport", 636, "ldapport override")
	var pldapsec = flag.Bool("ldapsecure", false, "ldap secure")

	var pvborestserver = flag.String("vboserver", "localhost", "Rest API for VBO 365")
	var pvboport = flag.Int("vboport", 4443, "Rest API port for VBO365")
	var pvboversion = flag.String("vboversion", "v1", "Rest API version")

	var pvbousername = flag.String("vbouser", "", "Rest API user")
	var pvbopassword = flag.String("vbopassword", "", "Rest API password")

	var pvboorg = flag.String("vboorganization", "", "Will default to first one found if none is supplied")
	var pmailboxuser = flag.String("mailboxuser", "", "supply email address")

	var pconfig = flag.String("config", "maestroconf.json", "Use a json config file")
	var pdump = flag.String("dump", "", "Dump config to file passed as a parameter")

	flag.Parse()

	//making an empty config and setting defaults
	//then parse a config if it is given
	config := MailMaestroConfig{}
	_, configerr := os.Stat(*pconfig)
	if *pconfig != "" && !os.IsNotExist(configerr) {
		jf, configerr := ioutil.ReadFile(*pconfig)
		if configerr == nil {
			//log.Printf("%s", jf)
			json.Unmarshal(jf, &config)
		} else {
			log.Fatal("Problem parsing config file")
		}
	}
	if config.LocalPort == 0 {
		config.LocalPort = 4123
	}
	if config.LDAPPort == 0 {
		config.LDAPPort = 389
	}
	if config.LDAPPortSec == 0 {
		config.LDAPPortSec = 636
	}
	if config.VBOPort == 0 {
		config.LDAPPortSec = 4443
	}
	if config.VBOVersion == "" {
		config.VBOVersion = "v1"
	}
	if config.VBOPort == 0 {
		config.VBOPort = 4443
	}

	//adding with cmd line feature
	if *pport != 4123 {
		config.LocalPort = *pport
	}
	if *pkey != "" {
		config.LocalKey = *pkey
	}
	if *pcert != "" {
		config.LocalCert = *pcert
	}
	if *plocalstop {
		config.LocalStop = true
	}
	if *pldapserver != "" {
		config.LDAPServer = *pldapserver
	}
	if *pldapreadonlyuser != "" {
		config.LDAPUser = *pldapreadonlyuser
	}
	if *pldapreadonlypassword != "" {
		config.LDAPPassword = *pldapreadonlypassword
	}
	if *pldapbaseDN != "" {
		config.LDAPBase = *pldapbaseDN
	}
	if *pldapport != 389 {
		config.LDAPPort = *pldapport
	}
	if *pldapportsec != 636 {
		config.LDAPPortSec = *pldapportsec
	}
	if *pldapsec {
		config.LDAPSecure = true
	}
	if *pvboorg != "" {
		config.VBOOrg = *pvboorg
	}
	if *pmailboxuser != "" {
		config.VBOMailBox = *pmailboxuser
	}
	if *pvbousername != "" {
		config.VBOUsername = *pvbousername
	}
	if *pvbopassword != "" {
		config.VBOPassword = *pvbopassword
	}
	if *pvboport != 4443 {
		config.VBOPort = 4443
	}
	if *pvborestserver != "localhost" {
		config.VBOServer = *pvborestserver
	}
	if *pvboversion != "v1" {
		config.VBOVersion = *pvboversion
	}

	if *pdump != "" {
		log.Printf(*pdump)
		exportcopy := config
		exportcopy.LDAPPassword = ""
		exportcopy.VBOPassword = ""
		b, derr := json.MarshalIndent(exportcopy, "", " ")
		if derr == nil {
			derr = ioutil.WriteFile(*pdump, b, 0644)
		}
		if derr != nil {
			log.Printf("Error dumping config %s", derr.Error())
		} else {
			log.Printf("Dumped config to %s", *pdump)
		}
	}

	starthttps := false
	if (config.LocalKey == "" && config.LocalCert != "") || (config.LocalKey != "" && config.LocalCert == "") {
		fmt.Println("Please provide both key and cert, one of them is empty. If both are empty, this will start in http mode")
		return
	}

	if config.LocalKey != "" {
		if _, err := os.Stat(config.LocalKey); os.IsNotExist(err) {
			fmt.Println("Key does not exist")
			return
		}
		if _, err := os.Stat(config.LocalCert); os.IsNotExist(err) {
			fmt.Println("Cert does not exist")
			return
		}
		starthttps = true
	}

	if checkvalnotempty(config.LDAPServer, "ldap server name") {
		return
	}
	if checkvalnotempty(config.LDAPUser, "read only user to bind") {
		return
	}
	if config.LDAPPassword == "" {
		fmt.Printf("Please provide a password for the LDAP User: ")
		bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Println(err)
			return
		}
		config.LDAPPassword = string(bytePassword)
		fmt.Println("")
	}
	if checkvalnotempty(config.VBOMailBox, "mailbox user") {
		return
	}

	if checkvalnotempty(config.VBOUsername, "rest api username") {
		return
	}

	if config.VBOPassword == "" {
		fmt.Printf("Please provide a password for the Rest API User: ")
		bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Println(err)
			return
		}
		config.VBOPassword = string(bytePassword)
		fmt.Println("")
	}

	if config.LDAPBase == "" {
		config.LDAPBase = fmt.Sprintf("DC=%s", strings.Join(strings.Split(config.LDAPServer, "."), ",DC="))
		log.Printf("Base DN is empty, tried constructing it from ldapserver name %s", config.LDAPBase)
	}

	log.Println("Trying to contact the LDAP Controller...")

	ldapc := LDAPController{Server: config.LDAPServer, BaseDN: config.LDAPBase, Rouser: config.LDAPUser, Ropass: config.LDAPPassword, Port: config.LDAPPort, Portsecure: config.LDAPPortSec, Execsecure: config.LDAPSecure}
	err := ldapc.Init()
	defer ldapc.Close()
	if err != nil {
		log.Fatal(err)
	}
	err, user := ldapc.SearchEmailUser(config.VBOMailBox)
	if err != nil {
		log.Fatal(err)
	}
	if user == nil {
		log.Fatal("Could not find user emailbox")
	} else {
		log.Printf("Found %s", user.DN)
		log.Printf("User will be able to login with SAM: %s or Email : %s", user.SAM, user.Email)

		restapi := RestAPI{BaseURL: fmt.Sprintf("https://%s:%d/%s", config.VBOServer, config.VBOPort, config.VBOVersion), Username: config.VBOUsername, Password: config.VBOPassword}
		err := restapi.Init()
		defer restapi.Close()
		if err == nil {
			log.Printf("Was able to login on rest %s, token is %s", restapi.BaseURL, restapi.Authentication.AccessToken)
			var orgs []Organization

			err, orgs = restapi.GetOrganizations()
			if err == nil && len(orgs) > 0 {
				orgSelected := orgs[0]
				if config.VBOOrg != "" {
					for _, org := range orgs {
						if org.Name == config.VBOOrg {
							orgSelected = org
						}
					}
				}
				log.Printf("Using Organization %s (%s) - point %s", orgSelected.Id, orgSelected.Name, orgSelected.LastBackupTime)
				var restoreSession RestoreSession
				err, restoreSession = restapi.StartRestoreAction(orgSelected.Id, orgSelected.LastBackupTime)

				for key, link := range restoreSession.Links {
					log.Printf("%s %s", key, link.Href)
				}

				if err == nil {
					defer restapi.StopRestoreSession(restoreSession.Id)
					log.Printf("Started restore session %s State %s Result %s", restoreSession.PointInTime, restoreSession.State, restoreSession.Result)
					var mbxp RestoreMailboxPage
					var mbxUser *RestoreMailbox
					err, mbxp = restapi.GetRestoreMailboxes(restoreSession.Id)
					if err == nil {
						//log.Printf("Count MBX %d",mbxp.Limit)
						for _, mbx := range mbxp.MailBoxes {
							if mbx.Email == config.VBOMailBox {
								mbxUser = mbx
							}
						}
						if mbxUser.Id != "" {
							log.Printf("Found users mailbox %s (%s)", mbxUser.Id, mbxUser.Email)
							log.Printf("Can start the demo page")

							stop := make(chan int, 1)
							sess := make(map[string]GUIServerSession)
							users := make(map[string]*LDAPEmailUserResult)

							users[user.Email] = user

							gui := GUIServeHTTP{LDAPController: &ldapc,
								LDAPUsers:             &(users),
								RestAPI:               &restapi,
								SessionIDLock:         &(sync.RWMutex{}),
								Sessions:              &(sess),
								RestoreSession:        &restoreSession,
								RestoreMailbox:        mbxUser,
								GracefullshutdownChan: &stop,
							}
							server := &http.Server{Addr: fmt.Sprintf(":%d", config.LocalPort), Handler: gui}

							go func() {
								if starthttps {

									if err := server.ListenAndServeTLS(config.LocalCert, config.LocalKey); err != nil && err != http.ErrServerClosed {
										log.Fatal(fmt.Sprintf("Server error %s", err))
									}
								} else {
									if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
										log.Fatal(fmt.Sprintf("Server error %s", err))
									}
								}
							}()
							sec := ""
							if starthttps {
								sec = "s"
							}
							log.Printf("Serving on http%s://localhost:%d", sec, config.LocalPort)

							if config.LocalStop {
								commandchan := make(chan string)

								go func(commandchan chan string) {
									reader := bufio.NewReader(os.Stdin)

									for stopreading := false; !stopreading; {
										fmt.Printf("> ")
										s, err := reader.ReadString('\n')
										if err == nil {
											if len(s) > 3 && strings.ToLower(s)[:4] == "stop" {
												stopreading = true
												commandchan <- "stop"
											}
										}
									}

								}(commandchan)

								for needstop := false; !needstop; {
									select {
									case <-stop:
										needstop = true
									case <-commandchan:
										needstop = true
									}
								}
							} else {
								<-stop
							}

							server.Shutdown(context.Background())
							log.Printf("Cleaned Up...")

						} else {
							log.Printf("Did not find users mailbox, this demo is limited to one page (30 mailboxes)")
							log.Fatal(err)
						}
					} else {
						log.Printf("Error getting mailboxes")
						log.Fatal(err)
					}
				} else {
					log.Printf("Error starting restore session")
					log.Fatal(err)
				}
			} else {
				log.Printf("Error fetching organization or non defined")
				log.Fatal(err)
			}

		} else {
			log.Printf("Could not login on the rest server")
			log.Fatal(err)
		}
	}
}
