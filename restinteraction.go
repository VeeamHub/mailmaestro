package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	// "time"
)

type RestAPI struct {
	BaseURL        string
	Username       string
	Password       string
	LoggedIn       bool
	Authentication *AuthenticationHeader
	Client         *http.Client
	APIUpdate      sync.RWMutex
	StopUpdate     chan bool
}

func (r *RestAPI) CreateUrl(req string) string {
	return fmt.Sprintf("%s/%s", r.BaseURL, req)
}
func (r *RestAPI) CreateRequest(request string, method string, body string) (*http.Request, error) {
	var err error
	var req *http.Request
	uri := request

	if len(uri) < 6 || uri[:6] != "https:" {
		uri = r.CreateUrl(request)
	}
	//fmt.Println(uri)
	if body != "" {
		req, err = http.NewRequest(method, uri, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, uri, nil)
	}

	if err == nil {
		authheader := fmt.Sprintf("Bearer %s", r.Authentication.AccessToken)
		//fmt.Println(authheader)

		req.Header.Add("Authorization", authheader)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
	}

	return req, err
}

func (r *RestAPI) Update() {
	r.APIUpdate.Lock()
	defer r.APIUpdate.Unlock()
	log.Printf("Updating ...")

	var err error
	var resp *http.Response
	resp, err = r.Client.PostForm(r.CreateUrl("token"), url.Values{"grant_type": {"refresh_token"}, "refresh_token": {r.Authentication.RefreshToken}})

	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			if resp.StatusCode < 299 {
				auth := AuthenticationHeader{}
				err = json.Unmarshal(body, &auth)
				if err == nil {
					r.Authentication = &auth
					r.LoggedIn = true

					log.Printf("New Expire Date %s", auth.Expires)
				}
			}
		}
	}

	if err != nil {
		log.Printf("Failure to update token : %s", err.Error())
		r.LoggedIn = false
	}
}

func (r *RestAPI) Init() error {
	r.APIUpdate.Lock()
	defer r.APIUpdate.Unlock()

	r.StopUpdate = make(chan bool)

	var err error
	r.LoggedIn = false
	//ignoring self signed certs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	r.Client = &http.Client{Transport: tr}

	var resp *http.Response
	resp, err = r.Client.PostForm(r.CreateUrl("token"), url.Values{"grant_type": {"password"}, "username": {r.Username}, "password": {r.Password}})

	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err == nil {
			if resp.StatusCode < 299 {
				auth := AuthenticationHeader{}
				err = json.Unmarshal(body, &auth)
				if err == nil {
					r.Authentication = &auth
					r.LoggedIn = true

					//still need to do clean up
					go func() {
						sec := (auth.ExpiresIn / 3)
						//log.Printf("Updating auth header in %d secs", sec)
						tick := time.Tick(time.Duration(sec) * time.Second)
						for keepupdating := true; keepupdating; {
							select {
							case <-r.StopUpdate:
								//log.Printf("Stopping ...")
								keepupdating = false
							case <-tick:
								r.Update()
							}
						}
					}()
				}
			} else {
				err = fmt.Errorf("Could not auth, got return code %d", resp.StatusCode)
			}
		}
	}
	return err
}

func (r *RestAPI) StartRestoreAction(orgId string, time string) (error, RestoreSession) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	var restoreSession RestoreSession

	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		jsonreq := fmt.Sprintf(`{  "explore":  {"datetime": "%s"  } }`, time)
		req, err = r.CreateRequest(fmt.Sprintf("Organizations/%s/action", orgId), "Post", jsonreq)
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode < 299 {
					var body []byte
					body, err = ioutil.ReadAll(resp.Body)
					if err == nil {
						restoreSession = RestoreSession{}
						err = json.Unmarshal(body, &restoreSession)
					}
				} else {
					err = fmt.Errorf("Got bad response code %d", resp.StatusCode)
				}
			}
		}
	}
	return err, restoreSession
}
func (r *RestAPI) StopRestoreSession(sessionId string) error {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		jsonreq := `{ "stop": null }`
		// time.Sleep(5 * time.Second)
		req, err = r.CreateRequest(fmt.Sprintf("RestoreSessions/%s/Action", sessionId), "Post", jsonreq)
		// uri := req.URL
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode > 299 {
					err = fmt.Errorf("Got bad response code %d", resp.StatusCode)
					// fmt.Println(err)
				} else {
					// var body []byte
					// body, _ = ioutil.ReadAll(resp.Body)
					// fmt.Printf("\nStopped %s code %d (%s)",sessionId,resp.StatusCode,body)
					// fmt.Printf("\nRequest %s",uri.String())
					// fmt.Printf("\nRequest body %s",jsonreq)
				}
			} else {
				// fmt.Printf("not stopped resp")
				// fmt.Println(err)
			}
		} else {
			// fmt.Printf("not stopped req")
			// fmt.Println(err)
		}
	}
	return err
}
func (r *RestAPI) GetRestoreMailboxesPage(sessionId string, offset int, limit int) (error, RestoreMailboxPage) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	var page RestoreMailboxPage
	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		req, err = r.CreateRequest(fmt.Sprintf("/RestoreSessions/%s/Organization/Mailboxes?offset=%d&limit=%d", sessionId, offset, limit), "Get", "")
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode < 299 {
					var body []byte
					body, err = ioutil.ReadAll(resp.Body)
					if err == nil {
						page = RestoreMailboxPage{}
						err = json.Unmarshal(body, &page)
					}
				} else {
					err = fmt.Errorf("Got bad response code %d", resp.StatusCode)
				}
			}
		}
	}
	return err, page
}
func (r *RestAPI) GetRestoreMailboxes(sessionId string) (error, RestoreMailboxPage) {
	return r.GetRestoreMailboxesPage(sessionId, 0, 5000)
}
func (r *RestAPI) GetRestoreMailboxFoldersPage(sessionId string, mailboxId string, offset int, limit int) (error, RestoreMailboxFolderPage) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	var page RestoreMailboxFolderPage
	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		req, err = r.CreateRequest(fmt.Sprintf("/RestoreSessions/%s/Organization/Mailboxes/%s/Folders?offset=%d&limit=%d", sessionId, mailboxId, offset, limit), "Get", "")
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode < 299 {
					var body []byte
					body, err = ioutil.ReadAll(resp.Body)
					if err == nil {
						page = RestoreMailboxFolderPage{}
						err = json.Unmarshal(body, &page)
					}
				} else {
					err = fmt.Errorf("Got bad response code %d", resp.StatusCode)
				}
			}
		}
	}
	return err, page
}

