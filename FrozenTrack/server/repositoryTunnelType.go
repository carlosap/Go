package main

import (
	"fmt"
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
)

// RepositoryGetAllTunnelTypes - returns all records from the table
func RepositoryGetAllTunnelTypes() (dbcontext.TunnelTypes, error) {
	t := dbcontext.TunnelType{}
	types, err := t.GetAll()
	if err != nil {
		return types, fmt.Errorf("error [RepositoryGetAllTunnelTypes]: %v", err)
	}

	return types, err
}

// RepositoryFindTunnelTypeById - returns a single records queried by id
func RepositoryFindTunnelTypeById(id string) (dbcontext.TunnelType, error) {

	t := dbcontext.TunnelType{}
	tunnelType, err := t.GetTunnelTypeByID(id)
	if err != nil {
		return t, fmt.Errorf("error [RepositoryFindTunnelTypeById]: %v", err)
	}

	return tunnelType, nil
}

// RepositoryDestroyTunnelTypeById - deletes a single records by id
func RepositoryDestroyTunnelTypeById(id string) (dbcontext.TunnelTypes, error) {

	tunneltype, err := RepositoryFindTunnelTypeById(id)
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyTunnelTypeById %s ", id)
	}

	err = tunneltype.Delete()
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyTunnelStateById -> failed to delete node %v", err)
	}

	return RepositoryGetAllTunnelTypes()
}

// RepositoryCreateTunnelType - creates a new record via payload (json)
func RepositoryCreateTunnelType(tunnelType *dbcontext.TunnelType) (dbcontext.TunnelTypes, error) {

	err := tunnelType.Insert()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryCreateTunnelType]: %v", err)
	}

	return RepositoryGetAllTunnelTypes()
}

// RepositoryUpdateTunnelType - updates an existing record via payload (json)
func RepositoryUpdateTunnelType(tunnelType *dbcontext.TunnelType) (dbcontext.TunnelTypes, error) {

	err := tunnelType.Update()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryUpdateTunnelType]: %v", err)
	}

	return RepositoryGetAllTunnelTypes()
}
