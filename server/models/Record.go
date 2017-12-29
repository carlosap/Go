package models

//Record specifies data in a specific state.
type Record struct {
	LinkType string `json:"linkType"`
	Source   Data   `json:"source"`
	Target   Data   `json:"target"`
}
