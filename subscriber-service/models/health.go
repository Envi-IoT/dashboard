package models

//device health model
type Health struct {
	Device    string `json:"device"`
	BootCount int    `json:"Bootcount"`
	Battery   int    `json:"Battery"`
}
