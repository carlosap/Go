package main

import (
	"fmt"
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
)

// RepositoryGetAllTunnelStates - returns all records from the table
func RepositoryGetAllTunnelStates() (dbcontext.TunnelStates, error) {
	t := dbcontext.TunnelState{}
	states, err := t.GetAll()
	if err != nil {
		return states, fmt.Errorf("error [RepositoryGetAllTunnelStates]: %v", err)
	}

	return states, err
}

// RepositoryFindTunnelStateById - returns a single records queried by id
func RepositoryFindTunnelStateById(id string) (dbcontext.TunnelState, error) {

	t := dbcontext.TunnelState{}
	state, err := t.GetTunnelStateByID(id)
	if err != nil {
		return t, fmt.Errorf("error [RepositoryFindTunnelStateById]: %v", err)
	}

	return state, nil
}

// RepositoryDestroyTunnelStateById - deletes a single records by id
func RepositoryDestroyTunnelStateById(id string) (dbcontext.TunnelStates, error) {

	state, err := RepositoryFindTunnelStateById(id)
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyTunnelStateById %s ", id)
	}

	err = state.Delete()
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyTunnelStateById -> failed to delete node %v", err)
	}

	return RepositoryGetAllTunnelStates()
}

// RepositoryCreateTunnelState - creates a new record via payload (json)
func RepositoryCreateTunnelState(state *dbcontext.TunnelState) (dbcontext.TunnelStates, error) {

	err := state.Insert()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryCreateTunnelState]: %v", err)
	}

	return RepositoryGetAllTunnelStates()
}

// RepositoryUpdateTunnelState - updates an existing record via payload (json)
func RepositoryUpdateTunnelState(state *dbcontext.TunnelState) (dbcontext.TunnelStates, error) {

	err := state.Update()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryUpdateTunnelState]: %v", err)
	}

	return RepositoryGetAllTunnelStates()
}
