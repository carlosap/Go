package models

import (
	"time"
)

//File Type
type File struct {
	ID          string    `json:"id"`
	BatchID     string    `json:"batchId"`
	AnalystID   string    `json:"analystId"`
	Priority    []string  `json:"priority"`
	Objective   string    `json:"objective"`
	CaptureTime time.Time `json:"captureTime"`
	Step        string    `json:"step"`
	SourceDB    string    `json:"sourceDb"`
	Languages   []string  `json:"languages"`
	Type        string    `json:"type"`
	MGRS        string    `json:"mgrs"`
	Unit        string    `json:"unit"`
	Comments    []string   `json:"comments"`
	Hash        string    `json:"md5hash"`
	MediumUUID  string    `json:"mediumUuid"`
	DeviceType  string    `json:"deviceType"`
	DeviceName  string    `json:"deviceName"`
	Name        string    `json:"name"`
	Created     time.Time `json:"dateCreated"`
	Modified    time.Time `json:"lastModified"`
	User        string    `json:"user"`
	Dismissed   bool      `json:"dismissed"`
	Description string    `json:"description"`
	Size        string    `json:"size"`
	EXIF        string    `json:"exif"`
	Selectors   []string  `json:"selectors"`
	DataSecurity
}

//Files Type
type Files []File

