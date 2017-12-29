package models

import (
	"net/http"
)
//JEncoder is the type that handles the response encoding
type JEncoder struct {
	Writer http.ResponseWriter
}