package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func newUUID() string {
	return uuid.New().String()
}

type GUIServerSession struct {
	SessionId string
	LDAPUser  *LDAPEmailUserResult
}
type GUIServeHTTP struct {
	LDAPController        *LDAPController
	LDAPUsers             *map[string]*LDAPEmailUserResult
	RestAPI               *RestAPI
	Sessions              *map[string]GUIServerSession
	SessionIDLock         *sync.RWMutex
	RestoreSession        *RestoreSession
	RestoreMailbox        *RestoreMailbox
	GracefullshutdownChan *chan int
}

func (g GUIServeHTTP) Auth(pw *http.ResponseWriter, r *http.Request) {
	w := *pw
	r.ParseForm()

	var err error
	var id string

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	log.Printf("Trying to validate %s", username)
	if username != "" && password != "" {
		g.SessionIDLock.Lock()
		defer g.SessionIDLock.Unlock()

		var foundUser *LDAPEmailUserResult
		for email, user := range *(g.LDAPUsers) {
			if username == email || username == user.SAM {
				foundUser = user
			}
		}

		if foundUser != nil {
			if g.LDAPController.ValidateLogin(foundUser, password) {
				id = newUUID()

				foundUser.Password = password
				(*(g.Sessions))[id] = GUIServerSession{SessionId: id, LDAPUser: foundUser}
				// fmt.Printf("\nValidated")
			}
		}
	}

	if err == nil && id != "" {
		http.Redirect(w, r, fmt.Sprintf("/%s/", id), 301)
		// fmt.Printf("\nRedirected to restore page")
	} else {
		http.Redirect(w, r, fmt.Sprintf("/"), 301)
		// fmt.Printf("\nRedirectered to main")
	}

}

