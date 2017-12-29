package models

//Selector types
type Selector struct {
	Objective   string         `json:"objective"`
	ID          string         `json:"selectorID"`
	BatchNumber string         `json:"batchNumber"`
	InventoryID string         `json:"inventoryID"`
	Type        string         `json:"selectorType"`
	Value       string         `json:"selectorValue"`
	MD5Hash     string         `json:"md5Hash"`
	Comments    []string       `json:"gist"`
	MGRS        string         `json:"capturePoint"`
	SourceDB    string         `json:"sourceDb"`
	User        string 			`json:"user"`
	DataSecurity
}

//Selectors Type
type Selectors []Selector