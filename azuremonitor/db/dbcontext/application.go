package dbcontext

import (
	"fmt"
	"time"
)

//============================================Functional Requirements ===================================

//GetAll returns all the alias info table rows
func (a Application) GetAll() ([]Application, error) {
	db, err := NewDbContext()
	if err != nil {
		return nil, fmt.Errorf("error: failed to connect to db %v", err)
	}

	var aliases []Application
	err = db.Pgdb.Model(&aliases).Select()
	if err != nil {
		return nil, fmt.Errorf("error: failed to fetch all Application %v", err)
	}

	_ = db.Close()
	return aliases, err
}

//Insert
func (a *Application) Insert() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Insert(a)

	if err != nil {
		return fmt.Errorf("error: failed insert Application %v", err)
	}

	_ = db.Close()
	return err
}

//Delete
func (a *Application) Delete() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Delete(a)

	if err != nil {
		return fmt.Errorf("error: failed to delete Application %v", err)
	}

	_ = db.Close()
	return err
}

//Update
func (a *Application) Update() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Update(a)

	if err != nil {
		return fmt.Errorf("error: failed to update Application %v", err)
	}

	_ = db.Close()
	return err
}

//============================================Common Driver Requirements===================================
// application
type Application struct {
	tableName struct{} `pg:"azmonitor.application,alias:t"`

	Applicationid  int        `pg:"applicationid,pk"`
	SubscriptionID *string    `pg:"subscription_id"`
	Name           *string    `pg:"name"`
	TenantID       *string    `pg:"tenant_id"`
	GrantType      *string    `pg:"grant_type"`
	ClientID       *string    `pg:"client_id"`
	ClientSecret   *string    `pg:"client_secret"`
	Lastmodified   *time.Time `pg:"lastmodified"`
}

//AliasInfos
type Applications []Application
