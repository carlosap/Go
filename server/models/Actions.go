package models

import (
	"time"
)

// Action type
type Action struct {
	ID     string         `json:"id"`
	Action string         `json:"action"` // a description of what happened
	Time   time.Time      `json:"time"`   // the time at which it occurred
}

//Actions Multiple actions
type Actions []Action

func (a Actions) Len() int {
	return len(a)
}

func (a Actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Actions) Less(i, j int) bool {
	return a[i].Time.Before(a[j].Time)
}
