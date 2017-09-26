package main

type GUIRestoreItemMail struct {
	Id      string
	Parent  string
	From    string
	To      string
	Subject string
}

type GUIRestoreItemAppointment struct {
	Id        string
	Parent    string
	Organizer string
	Location  string
	StartTime string
}

type GUIRestoreItemContact struct {
	Id      string
	Parent  string
	Display string
	Email   string
}

type GUIRestoreItemTask struct {
	Id     string
	Parent string
	Owner  string
}

type GUIRestoreFolder struct {
	Id   string
	Name string
}

type GUIRestoreItems struct {
	Sessionid     string
	Mails         []GUIRestoreItemMail
	Appointments  []GUIRestoreItemAppointment
	Contacts      []GUIRestoreItemContact
	Folders       []GUIRestoreFolder
	Tasks         []GUIRestoreItemTask
	ErrorString   string
	SuccessString string
}

type GUIRestoreResponse struct {
	Result            string `json:"result"`
	Message           string `json:"message"`
	CreatedItemsCount int    `json:"createdItemsCount"`
	MergedItemsCount  int    `json:"mergedItemsCount"`
	FailedItemsCount  int    `json:"failedItemsCount"`
	SkippedItemsCount int    `json:"skippedItemsCount"`
}
