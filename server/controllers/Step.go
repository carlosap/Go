package controllers

import (
	"sort"
	"net/http"
	"encoding/json"
	"github.com/Go/server/util/logging"
	"github.com/martini-contrib/encoder"
	"github.com/Go/server/models"
	"github.com/Go/server/dbContext"
	"github.com/go-martini/martini"
)

func RegisterStepEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/steps", func(martini.Router) {
		m.Get("", GetStepsHandler)
		m.Get("/:stepID/files", GetFilesByStepHandler)
		m.Get("/:stepID/batches", GetBatchesByStepHandler)
	})
}

func GetStepsHandler(r *http.Request, w http.ResponseWriter,enc encoder.Encoder) (int, []byte) {
	var err error
	var records = make([]models.Record, 0)
	var steps = make([]models.Step, 0)
	records, err = dbContext.GetTargetData("step")
	if err != nil {
		logging.Errorf("Unable to get steps, Error: %+v.", err)
		return http.StatusNotFound, encoder.Must(enc.Encode(steps))
	}
	steps = convertRecordToSteps(records)
	return http.StatusOK, encoder.Must(enc.Encode(steps))
}

func convertRecordToSteps(records []models.Record) models.Steps {
	steps := make(models.Steps, len(records))
	for i, r := range records {
		s := models.Step{}
		if err := json.Unmarshal([]byte(r.Target.Data), &s); err == nil {
			steps[i] = s
		}
	}
	sort.Sort(models.OrderSteps{steps})
	return steps
}

func CreateDataFromStep(step models.Step, sec models.DataSecurity) (models.Data, error) {
	var sData models.Data
	sString, err := json.Marshal(step)
	if err != nil {
		return sData, err
	}
	sData = models.Data{
		ID:           step.ID,
		Data:         string(sString),
		DataSecurity: sec,
	}
	return sData, nil
}