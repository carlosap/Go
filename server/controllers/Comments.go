package controllers
import (

	"github.com/Go/server/util/logging"
	"encoding/json"
	"net/http"	

	"github.com/Go/server/dbContext"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
	"github.com/Go/server/models"

) 

func RegisterCommentsEndpoints(h http.Handler) {
	m := h.(*martini.ClassicMartini)
	m.Group("/comments", func(r martini.Router) {
		m.Get("/:hashValue/files", GetCommentFilesByHashHandler)
		m.Get("/:hashValue/selectors", GetCommentSelectorByHashHandler)
	})
}

//Example: http://localhost:5001/comments/7c630248a4cba073e33f49f9795453e6/files
func GetCommentFilesByHashHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var files = make([]models.File, 0)
	var records = make([]models.Record, 0)
	var err error
	hashValue, ok := params["hashValue"]
	if !ok || len(hashValue) == 0 {
		return http.StatusBadRequest, encoder.Must(enc.Encode(files))
	}
	records, err = dbContext.GetRecordsByHash(hashValue)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(records))
	}
	if len(records) > 0 {
		for _, r := range records {
			if r.Target.Data != "" {
				if r.LinkType == "files" {
					var file models.File
					err := json.Unmarshal([]byte(r.Target.Data), &file)
					if err != nil {
						logging.Errorf("Error %+v Unmarshal data from %s", err, hashValue)
					}	
					files = append(files, file)
				}		
			}
		}
	}
	return http.StatusOK, encoder.Must(enc.Encode(files))
}

func GetCommentSelectorByHashHandler(r *http.Request, w http.ResponseWriter, params martini.Params, enc encoder.Encoder) (int, []byte) {
	var selectors = make([]models.Selector, 0)
	var records = make([]models.Record, 0)
	var err error
	hashValue, ok := params["hashValue"]
	if !ok || len(hashValue) == 0 {
		return http.StatusBadRequest, encoder.Must(enc.Encode(selectors))
	}
	records, err = dbContext.GetRecordsByHash(hashValue)
	if err != nil {
		return http.StatusBadRequest, encoder.Must(enc.Encode(records))
	}
	if len(records) > 0 {
		for _, r := range records {
			if r.Target.Data != "" {
				if r.LinkType == "selectors" {
					var selector models.Selector
					err := json.Unmarshal([]byte(r.Target.Data), &selector)
					if err != nil {
						logging.Errorf("Error %+v Unmarshal data from %s", err, hashValue)
					}	
					selectors = append(selectors, selector)
				}		
			}
		}
	}
	return http.StatusOK, encoder.Must(enc.Encode(selectors))
}
