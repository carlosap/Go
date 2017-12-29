package models

//Inventory Type
type Inventory struct {
	Objective    string         `json:"objective"`
	BatchNumber  string         `json:"batchNumber"`
	ID           string         `json:"id"`
	Description  string         `json:"description"`
	DeviceType   string         `json:"deviceType"`
	DeviceName   string         `json:"deviceName"`
	SerialNumber string         `json:"serialNumber"`
	Status       string         `json:"status"`
	User         string         `json:"user"`
	Comments     []Note         `json:"comments, omitempty" cql:"comments"`
	SourceDB     string         `json:"sourceDb"`
	DataSecurity
}

//Inventories Type
type Inventories []Inventory

//SortByInventoryID for sorting InventoryID slice items
type SortByInventoryID struct {
	Inventories
}

func (inventories Inventories) Len() int {
	return len(inventories)
}

func (inventories Inventories) Swap(i, j int) {
	inventories[i], inventories[j] = inventories[j], inventories[i]
}

func (s SortByInventoryID) Less(i, j int) bool {
	return s.Inventories[i].ID < s.Inventories[j].ID
}