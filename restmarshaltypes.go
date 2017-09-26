package main

type AuthenticationHeader struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Issued       string `json:".issued"`
	Expires      string `json:".expires"`
}
type Link struct {
	Href string `json:"href"`
}
type Action struct {
	Uri    string `json:"uri"`
	Method string `json:"method"`
}

type Organization struct {
	Type            string             `json:"type"`
	Username        string             `json:"username"`
	ServerName      string             `json:"serverName"`
	UseSSL          bool               `json:"useSSL"`
	Id              string             `json:"Id"`
	Name            string             `json:"Name"`
	IsBackedup      bool               `json:"isBackedUp"`
	FirstBackupTime string             `json:"firstBackuptime"`
	LastBackupTime  string             `json:"lastBackuptime"`
	Links           map[string]*Link   `json:"_links"`
	Actions         map[string]*Action `json:"_actions"`
}

type RestoreSession struct {
	Id           string             `json:"Id"`
	Type         string             `json:"type"`
	PointInTime  string             `json:"pointInTime"`
	CreationTime string             `json:"creationTime"`
	State        string             `json:"state"`
	Result       string             `json:"result"`
	Links        map[string]*Link   `json:"_links"`
	Actions      map[string]*Action `json:"_actions"`
}
type RestoreMailbox struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	IsArchive bool   `json:"isArchive"`
}
type RestoreMailboxPage struct {
	Offset    int                `json:"offset"`
	Limit     int                `json:"limit"`
	MailBoxes []*RestoreMailbox  `json:"results"`
	Links     map[string]*Link   `json:"_links"`
	Actions   map[string]*Action `json:"_actions"`
}

type RestoreMailboxFolder struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Description string             `json:"description"`
	Links       map[string]*Link   `json:"_links"`
	Actions     map[string]*Action `json:"_actions"`
}

type RestoreMailboxFolderPage struct {
	Offset  int                     `json:"offset"`
	Limit   int                     `json:"limit"`
	Folders []*RestoreMailboxFolder `json:"results"`
	Links   map[string]*Link        `json:"_links"`
	Actions map[string]*Action      `json:"_actions"`
}

type RestoreMailboxItem struct {
	Id              string             `json:"id"`
	ItemClass       string             `json:"itemClass"`
	From            string             `json:"from"`
	CC              string             `json:"cc"`
	BCC             string             `json:"bcc"`
	To              string             `json:"to"`
	Send            string             `json:"send"`
	Received        string             `json:"received"`
	Reminder        bool               `json:"reminder"`
	Subject         string             `json:"subject"`
	StartTime       string             `json:"startTime"`
	EndTime         string             `json:"endTime"`
	Organizer       string             `json:"organizer"`
	Location        string             `json:"location"`
	Attendees       string             `json:"attendees"`
	Recurring       bool               `json:"recurring"`
	Address         string             `json:"address"`
	BusinesPhone    string             `json:"businesPhone"`
	Company         string             `json:"company"`
	DisplayAs       string             `json:"displayAs"`
	Email           string             `json:"email"`
	Fax             string             `json:"fax"`
	FileAs          string             `json:"fileAs"`
	FullName        string             `json:"fullName"`
	HomePhone       string             `json:"homePhone"`
	ImAddress       string             `json:"imAddress"`
	JobTitle        string             `json:"jobTitle"`
	Mobile          string             `json:"mobile"`
	WebPage         string             `json:"webPage"`
	Status          string             `json:"status"`
	PercentComplete float64            `json:"percentComplete"`
	StartDate       string             `json:"startDate"`
	DueDate         string             `json:"dueDate"`
	Owner           string             `json:"owner"`
	Links           map[string]*Link   `json:"_links"`
	Actions         map[string]*Action `json:"_actions"`
}
type RestoreMailboxItemPage struct {
	Offset  int                   `json:"offset"`
	Limit   int                   `json:"limit"`
	Items   []*RestoreMailboxItem `json:"results"`
	Links   map[string]*Link      `json:"_links"`
	Actions map[string]*Action    `json:"_actions"`
}

type RestoreResponse struct {
	CreatedItemsCount int      `json:"createdItemsCount"`
	MergedItemsCount  int      `json:"mergedItemsCount"`
	FailedItemsCount  int      `json:"failedItemsCount"`
	SkippedItemsCount int      `json:"skippedItemsCount"`
	Exceptions        []string `json:"exceptions"`
}
