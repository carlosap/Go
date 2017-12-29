package controllers

import (
	"sort"
	"encoding/json"
	"github.com/Go/server/models"
)

//ConvertRecordToInventories deserializes target data (db) and
//converts records into inventories
func ConvertRecordToInventories(records []models.Record) []models.Inventory {
	inventories := make([]models.Inventory, len(records))
	for i, r := range records {
		inventoryItem := models.Inventory{}
		if err := json.Unmarshal([]byte(r.Target.Data), &inventoryItem); err == nil {
			inventories[i] = inventoryItem
		}
	}
	sort.Sort(models.SortByInventoryID{inventories})
	return inventories
}