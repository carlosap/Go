package controllers
import (
	"net/http"
	"encoding/json"
	"sort"
	"github.com/Go/server/models"
	"github.com/Go/server/dbContext"
	"github.com/Go/server/util/logging"
	"github.com/martini-contrib/encoder"
	"github.com/go-martini/martini"
	
) 
func registerPriorityEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/priorities", func(r martini.Router) {
		m.Get("", GetPrioritiesHandler)
		m.Get("/:priorityID/files", GetFilesByPriorityHandler)
		m.Get("/:priorityID/batches", GetBatchesByPriorityHandler)
	})
}

func GetPrioritiesHandler(r *http.Request, w http.ResponseWriter, enc encoder.Encoder) (int, []byte) {
	var priorities = make([]models.Priority, 0)
	var records = make([]models.Record, 0)
	var err error
	records, err = dbContext.GetTargetData("priorityID")
	if err != nil {
		logging.Errorf("Unable to get priorities, Error: %+v.", err)
		return http.StatusNotFound, encoder.Must(enc.Encode(priorities))
	}
	priorities = convertRecordToPriorities(records)
	return http.StatusOK, encoder.Must(enc.Encode(priorities))
}
func convertRecordToPriorities(records []models.Record) []models.Priority {
	priorities := make([]models.Priority, len(records))
	for i, r := range records {
		p := models.Priority{}
		if err := json.Unmarshal([]byte(r.Target.Data), &p); err == nil {
			priorities[i] = p
		}
	}
	sort.Sort(models.Priorities(priorities))
	return priorities
}