package costmanagement


type BatchResponse struct {
	Responses []Responses `json:"responses"`
}
type Name struct {
	Value          string `json:"value"`
	LocalizedValue string `json:"localizedValue"`
}
type Data struct {
	TimeStamp string `json:"timeStamp"`
	Total float64 `json:"total"`
	Average   float64   `json:"average"`
}
type Timeseries struct {
	Metadatavalues []interface{} `json:"metadatavalues"`
	Data           []Data        `json:"data"`
}
type Value struct {
	ID                 string       `json:"id"`
	Type               string       `json:"type"`
	Name               Name         `json:"name"`
	DisplayDescription string       `json:"displayDescription"`
	Unit               string       `json:"unit"`
	Timeseries         []Timeseries `json:"timeseries"`
	ErrorCode          string       `json:"errorCode"`
}
type Content struct {
	Cost           int       `json:"cost"`
	Timespan       string `json:"timespan"`
	Interval       string    `json:"interval"`
	Value          []Value   `json:"value"`
	Namespace      string    `json:"namespace"`
	Resourceregion string    `json:"resourceregion"`
}
type Responses struct {
	HTTPStatusCode int     `json:"httpStatusCode"`
	Content        Content `json:"content"`
}

