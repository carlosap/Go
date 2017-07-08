package httputil

import (
	"encoding/json"
	"net/http"
)

//JsonEncoder is an interface for encoding json http responses
type JsonEncoder interface {
	Encode(v interface{}) ([]byte, error)
	EncodeResponse(status int, v interface{}) error
}

//JEncoder is the type that handles the response encoding
type JEncoder struct {
	Writer http.ResponseWriter
}

//Encode encodes the given type as json and sets the appropriate http response headers
//for a json response
func (j JEncoder) Encode(v interface{}) ([]byte, error) {
	var data []byte
	var err error
	if v != nil {
		j.Writer.Header().Set("Content-Type", "application/json")
		data, err = json.MarshalIndent(v, "", " ")
	}
	return data, err
}

//EncodeResponse writes the response to the network in addition to Encode above
func (j JEncoder) EncodeResponse(status int, v interface{}) error {
	data, err := j.Encode(v)
	j.Writer.WriteHeader(status)
	j.Writer.Write(data)

	return err
}
