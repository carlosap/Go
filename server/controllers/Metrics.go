package controllers
import (
	"time"
	"math"
	"strings"
	"net/http"
	"github.com/Go/server/util/logging"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"github.com/Go/server/models"

) 

const (
	unitPrecision = 0.05
	decimalPlaces = 2
	refNumPivot				= 2
	triageStartTimeToken = "ingested from cypher"
	triageEndTimeToken = "moved from triage to translate"
	translateStartTimeToken = "moved from triage to translate"
	translateEndTimeToken   = "moved from translate to exploit"

	exploitStartTimeToken   = "moved from translate to exploit"
	exploitEndTimeToken   = "moved from exploit to report"

)

func RegisterMetricsEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/metrics", func(r martini.Router) {
		m.Get("/:fileIDParam", getTimeCompletionMetricsHandler)
	})
}

//TODO MAKE STEP ALSO DYMIC
func getTimeCompletionMetricsHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {

	fileIds, ok := params["fileIDParam"]
	if !ok || len(fileIds) == 0 {
		logging.Errorf("File Id cannot be empty!")
		return http.StatusBadRequest, encoder.Must(enc.Encode(fileIds))
	}

	vTriage := stepTimeCompletionAvg(fileIds, "triage")
	vTranslate := stepTimeCompletionAvg(fileIds, "translate")
	vExploit := stepTimeCompletionAvg(fileIds, "exploit")

	m := map[string]interface{}{
		"triage": map[string]interface{} {
			"notes": triageStartTimeToken + " - " + triageEndTimeToken,
			"hours": vTriage,
		},
		"translate": map[string]interface{} {
			"notes": translateStartTimeToken + " - " + translateEndTimeToken,
			"hours": vTranslate,
		},
		"exploit": map[string]interface{} {
			"notes": exploitStartTimeToken + " - " + exploitEndTimeToken,
			"hours": vExploit,
		},
	}

	return http.StatusOK, encoder.Must(enc.Encode(m))

}

func stepTimeCompletionAvg(strIDs string, step string) float64 {
	elapseList := stepTimeCompletionList(strIDs, step)
  avg := avgFloat64List(elapseList)
	return roundFloat64(avg, unitPrecision, decimalPlaces)
}

//getTriageTimeCompletionList - 
func stepTimeCompletionList(strIDs string, step string) []float64 {
	var actions = make(models.Actions, 0)
	var retVal = make([]float64, 0)
	var err error
	ids := strings.Split(strIDs, ",")
	for _, id := range ids {
		actions, err = getActionsByFileID(id)
		if err != nil {
			logging.Errorf("Unable to get Actions for file %s, Error: %+v.", id, err)
		}

		switch step {

		case "triage":
			filterTriage := filterActionTokens(actions, triageStartTimeToken, triageEndTimeToken)
			if len(filterTriage) == refNumPivot {
				retVal = append(retVal, getActionHoursByEventNames(filterTriage, triageStartTimeToken, triageEndTimeToken)...)
			}

		case "translate":
			filterTranslate := filterActionTokens(actions, translateStartTimeToken, translateEndTimeToken)
			if len(filterTranslate) == refNumPivot {
				retVal = append(retVal, getActionHoursByEventNames(filterTranslate, translateStartTimeToken, translateEndTimeToken)...)
			}

		case "exploit":
			filterExploit := filterActionTokens(actions, exploitStartTimeToken, exploitEndTimeToken)
			if len(filterExploit) == refNumPivot {
				retVal = append(retVal, getActionHoursByEventNames(filterExploit, exploitStartTimeToken, exploitEndTimeToken)...)
			}

		}
	}
	return retVal
}

func getActionHoursByEventNames(filterActions models.Actions, startStrToken string, endStrToken string) []float64 {
	var duration float64
	var hours = make([]float64, 0)
	if len(filterActions) > 0 && len(startStrToken) > 0 && len(endStrToken) > 0 {
		duration = elapsedLookUp(filterActions, startStrToken, endStrToken)
		if duration > 0 {
			hours = append(hours, duration)
		}
	}
	return hours
}

func filterActionTokens(actions models.Actions, startStrToken string, endStrToken string) models.Actions {
	var fAction = make(models.Actions, 0)
	for _, a := range actions {
		name := strings.ToLower(a.Action)
		if name == startStrToken || name == endStrToken {
			fAction = append(fAction, a)
		}
	}

	return fAction
}

func elapsedLookUp(actions models.Actions, startStrToken string, endStrToken string) float64 {
	t1 := time.Time{}
	t2 := time.Time{}
	var h float64

  for _,a := range actions {
		aName := strings.ToLower(a.Action)
		switch aName {

		case startStrToken:
			t1 = a.Time

		case endStrToken:
			t2 = a.Time

		}
	}

	if t2.Second() > 0 && t1.Second() > 0 {
		diff := t2.Sub(t1)
	  h = roundFloat64(diff.Hours(), unitPrecision, decimalPlaces)
	}

	return h
}

//Round adding precision to float64 and decimal places
func roundFloat64(x, unit float64, decimalPlaces int) float64 {
	var round float64
	pow := math.Pow(10, float64(decimalPlaces))
	digit := pow * x
	_, d := math.Modf(digit)
	if x > 0 {
		if d >= unit {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
	} else {
		if d >= unit {
			round = math.Floor(digit)
		} else {
			round = math.Ceil(digit)
		}
	}

	return round / pow
}

func avgFloat64List(x []float64) float64 {
	var t float64
	if len(x) > 0 {
		for _, v := range x {
			t += v
		}
	}

	if t == 0 {
		return 0.0
	}

	return t /float64(len(x))
}