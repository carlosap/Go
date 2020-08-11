package dbcontext

import (
	"fmt"
)

//============================================Functional Requirements ===================================

//GetAll returns all the alias info table rows
func (a Storageaccount) GetAll() ([]Storageaccount, error) {
	db, err := NewDbContext()
	if err != nil {
		return nil, fmt.Errorf("error: failed to connect to db %v", err)
	}

	var aliases []Storageaccount
	err = db.Pgdb.Model(&aliases).Select()
	if err != nil {
		return nil, fmt.Errorf("error: failed to fetch all storageaccount %v", err)
	}

	_ = db.Close()
	return aliases, err
}

//Insert
func (a *Storageaccount) Insert() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Insert(a)

	if err != nil {
		return fmt.Errorf("error: failed insert storageaccount %v", err)
	}

	_ = db.Close()
	return err
}

//Delete
func (a *Storageaccount) Delete() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Delete(a)

	if err != nil {
		return fmt.Errorf("error: failed to delete storageaccount %v", err)
	}

	_ = db.Close()
	return err
}

//Update
func (a *Storageaccount) Update() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Update(a)

	if err != nil {
		return fmt.Errorf("error: failed to update storageaccount %v", err)
	}

	_ = db.Close()
	return err
}

//============================================Common Driver Requirements===================================
// storageaccount
type Storageaccount struct {
	tableName struct{} `pg:"azmonitor.storageaccount,alias:t"`

	Resourceid        *string `pg:"resourceid"`
	Resourcegroup     *string `pg:"resourcegroup"`
	Servicename       *string `pg:"servicename"`
	Cost              *string `pg:"cost"`
	Resourcetype      *string `pg:"resourcetype"`
	Resourcelocation  *string `pg:"resourcelocation"`
	Consumptiontype   *string `pg:"consumptiontype"`
	Meter             *string `pg:"meter"`
	Availability      *string `pg:"availability"`
	Totaltransactions *string `pg:"totaltransactions"`
	E2elatency        *string `pg:"e2elatency"`
	Serverlantency    *string `pg:"serverlantency"`
	Failures          *string `pg:"failures"`
	Capacity          *string `pg:"capacity"`
}

//AliasInfos
type Storageaccounts []Storageaccount
