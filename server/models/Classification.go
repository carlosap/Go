package models

// Classification holds a classification object.
type Classification struct {
	Level                string   `json:"level" cql:"level"`
	SCIs                 Markings `json:"scis" cql:"scis"`
	SAPs                 Markings `json:"saps" cql:"saps"`
	AEA                  Marking  `json:"aea" cql:"aea"`
	Owners               []string `json:"owners" cql:"owners"`
	Disseminations       Markings `json:"disseminations" cql:"disseminations"`
	OtherDisseminations  Markings `json:"otherDisseminations" cql:"otherdisseminations"`
	ClassificationString string   `json:"classificationString"`
	PortionString        string   `json:"portionString"`
}
