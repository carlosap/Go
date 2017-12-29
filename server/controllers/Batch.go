package controllers
import (
	"sync"
	"runtime"
	"fmt"
	"strings"
	"errors"
	"net/http"
	"log"
	"encoding/json"
	"sort"
	"github.com/Go/server/models"
	"github.com/Go/server/dbContext"
	"github.com/Go/server/util/logging"
	"github.com/martini-contrib/encoder"
	"github.com/go-martini/martini"
	
) 

func RegisterBatchEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/batches", func(r martini.Router) {
		m.Get("/:batchID/files", GetFilesByBatchHandler)
		m.Get("", GetBatchesHandler)
	})
}

func GetBatchesHandler(r *http.Request, w http.ResponseWriter,enc encoder.Encoder) (int, []byte) {
	var batches = make([]models.Batch, 0)
	var records = make([]models.Record, 0)
	steps := []string{"initial", "triage", "report", "translate", "analysis", "exploit"}
	for _, s := range steps {
		records = make([]models.Record, 0)
		var err error
		records, err = dbContext.GetTargetDataByLinkType(s, "batches")
		if err != nil {
			log.Printf("Unable to get batches by step %+v.", err)
		}
		batches = append(batches, ConvertRecordToBatches(records)...)
	}
	
	if len(batches) > 0 {
		addMediaSelectorCount(batches)
	}

	return http.StatusOK, encoder.Must(enc.Encode(batches))
}
func GetBatchesByPriorityHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var batches = make([]models.Batch, 0)
	var records = make([]models.Record, 0)
	var err error
	priorityID, ok := params["priorityID"]
	if !ok || len(priorityID) == 0 {
		return http.StatusBadRequest, encoder.Must(enc.Encode(batches))
	}
	records, err = dbContext.GetTargetDataByLinkType(priorityID, "batches")
	if err != nil {
		logging.Errorf("Unable to get batches by priority, Error: %+v.", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(batches))
	}
	batches = ConvertRecordToBatches(records)
	return http.StatusOK, encoder.Must(enc.Encode(batches))
}

