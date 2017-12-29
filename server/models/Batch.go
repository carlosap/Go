package models

import (
	"time"
)
//Batch entity
type Batch struct {
	ID             string    `json:"id"`
	AnalystID      string    `json:"analystId"`
	CaptureTime    time.Time `json:"captureTime"`
	Comments       string    `json:"comments"`
	Objective      string    `json:"objective"`
	Priority       []string  `json:"priority"`
	MGRS           string    `json:"mgrs"`
	Selectors      []string  `json:"selectors,omitempty"`
	FilesRemaining int       `json:"filesRemaining"`
	SourceDB       string    `json:"sourceDb"`
	Languages      []string  `json:"languages"`
	Step           string    `json:"step"`
	TaskForce      string    `json:"taskForce"`
	DataSecurity
}

//Batches Array Collection
type Batches []Batch

func (b Batches) Len() int {
	return len(b)
}

func (b Batches) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Batches) Less(i, j int) bool {
	return b[i].ID < b[j].ID
}