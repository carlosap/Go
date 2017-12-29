package models

//Data describes a specific piece of data.
type Data struct {
	ID   string `json:"id"`
	Data string `json:"data"`
	DataSecurity
}