func GetBatchesByStepHandler(r *http.Request, w http.ResponseWriter, params martini.Params,enc encoder.Encoder) (int, []byte) {

	var batches = make([]models.Batch, 0)
	var records = make([]models.Record, 0)
	var err error
	stepID, ok := params["stepID"]
	if !ok || len(stepID) == 0 {
		logging.Errorf("Step ID can not be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(batches))
	}
	records, err = dbContext.GetTargetDataByLinkType(stepID, "batches")
	if err != nil {
		logging.Errorf("Unable to get batches by step, Error: %+v.", err)
		return http.StatusBadRequest, encoder.Must(enc.Encode(batches))
	}
	batches = ConvertRecordToBatches(records)
	if len(batches) > 0 {
		batches = addMediaSelectorCount(batches)
	}
	return http.StatusOK, encoder.Must(enc.Encode(batches))
}

func addMediaSelectorCount(batches []models.Batch) []models.Batch {
	runtime.GOMAXPROCS(2) 
	var asyncAwait sync.WaitGroup
	for i := 0; i < len(batches); i++ {
		cbatch := batches[i]
		asyncAwait.Add(2)
		go func(){		
			defer asyncAwait.Done()
			var ct = GetMediaCount(cbatch.ID)
			for k, v := range ct {
				batches[i].Selectors = append(batches[i].Selectors, fmt.Sprintf("%s:%d",k,v))
			}	
		}()
		go func(){		
			defer asyncAwait.Done()
			var ct = GetSelectorsCount(cbatch.ID)
			for j, p := range ct {
				batches[i].Selectors = append(batches[i].Selectors, fmt.Sprintf("%s:%d",j,p))
			}	
		}()
		asyncAwait.Wait()
	}
	return batches
}

//getBatchMediaCountHandler: audio, document, geo, image, report, translateddoc, video
//end point: http://localhost:5001/batches/134734-2AEP/media/count
//func GetBatchMediaCountHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
func GetMediaCount(batchID string) (map[string]int) {
	var files = make([]models.File, 0)
	t := make(map[string]int)
	files, _ = getFilesByBatchID(batchID)
	for _, fileItem := range files {	
			getCountByCategory(fileItem.Type, &t)					
	}
	return t
}

//getBatchSelectorsCountHandler: email, burners (cellular), social media, macaddress
//http://localhost:5001/batches/134734-2AEP/selectors/count
//http://localhost:5001/batches/175336-1DJY/selectors/count
func GetSelectorsCount(batchID string) (map[string]int) {
	var selectors = make([]models.Selector, 0)
	//m := map[string]interface{}{}
	t := make(map[string]int)
	selectors, _ = getSelectorsByBatchID(batchID)
	for _, selectorItem := range selectors {	
			getCountByCategory(selectorItem.Type, &t)					
	}
	// m["batchid"] = map[string]interface{}{
    // 	batchID: &t,
	// }
	return t
}

func getCountByCategory(fileType string, t *map[string]int) (map[string]int, error){
	if len(fileType) <= 0 {
		return *t, errors.New("Error: File Type field is missing")
	}
	fileType = strings.ToLower(fileType)
	(*t)[fileType] = (*t)[fileType] + 1
	return *t, nil
}



func getSelectorsByBatchID(batchID string) ([]models.Selector, error) {
	var selectors = make([]models.Selector, 0)
	var records = make([]models.Record, 0)
	if len(batchID) <= 0 {
		return selectors, errors.New("Error: BatchId is require")
	}
	
	records, err := dbContext.GetTargetDataByLinkType(batchID, "selectors")
	if err != nil {
		return selectors, err
	}
	
	selectors = ConvertRecordToSelectors(records)
	return selectors, nil
}

func getFilesByBatchID(batchID string) ([]models.File, error) {
	var files = make([]models.File, 0)
	var records = make([]models.Record, 0)
	if len(batchID) <= 0 {
		return files, errors.New("Error: BatchId is require")
	}
	
	records, err := dbContext.GetTargetDataByLinkType(batchID, "files")
	if err != nil {
		return files, err
	}
	
	files = ConvertRecordToFiles(records)
	return files, nil
}

//ModifyFilesRemaining Retrieves the batch, modifies files remaining
// updates all references to that batch
func ModifyFilesRemaining(batchID string, modifier int) error {
	// Retrieve the static columns
	bData, err := dbContext.GetStaticCols(batchID)
	if err != nil {
		return err
	}
	// Update the files remaining
	var newBatch models.Batch
	err = json.Unmarshal([]byte(bData.Data), &newBatch)
	if err != nil {
		return err
	}
	newBatch.FilesRemaining += modifier
	tmp, err := json.Marshal(newBatch)
	bData.Data = string(tmp)
	err = dbContext.UpdateRecordsByID(bData,"batches")
	if err != nil {
		return err
	}
	return err
}

func MoveBatchToNewStep(batchID string, newStepID string) error {
	var batch models.Batch
	// Grab the batch
	bData, err := dbContext.GetStaticCols(batchID)
	if err != nil {
		logging.Errorf("Unable to get batch %s: %+v", batchID, err)
		return err
	}

	err = json.Unmarshal([]byte(bData.Data), &batch)
	if err != nil {
		logging.Errorf("Unable to unmarshal batch %s: %+v", batchID, err)
		return err
	}
	// Update it's step
	oldStepID := batch.Step
	batch.Step = newStepID
	tmp, err := json.Marshal(batch)
	bData.Data = string(tmp)
	// Update all references to its step
	//BatchLinkType = "batches"
	err = dbContext.UpdateRecordsByID(bData, "batches")
	if err != nil {
		logging.Errorf("Unable to update batch %s: %+v", batchID, err)
		return err
	}
	// Delete oldStep->batch
	err = dbContext.DeleteData(oldStepID, "batches", batchID)
	if err != nil {
		logging.Errorf("Unable to delete src=%s linktype=%s target=%s: %+v",oldStepID, "batches", batchID, err)
		return err
	}

	// Create newStep->batch
	step := models.Step{
		ID:    newStepID,
		Order: 2,
	}

	sData, err := CreateDataFromStep(step, bData.DataSecurity)
	if err != nil {
		logging.Errorf("Unable to create step data: %+v", err)
		return err
	}

	rec := models.Record{
		Source:   sData,
		LinkType: "batches",
		Target:   bData,
	}
	_, err = dbContext.NewTargetData(rec)
	if err != nil {
		logging.Errorf("Error %+v creating src=%s link=%s target=%s",err, newStepID, "batches", batchID)
		return err
	}
	
	err = dbContext.DeleteData(batchID, "steps", oldStepID)
	if err != nil {
		logging.Errorf("Unable to delete src=%s linktype=%s target=%s: %+v",batchID, "steps", oldStepID, err)
		return err
	}
	// Create batch->newStep
	rec = models.Record{
		Source:   bData,
		LinkType: "steps",
		Target:   sData,
	}

	_, err = dbContext.NewTargetData(rec)
	if err != nil {
		logging.Errorf("Error %+v creating src=%s link=%s target=%s",err, batchID, "steps", newStepID)
		return err
	}

	return err
}

func createDataFromBatch(batch models.Batch) (models.Data, error) {
	var bData models.Data
	bString, err := json.Marshal(batch)

	if err != nil {
		return bData, err
	}

	bData = models.Data{
		ID:           batch.ID,
		Data:         string(bString),
		DataSecurity: batch.DataSecurity,
	}

	return bData, err
}

func ConvertRecordToBatches(records []models.Record) []models.Batch {
	batches := make([]models.Batch, len(records))
	for i, r := range records {
		b := models.Batch{}
		if err := json.Unmarshal([]byte(r.Target.Data), &b); err == nil {
			batches[i] = b
		}
	}
	sort.Sort(models.Batches(batches))
	return batches
}
