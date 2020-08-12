package dbcontext

import (
	"fmt"
	"time"
)

//============================================Functional Requirements ===================================

//GetAll returns all the alias info table rows
func (a Virtualmachine) GetAll() ([]Virtualmachine, error) {
	db, err := NewDbContext()
	if err != nil {
		return nil, fmt.Errorf("error: failed to connect to db %v", err)
	}

	var aliases []Virtualmachine
	err = db.Pgdb.Model(&aliases).Select()
	if err != nil {
		return nil, fmt.Errorf("error: failed to fetch all Virtualmachine %v", err)
	}

	_ = db.Close()
	return aliases, err
}

//Insert
func (a *Virtualmachine) Insert() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Insert(a)

	if err != nil {
		return fmt.Errorf("error: failed insert Virtualmachine %v", err)
	}

	_ = db.Close()
	return err
}

//Delete
func (a *Virtualmachine) Delete() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Delete(a)

	if err != nil {
		return fmt.Errorf("error: failed to delete Virtualmachine %v", err)
	}

	_ = db.Close()
	return err
}

//Update
func (a *Virtualmachine) Update() error {
	db, err := NewDbContext()

	if err != nil {
		return fmt.Errorf("error: failed to connect to db %v", err)
	}

	err = db.Pgdb.Update(a)

	if err != nil {
		return fmt.Errorf("error: failed to update Virtualmachine %v", err)
	}

	_ = db.Close()
	return err
}

//============================================Common Driver Requirements===================================
// Virtualmachine
type Virtualmachine struct {
	tableName struct{} `pg:"azmonitor.virtualmachine,alias:t"`

	ID                  string                 `pg:"id,pk,type:uuid"`
	Resourceid          *string                `pg:"resourceid"`
	Resourcegroup       *string                `pg:"resourcegroup"`
	Servicename         *string                `pg:"servicename"`
	Cost                *string                `pg:"cost"`
	Resourcetype        *string                `pg:"resourcetype"`
	Resourcelocation    *string                `pg:"resourcelocation"`
	Consumptiontype     *string                `pg:"consumptiontype"`
	Meter               *string                `pg:"meter"`
	Cpuutilization      *string                `pg:"cpuutilization"`
	Availablememory     *string                `pg:"availablememory"`
	Disklatency         *string                `pg:"disklatency"`
	Diskiops            *string                `pg:"diskiops"`
	Diskbytespersec     *string                `pg:"diskbytespersec"`
	Networksentrate     *string                `pg:"networksentrate"`
	Networkreceivedrate *string                `pg:"networkreceivedrate"`
	Datecreated         *time.Time             `pg:"datecreated"`
	Lastupdated         *time.Time             `pg:"lastupdated"`
	Reportstartdate     *string                `pg:"reportstartdate"`
	Reportenddate       *string                `pg:"reportenddate"`
	Data                map[string]interface{} `pg:"data"`
}

//AliasInfos
type Virtualmachines []Virtualmachine