func (r *RestAPI) DoRestMailboxItem(sessionId string, mailboxId string, mailboxItemId string, mailboxusername string, mailboxpassword string) (error, RestoreResponse) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	//var feedback string
	var result RestoreResponse

	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		jsonreq := fmt.Sprintf(`{ "RestoreToOriginalLocation": {
			"userName": "%s",
			"userPassword": "%s",
			"ChangedItems": "True",
			"DeletedItems": "True",
			"MarkRestoredAsUnread": "True"		  
		}}`, mailboxusername, mailboxpassword)
		suburi := fmt.Sprintf("/RestoreSessions/%s/Organization/Mailboxes/%s/Items/%s/Action", sessionId, mailboxId, mailboxItemId)
		req, err = r.CreateRequest(suburi, "POST", jsonreq)
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode > 299 {
					var body []byte
					body, _ = ioutil.ReadAll(resp.Body)

					err = fmt.Errorf("Got bad response code %d : %s (%s)", resp.StatusCode, body, suburi)
				} else {
					var body []byte
					body, _ = ioutil.ReadAll(resp.Body)

					err = json.Unmarshal(body, &result)

					//feedback = fmt.Sprintf("%s", body)
				}
			}
		}
	}
	return err, result
}

func (r *RestAPI) GetRestoreMailboxFolders(sessionId string, mailboxId string) (error, RestoreMailboxFolderPage) {
	return r.GetRestoreMailboxFoldersPage(sessionId, mailboxId, 0, 5000)
}

func (r *RestAPI) GetRestoreMailboxItemsPage(sessionId string, mailboxId string, offset int, limit int) (error, RestoreMailboxItemPage) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	var page RestoreMailboxItemPage
	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		req, err = r.CreateRequest(fmt.Sprintf("/RestoreSessions/%s/Organization/Mailboxes/%s/Items?offset=%d&limit=%d", sessionId, mailboxId, offset, limit), "Get", "")
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode < 299 {
					var body []byte
					body, err = ioutil.ReadAll(resp.Body)
					if err == nil {
						page = RestoreMailboxItemPage{}
						err = json.Unmarshal(body, &page)

						// fmt.Printf("%s",body)
					}
				} else {
					err = fmt.Errorf("Got bad response code %d", resp.StatusCode)
				}
			}
		}
	}
	return err, page
}

func (r *RestAPI) GetRestoreMailboxItems(sessionId string, mailboxId string) (error, RestoreMailboxItemPage) {
	return r.GetRestoreMailboxItemsPage(sessionId, mailboxId, 0, 5000)
}

func (r *RestAPI) GetRestoreMailboxFolderItemsPage(sessionId string, mailboxId string, folderId string, offset int, limit int) (error, RestoreMailboxItemPage) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	var page RestoreMailboxItemPage
	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request

		folderQuery := ""
		if folderId != "" {
			folderQuery = fmt.Sprintf("&folderId=%s", folderId)
		}

		req, err = r.CreateRequest(fmt.Sprintf("/RestoreSessions/%s/Organization/Mailboxes/%s/Items?offset=%d&limit=%d%s", sessionId, mailboxId, offset, limit, folderQuery), "Get", "")
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode < 299 {
					var body []byte
					body, err = ioutil.ReadAll(resp.Body)
					if err == nil {
						page = RestoreMailboxItemPage{}
						err = json.Unmarshal(body, &page)

						//fmt.Printf("%s", body)
					}
				} else {
					err = fmt.Errorf("Got bad response code %d", resp.StatusCode)
				}
			}
		}
	}
	return err, page
}

func (r *RestAPI) GetRestoreMailboxFolderItems(sessionId string, mailboxId string, folderId string) (error, RestoreMailboxItemPage) {
	return r.GetRestoreMailboxFolderItemsPage(sessionId, mailboxId, folderId, 0, 5000)
}

func (r *RestAPI) GetOrganizations() (error, []Organization) {
	r.APIUpdate.RLock()
	defer r.APIUpdate.RUnlock()

	var err error
	var orgs []Organization

	if r.Client != nil && r.LoggedIn {
		var resp *http.Response
		var req *http.Request
		req, err = r.CreateRequest("Organizations", "Get", "")
		if err == nil {
			resp, err = r.Client.Do(req)
			if err == nil {
				if resp.StatusCode < 299 {
					var body []byte
					body, err = ioutil.ReadAll(resp.Body)
					if err == nil {
						orgs = []Organization{}
						err = json.Unmarshal(body, &orgs)
					}
				} else {
					err = fmt.Errorf("Bad return code %d", resp.StatusCode)
				}
			}
		}
	} else {
		err = fmt.Errorf("Not Logged In")
	}
	return err, orgs
}
func (r *RestAPI) Close() {
	if r.Authentication != nil {
		log.Printf("Stopping Auth API Update Loop")
		r.StopUpdate <- true
	}
}
