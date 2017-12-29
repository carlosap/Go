package controllers
import (
	"encoding/json"
	"sort"
	"github.com/Go/server/models"
	"github.com/Go/server/dbContext"
) 

func convertRecordToActions(records []models.Record) models.Actions {
	actions := make(models.Actions, len(records))

	for i, r := range records {
		a := models.Action{}

		err := json.Unmarshal([]byte(r.Target.Data), &a)

		if err == nil {
			actions[i] = a
		}
	}

	sort.Sort(models.Actions(actions))
	return actions
}

func getActionsByFileID(fileID string) (models.Actions, error) {
	var actions = make(models.Actions, 0)
	var records = make([]models.Record, 0)
	var err error

	records, err = dbContext.GetTargetDataByLinkType(fileID, "actions",)

	if err != nil {
		return actions, err
	}

	actions = convertRecordToActions(records)

	return actions, nil
}


