package main

import (
	"fmt"
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
)

// RepositoryGetAllTunnelActions - returns all records from the table
func RepositoryGetAllTunnelActions() (dbcontext.TunnelActions, error) {
	t := dbcontext.TunnelAction{}
	actions, err := t.GetAll()
	if err != nil {
		return actions, fmt.Errorf("error [RepositoryGetAllTunnelActions]: %v", err)
	}

	return actions, err
}

// RepositoryFindTunnelActionById - returns a single records queried by id
func RepositoryFindTunnelActionById(id string) (dbcontext.TunnelAction, error) {

	t := dbcontext.TunnelAction{}
	action, err := t.GetTunnelActionByID(id)
	if err != nil {
		return t, fmt.Errorf("error [RepositoryFindTunnelActionById]: %v", err)
	}

	return action, nil
}

// RepositoryDestroyTunnelActionById - deletes a single records by id
func RepositoryDestroyTunnelActionById(id string) (dbcontext.TunnelActions, error) {

	action, err := RepositoryFindTunnelActionById(id)
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyTunnelActionById %s ", id)
	}

	err = action.Delete()
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyTunnelActionById -> failed to delete node %v", err)
	}

	return RepositoryGetAllTunnelActions()
}

// RepositoryCreateTunnelAction - creates a new record via payload (json)
func RepositoryCreateTunnelAction(action *dbcontext.TunnelAction) (dbcontext.TunnelActions, error) {

	err := action.Insert()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryCreateTunnelAction]: %v", err)
	}

	return RepositoryGetAllTunnelActions()
}

// RepositoryUpdateTunnelAction - updates an existing record via payload (json)
func RepositoryUpdateTunnelAction(action *dbcontext.TunnelAction) (dbcontext.TunnelActions, error) {

	err := action.Update()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryUpdateTunnelAction]: %v", err)
	}

	return RepositoryGetAllTunnelActions()
}
