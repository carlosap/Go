package models

// Marking is mostly used to hold a marking values
type Marking struct {
	Display string   `yaml:"display" json:"display,omitempty"`
	Value   string   `yaml:"value" json:"value" cql:"value"`
	Portion string   `yaml:"portion" json:"portion,omitempty"`
	Groups  []string `yaml:"groups" json:"-"`
}

// Markings is an array of Marking
type Markings []Marking