package controllers
import (
	"strings"
	"reflect"
	"encoding/json"
	"net/http"
	
	"github.com/Go/server/util/logging"
	"github.com/Go/server/dbContext"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"github.com/Go/server/models"

) 


func RegisterFileEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/files", func(r martini.Router) {
		m.Post("/:fileID/dismiss", DismissFileHandler)
		m.Put("/:fileID/update", UpdateFileHandler)
		m.Post("/:fileID/move/:stepID", MoveFileHandler)
		m.Get("/:fileID/highlights", GetFileHighlightsHandler)
		m.Put("/:fileID/highlights", UpdateFileHighlightsHandler)
	})
}

func UpdateFileHighlightsHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var retVal models.Positions

	fileID, ok := params["fileID"]
	if !ok || len(fileID) == 0 {
		logging.Errorf("File ID cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}

	err := json.NewDecoder(r.Body).Decode(&retVal)
	if err != nil {
		logging.Errorf("Error decoding updated positions:  %+v", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}

	err = UpdateFileHighlightPositions(fileID,retVal)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	
	return http.StatusOK, encoder.Must(enc.Encode(retVal))
}

func GetFileHighlightsHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var positions = make([]models.Position, 0)
	var err error

	fileID, ok := params["fileID"]
	if !ok || len(fileID) == 0 {
		logging.Errorf("File ID cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(positions))
	}

	positions, err = GetFileHighlightsByID(fileID)
	if err != nil {
		logging.Errorf("Error: GetFileHighlightsByID %-v", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(positions))
	}
	return http.StatusOK, encoder.Must(enc.Encode(positions))
}


func GetFilesByBatchHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var files = make([]models.File, 0)
	var records = make([]models.Record, 0)
	var err error
	batchID, ok := params["batchID"]
	if !ok || len(batchID) == 0 {
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	records, err = dbContext.GetTargetDataByLinkType(batchID, "files")
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	files = ConvertRecordToFiles(records)
	
	return http.StatusOK, encoder.Must(enc.Encode(files))

}

func GetFilesByPriorityHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var files = make([]models.File, 0)
	var records = make([]models.Record, 0)
	var err error
	priorityID, ok := params["priorityID"]
	if !ok || len(priorityID) == 0 {
		logging.Errorf("Priority ID can not be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	records, err = dbContext.GetTargetDataByLinkType(priorityID, "files")
	if err != nil {
		logging.Errorf("Unable to get files by priority, Error: %+v.", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	files = ConvertRecordToFiles(records)
	return http.StatusOK, encoder.Must(enc.Encode(files))
}

func GetFilesByStepHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var files = make([]models.File, 0)
	var records = make([]models.Record, 0)
	var err error

	//fix me "translate not showing"
	stepID, ok := params["stepID"]
	if !ok || len(stepID) == 0 {
		logging.Errorf("Step ID can not be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	records, err = dbContext.GetTargetDataByLinkType("translate", "files")
	if err != nil {
		logging.Errorf("Unable to get files by step, Error: %+v.", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	files = ConvertRecordToFiles(records)
	return http.StatusOK, encoder.Must(enc.Encode(files))
}


func DismissFileHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var retVal models.File
	var records = make([]models.Record, 0)
	//fileIDParam       = "fileID"
	fileID, ok := params["fileID"]
	if !ok || len(fileID) == 0 {
		logging.Errorf("File ID cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}

	// update the static columns
	fData, err := dbContext.GetStaticCols(fileID)

	if err != nil {
		logging.Errorf("Unable to get the static columns for file %s: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}

	tmpFile := models.File{}
	err = json.Unmarshal([]byte(fData.Data), &tmpFile)
	if err != nil {
		logging.Errorf("Unable to unmarshal file data: %+v", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	tmpFile.Dismissed = true
	err = updateFileStaticCols(tmpFile)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}

	// delete fileId | batches | batchId
	//BatchLinkType = "batches"
	records, err = dbContext.GetTargetDataByLinkType(fileID, "batches")
	if err != nil || len(records) != 1 {
		logging.Errorf("Unable to get the batch for file %s, Error: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	batchID := records[0].Target.ID
	err = dbContext.DeleteData(fileID, "batches", batchID)
	if err != nil {
		logging.Errorf("Error %+v while trying to delete %s | %s | %s",err, fileID, "batches", batchID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	// delete batchId | files | fileId
	err = dbContext.DeleteData(batchID, "files", fileID)
	if err != nil {
		logging.Errorf("Error %+v while trying to delete %s | %s | %s",err, batchID, "files", fileID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	// subtract from batchID files remaining
	err = ModifyFilesRemaining(batchID, -1)
	if err != nil {
		logging.Errorf("Error %+v while trying to modify files remaining in batch %s",err, batchID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	records, err = dbContext.GetTargetDataByLinkType(fileID,"steps")
	if err != nil || len(records) != 1 {
		logging.Errorf("Unable to get the step for file %s, Error: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	stepID := records[0].Target.ID
	err = dbContext.DeleteData(stepID,"files", fileID)
	err = dbContext.DeleteData(fileID, "steps", stepID)
	disData := models.Data{
		ID:           "dismissed",
		Data:         "",
		DataSecurity: fData.DataSecurity,
	}
	disRecord := models.Record{
		Source:   disData,
		LinkType: "dismissed",
		Target:   fData,
	}
	_, err = dbContext.NewTargetData(disRecord)
	return http.StatusOK, encoder.Must(enc.Encode(retVal))
}

///files/014846df-bc14-8d3a-3fae-d6fcc04d44f8/update
func UpdateFileHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var languages = make([]string, 0)		
	var priorities = make([]string, 0)		
	var selectors = make([]string, 0)
	var retVal models.File
	fileID, ok := params["fileID"]
	if !ok || len(fileID) == 0 {
		logging.Errorf("File ID cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	err := json.NewDecoder(r.Body).Decode(&retVal)
	if err != nil {
		logging.Errorf("Error decoding updated file: %s: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	
	err = updateFileStaticCols(retVal)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
	
	fData, err := CreateDataFromFile(retVal)
	if err != nil {
		logging.Errorf("Error %+v creating data from %s", err, fileID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}

	err = dbContext.UpdateRecordsByID(fData, "files")
	if err != nil {
		logging.Errorf("Error %+v getting records by source: %s", err, fileID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
	}
						
	languages, priorities = getFilePriorityLanguage(retVal.BatchID,"files")
	sRecords, _ := dbContext.GetTargetDataByLinkType(retVal.ID,"batches")
	if len(sRecords) > 0 {
		for _, r := range sRecords {					
				if r.Target.Data != "" {
					var batch models.Batch
					err := json.Unmarshal([]byte(r.Target.Data), &batch)
					if err != nil {
						logging.Errorf("Error %+v Unmarshal data from %s", err, batch.ID)
						return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
					}
					if batch.ID != "" {
						if !reflect.DeepEqual(batch.Priority,languages) || 
								!reflect.DeepEqual(batch.Languages,priorities) {
									
									for _, rtvalLanguage :=  range retVal.Languages {									
										isFound :=  IsInSlice(rtvalLanguage, languages)
										if !isFound && rtvalLanguage != "" {
											languages = append(languages, rtvalLanguage)	
										}
									}			
															
									for _, rtvalPriority :=  range retVal.Priority {									
										isFound :=  IsInSlice(rtvalPriority, priorities)
										if !isFound && rtvalPriority != "" {
											priorities = append(priorities, rtvalPriority)	
										}
									}	
									//selectors
									if len(retVal.Selectors) > 0 {
										for _, rtvalSelector := range retVal.Selectors {
											isFound := IsInSlice(rtvalSelector, selectors)
											if !isFound && rtvalSelector != "" {
												v := strings.ToLower(fileID + "|" + rtvalSelector)
												selectors = append(selectors, v)
											}
										}	
										batch.Selectors = selectors	
									}																	
									batch.Languages = languages
									batch.Priority = priorities		
									
									bData, err := createDataFromBatch(batch)
									if err != nil {
										logging.Errorf("Error %+v creating data from %s", err, batch.ID)
										return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
									}		
									logging.Warn("Batch Data: %+v", bData)
									// err = dbContext.UpdateRecordsByID(bData, "batches")
									// if err != nil {
									// 	logging.Errorf("Error %+v UpdateRecordsByID: %s", err, bData.ID)
									// 	return http.StatusBadRequest, encoder.Must(enc.Encode(retVal))
									// }								
						}		
					}
				}
		}
	}
	return http.StatusOK,encoder.Must(enc.Encode(retVal))
}

func ConvertRecordToFiles(records []models.Record) []models.File {
	files := make([]models.File, len(records))
	for i, r := range records {
		f := models.File{}
		if err := json.Unmarshal([]byte(r.Target.Data), &f); err == nil {
			files[i] = f
		}
	}
	return files
}

//getFilePriorityLanguage discovers language and priorities' 
//associations by providing the sourceid annd the linktype. 
func getFilePriorityLanguage(sourceID, linkType string) ([]string, []string) {
	var languages = make([]string, 0)		
	var priorities = make([]string, 0)	
	fRecords, _ := dbContext.GetTargetDataByLinkType(sourceID,linkType)
	if len(fRecords) > 0 {
		for _, f := range fRecords {
			if f.Target.Data != "" {
				var file models.File
				err := json.Unmarshal([]byte(f.Target.Data), &file)
				if err != nil {
					logging.Errorf("Error %+v Unmarshal data from %s", err, sourceID)
					return languages, priorities
				}							
				//-------priorities----------------------
				for _, pValue :=  range file.Priority {									
					isFound :=  IsInSlice(pValue, priorities)
					if !isFound && pValue != "" {
						priorities = append(priorities, pValue)	
					}
				}
				//-------languages--------------------------
				for _, lValue := range file.Languages {
					isFound := IsInSlice(lValue, languages)
					if !isFound && lValue != "" {
						languages = append(languages, lValue)
					}
					
				}
			}
		}
	}
	return languages, priorities
}

//isInSlice returns boolean value to check and see
//if the strValue is part the string list/array. Recommened 
//to use this only when dealing with small collections
func IsInSlice(strValue string, list []string) bool {
 	for _, v := range list {
 		if v == strValue {
 			return true
 		}
 	}
 	return false
 }

 func updateFileStaticCols(file models.File) error {
	tmpData, err := json.Marshal(file)
	if err != nil {
		logging.Errorf("Unable to marshal updated file data: %+v", err)
		return err
	}
	fData := models.Data{
		ID:           file.ID,
		Data:         string(tmpData),
		DataSecurity: file.DataSecurity,
	}
	updatedRecord := models.Record{
		Source: fData,
	}
	_, err = dbContext.UpdateTargetData(updatedRecord)
	if err != nil {
		logging.Errorf("Unable to update file %s static columns: %+v", file.ID, err)
	}
	return err
}

func CreateDataFromFile(tFile models.File) (models.Data, error) {

	var fData models.Data
	fString, err := json.Marshal(tFile)
	if err != nil {
		return fData, err
	}
	fData = models.Data{
		ID:           tFile.ID,
		Data:         string(fString),
		DataSecurity: tFile.DataSecurity,
	}
	return fData, err
}

func MoveFileHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var file models.File
	//fileIDParam       = "fileID"
	fileID, ok := params["fileID"]
	if !ok || len(fileID) == 0 {
		logging.Errorf("File ID cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}
	//stepIDParam   = "stepID"
	stepID, ok := params["stepID"]
	if !ok || len(fileID) == 0 {
		logging.Errorf("Step ID cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	// Get static columns
	fData, err := dbContext.GetStaticCols(fileID)
	if err != nil {
		logging.Errorf("Unable to get file %s static columns: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	err = json.Unmarshal([]byte(fData.Data), &file)
	if err != nil {
		logging.Errorf("Unable to unmarshal file data: %+v", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}
	
	oldStepID := file.Step
	file.Step = stepID
	fData, err = CreateDataFromFile(file)
	if err != nil {
		logging.Errorf("Unable to create file data: %+v", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	// Update with the new step
	err = updateFileStaticCols(file)
	if err != nil {
		logging.Errorf("Unable to update file %s static columns: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	// Delete src=fileID linktype=step target=oldStepID
	//StepLinkType  = "steps"
	err = dbContext.DeleteData(fileID, "steps", oldStepID)
	if err != nil {
		logging.Errorf("Error %+v deleting src=%s link=%s target=%s",err, fileID, "steps", oldStepID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	// Create src=fileID linktype=step target=newStepID
	step := models.Step{
		ID:    stepID,
		Order: 2,
	}

	sData, err := CreateDataFromStep(step, fData.DataSecurity)
	if err != nil {
		logging.Errorf("Unable to create step data: %+v", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	//StepLinkType  = "steps"
	rec := models.Record{
		Source:   fData,
		LinkType: "steps",
		Target:   sData,
	}

	_, err = dbContext.NewTargetData(rec)
	if err != nil {
		logging.Errorf("Error %+v creating src=%s link=%s target=%s",err, fileID, "steps", stepID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}
	// Delete src=oldStepID linktype=files target=fileID
	err = dbContext.DeleteData(oldStepID, "files", fileID)
	if err != nil {
		logging.Errorf("Error %+v deleting src=%s link=%s target=%s",err, oldStepID, "files", fileID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}

	// Create src=newStepID linktype=files target=fileID
	rec = models.Record{
		Source:   sData,
		LinkType: "files",
		Target:   fData,
	}
	_, err = dbContext.NewTargetData(rec)
	if err != nil {
		logging.Errorf("Error %+v creating src=%s link=%s target=%s",err, stepID, "files", fileID)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}
	// Update all other references
	err = dbContext.UpdateRecordsByID(fData, "files")
	if err != nil {
		logging.Errorf("Unable to update file %s records: %+v", fileID, err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(file))
	}
	// If the original step is triage
	if oldStepID == "triage" {
		// subtract one from its batch's files remaining
		err = ModifyFilesRemaining(file.BatchID, -1)
		if err != nil {
			logging.Errorf("Error %+v while trying to modify files remaining in batch %s",err, file.BatchID)
			return http.StatusBadRequest, encoder.Must(enc.Encode(file))
		}
		// Grab the batch
		bData, err := dbContext.GetStaticCols(file.BatchID)
		if err != nil {
			logging.Errorf("Unable to get batch %s for file %s: %+v",file.BatchID, fileID, err)
			return http.StatusBadRequest, encoder.Must(enc.Encode(file))
		}

		// if the files remaining = 0
		var batch models.Batch
		err = json.Unmarshal([]byte(bData.Data), &batch)

		if err != nil {
			logging.Errorf("Unable to unmarshal batch %s: %+v",file.BatchID, err)
			return http.StatusBadRequest, encoder.Must(enc.Encode(file))
		}

		if batch.FilesRemaining == 0 {
			// move it to the next step
			err = MoveBatchToNewStep(file.BatchID, "translate")
			if err != nil {
				logging.Errorf("Unable to move batch %s to translate: %+v",file.BatchID, err)
				return http.StatusBadRequest, encoder.Must(enc.Encode(file))
			}
		}
	}

	return http.StatusOK, encoder.Must(enc.Encode(file))
}