func (g GUIServeHTTP) WriteRestorePage(pw *http.ResponseWriter, sessionid string, errorString string, successString string) {
	w := *pw
	restapi := g.RestAPI

	var err error

	var mbxFolders RestoreMailboxFolderPage
	err, mbxFolders = restapi.GetRestoreMailboxFolders(g.RestoreSession.Id, g.RestoreMailbox.Id)

	if err == nil {
		folders := []GUIRestoreFolder{}
		for _, fld := range mbxFolders.Folders {
			folders = append(folders, GUIRestoreFolder{fld.Id, fld.Name})
		}

		var restorePage *template.Template
		restorePage, err = template.New("restorePage").Parse(restoreTemplate)
		items := GUIRestoreItems{
			Sessionid:     sessionid,
			Folders:       folders,
			ErrorString:   errorString,
			SuccessString: successString,
		}
		restorePage.Execute(w, &items)

	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Println(err)
	}
}
func (g GUIServeHTTP) WriteLogin(pw *http.ResponseWriter) {
	w := *pw
	fmt.Fprintf(w, loginPageTemplate)
}
func (g GUIServeHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestlen := len(r.RequestURI)

	if requestlen > 8 && r.RequestURI[:8] == "/static/" {
		data, err := Asset(r.RequestURI[1:])
		if err != nil {
			http.Error(w, fmt.Sprintf("Could not find %s", r.RequestURI[1:]), http.StatusNotFound)
		} else {
			ext := strings.ToLower(filepath.Ext(r.RequestURI[1:]))
			ctype := "text/plain"
			switch ext {
			case ".css":
				{
					ctype = "text/css"
				}
			case ".js":
				{
					ctype = "application/javascript"
				}
			default:
				{
					// fmt.Printf("What is %s",ext)
				}
			}
			w.Header().Set("Content-Type", ctype)
			io.Copy(w, bytes.NewReader(data))
		}
	} else if requestlen > 37 {
		sessid := r.RequestURI[1:37]
		_, err := uuid.Parse(sessid)
		if err == nil && r.RequestURI[37] == '/' {
			g.SessionIDLock.RLock()
			defer g.SessionIDLock.RUnlock()

			currentSession, foundid := (*(g.Sessions))[sessid]

			if foundid {
				subrequest := strings.Split(r.RequestURI[37:], "/")
				if len(subrequest) > 1 {
					switch subrequest[1] {
					case "stop":
						{
							fmt.Fprintf(w, "Stopping...")
							go func(pc *chan int) {
								(*pc) <- 1
							}(g.GracefullshutdownChan)
						}
					case "folderItems":
						{
							var resp GUIRestoreItems

							if len(subrequest) > 2 && subrequest[2] != "" {
								var mbxItems RestoreMailboxItemPage
								err, mbxItems = g.RestAPI.GetRestoreMailboxFolderItems(g.RestoreSession.Id, g.RestoreMailbox.Id, subrequest[2])

								if err == nil {
									mails := []GUIRestoreItemMail{}
									cals := []GUIRestoreItemAppointment{}
									contacts := []GUIRestoreItemContact{}
									tasks := []GUIRestoreItemTask{}

									for _, item := range mbxItems.Items {
										parent := ""

										if refsplit := strings.Split(item.Links["parent"].Href, "/"); len(refsplit) > 1 {
											parent = refsplit[len(refsplit)-1]
										}

										switch item.ItemClass {
										case "IPM.Appointment":
											{
												// log.Printf("%sMeeting Organizer %s Location %s Time %s", parent, item.Organizer, item.Location, item.StartTime)
												cals = append(cals, GUIRestoreItemAppointment{item.Id, parent, item.Organizer, item.Location, item.StartTime})
											}
										case "IPM.Schedule.Meeting.Resp.Pos":
											{
												// log.Printf("%sMeeting Res From %s To %s CC %s BCC %s Subject %s", parent, item.From, item.To, item.CC, item.BCC, item.Subject)
												mails = append(mails, GUIRestoreItemMail{item.Id, parent, item.From, item.To, item.Subject})
											}
										case "IPM.Schedule.Meeting.Request":
											{
												// log.Printf("%sMeeting Req From %s To %s CC %s BCC %s Subject %s", parent, item.From, item.To, item.CC, item.BCC, item.Subject)
												mails = append(mails, GUIRestoreItemMail{item.Id, parent, item.From, item.To, item.Subject})
											}
										case "IPM.Note":
											{
												// log.Printf("%sMail From %s To %s CC %s BCC %s Subject %s", parent, item.From, item.To, item.CC, item.BCC, item.Subject)
												mails = append(mails, GUIRestoreItemMail{item.Id, parent, item.From, item.To, item.Subject})
											}
										case "IPM.Contact":
											{
												contacts = append(contacts, GUIRestoreItemContact{item.Id, parent, item.DisplayAs, item.Email})
											}
										case "IPM.Task":
											{
												tasks = append(tasks, GUIRestoreItemTask{item.Id, parent, item.Owner})
											}
										default:
											{
												log.Printf("%s%s From %s To %s CC %s BCC %s Subject %s", parent, item.ItemClass, item.From, item.To, item.CC, item.BCC, item.Subject)
											}
										}
									}
									resp = GUIRestoreItems{ErrorString: "", Mails: mails, Appointments: cals, Contacts: contacts, Tasks: tasks}
								} else {
									resp = GUIRestoreItems{ErrorString: fmt.Sprintf("Error : %s", err)}
								}
							} else {
								resp = GUIRestoreItems{ErrorString: "Please provide an Id"}
							}
							json.NewEncoder(w).Encode(resp)
						}
					case "restorejson":
						{
							var resp GUIRestoreResponse

							if len(subrequest) > 2 {
								errRestore, feedback := g.RestAPI.DoRestMailboxItem(g.RestoreSession.Id, g.RestoreMailbox.Id, subrequest[2], currentSession.LDAPUser.Email, currentSession.LDAPUser.Password)
								if errRestore == nil && len(feedback.Exceptions) == 0 {
									resp = GUIRestoreResponse{"success", "", feedback.CreatedItemsCount, feedback.MergedItemsCount, feedback.SkippedItemsCount, feedback.FailedItemsCount}
								} else if errRestore == nil {
									resp = GUIRestoreResponse{"failure", "Backend Rest API", 0, 0, 0, 0}

									log.Printf("Exceptions during restore: ")
									for _, excep := range feedback.Exceptions {
										log.Printf(" EX-%s", excep)
									}
								} else {
									resp = GUIRestoreResponse{"failure", "Internal Server Error", 0, 0, 0, 0}

									log.Printf("Error restoring %s", errRestore)
								}
							} else {
								resp = GUIRestoreResponse{"failure", "Not Enough Parameters in URL request", 0, 0, 0, 0}
							}
							json.NewEncoder(w).Encode(resp)
						}
					default:
						{
							g.WriteRestorePage(&w, sessid, "", "")
						}
					}
				} else {
					g.WriteRestorePage(&w, sessid, "", "")
				}
			} else {
				// fmt.Printf("\nGot %d ids",len(sesssl))
				http.Error(w, "You didn't say the magic word", http.StatusForbidden)
			}
		} else {
			http.Error(w, "You didn't say the magic word", http.StatusForbidden)
		}
	} else if requestlen > 4 && r.RequestURI[:5] == "/auth" {
		g.Auth(&w, r)
	} else {
		if requestlen == 0 || (requestlen == 1 && r.RequestURI[0] == '/') {
			g.WriteLogin(&w)
		} else {
			http.Redirect(w, r, "/", 301)
		}
	}
}
