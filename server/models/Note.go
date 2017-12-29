package models

import (
	"time"
	"github.com/gocql/gocql"
)

//Note structure
type Note struct {
	NoteID     gocql.UUID `json:"noteID" cql:"noteid"`
	Name       string     `json:"name" cql:"name"`
	Org        string     `json:"org" cql:"org"`
	Email      string     `json:"email" cql:"email"`
	Note       string     `json:"note" cql:"note"`
	DateOfNote time.Time  `json:"dateOfNote" cql:"dateofnote"`
}