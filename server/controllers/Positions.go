package controllers
import (
	"encoding/json"
	"github.com/Go/server/models"
	"github.com/Go/server/dbContext"
	"github.com/Go/server/util/logging"

) 

func GetFileHighlightsByID(fileID string) ([]models.Position, error) {
	var positions = make([]models.Position, 0)
	var records = make([]models.Record, 0)
	var err error


	records, err = dbContext.GetRecordsBySource(fileID)

	if err != nil {
		 logging.Errorf("Unable to retrieve positions, Error: %+v.", err)
		 return positions, err
	}
	
	if len(records) > 0 {

		for _, r := range records {					
				if r.Source.Data != "" && r.LinkType == models.PositionsLinkType {

					err := json.Unmarshal([]byte(r.Source.Data), &positions)
					if err != nil {
						logging.Errorf("Error %+v Unmarshal data from %s", err, fileID)
					}

				}
		}

	}

	return positions, nil
}


func UpdateFileHighlightPositions(fileID string, positions models.Positions) error {
	strData, err := json.Marshal(positions)

	if err != nil {
		logging.Errorf("Unable to marshal updated file data: %+v", err)
		return err
	}

	pData := models.Data{
		ID:           fileID,
		Data:         string(strData),
	}

	updatedRecord := models.Record{
		LinkType: models.PositionsLinkType,
		Source: pData,
	}

	_, err = dbContext.UpdateTargetData(updatedRecord)

	if err != nil {
		logging.Errorf("Unable to update file %s static columns: %+v", fileID, err)
	}

	return err
}
