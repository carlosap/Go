package dbcontext

import (
	"fmt"
	"time"
)

//============================================Functional Requirements ===================================

//GetAll returns all the alias info table rows
func (a AzureMonitor) GetAll() ([]AzureMonitor, error) {
	db, err := NewDbContext()
	if err != nil {
		return nil, fmt.Errorf("error: failed to connect to db %v", err)
	}

	var aliases []AzureMonitor
	err = db.Pgdb.Model(&aliases).Select()
	if err != nil {
		return nil, fmt.Errorf("error: failed to fetch all Authv2s %v", err)
	}

	_ = db.Close()
	return aliases, err
}

//Insert
func (a *AzureMonitor) Insert() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Insert(a)

	if err != nil {
		return fmt.Errorf("error: failed insert AzureMonitor %v", err)
	}

	_ = db.Close()
	return err
}

//Delete
func (a *AzureMonitor) Delete() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Delete(a)

	if err != nil {
		return fmt.Errorf("error: failed to delete AzureMonitor %v", err)
	}

	_ = db.Close()
	return err
}

//Update
func (a *AzureMonitor) Update() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Update(a)

	if err != nil {
		return fmt.Errorf("error: failed to update AzureMonitor %v", err)
	}

	_ = db.Close()
	return err
}

//============================================Common Driver Requirements===================================
// AliasInfo
type AzureMonitor struct {
	tableName struct{} `pg:"azuremonitor.azuremonitor,alias:t"`

	AzuremonitorID string     `pg:"azuremonitor_id,pk,type:uuid"`
	Name           *string    `pg:"name"`
	Hostname       *int       `pg:"hostname"`
	Lastmodified   *time.Time `pg:"lastmodified"`
}

//AliasInfos
type AzureMonitors []AzureMonitor
