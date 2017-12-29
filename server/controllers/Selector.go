package controllers

import (
	"encoding/json"
	"github.com/Go/server/models"
)

func ConvertRecordToSelectors(records []models.Record) []models.Selector {
	selectors := make([]models.Selector, len(records))
	for i, r := range records {
		s := models.Selector{}
		if err := json.Unmarshal([]byte(r.Target.Data), &s); err == nil {
			selectors[i] = s
		}
	}
	return selectors
}