package main

import (
	"fmt"
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
)

// RepositoryGetAllTunnelRequests - returns all records from the table
func RepositoryGetAllTunnelRequests() (dbcontext.TunnelRequests, error) {
	t := dbcontext.TunnelRequest{}
	requests, err := t.GetAll()
	if err != nil {
		return requests, fmt.Errorf("error [RepositoryGetAllTunnelRequests]: %v", err)
	}

	return requests, err
}

// RepositoryFindTunnelRequestById - returns a single records queried by id
func RepositoryFindTunnelRequestById(id string) (dbcontext.TunnelRequest, error) {

	r := dbcontext.TunnelRequest{}
	request, err := r.GetTunnelRequestByID(id)
	if err != nil {
		return r, fmt.Errorf("error [RepositoryFindTunnelRequestById]: %v", err)
	}

	return request, nil
}

// RepositoryDestroyTunnelRequestById - deletes a single records by id
func RepositoryDestroyTunnelRequestById(id string) (dbcontext.TunnelRequests, error) {

	r, err := RepositoryFindTunnelRequestById(id)
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryFindTunnelRequestById %s ", id)
	}

	err = r.Delete()
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryFindTunnelRequestById -> failed to delete request by id %v", err)
	}

	return RepositoryGetAllTunnelRequests()
}

// RepositoryCreateTunnelRequest - creates a new record via payload (json)
func RepositoryCreateTunnelRequest(tunnelRequest *dbcontext.TunnelRequest) (dbcontext.TunnelRequests, error) {

	err := tunnelRequest.Insert()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryCreateTunnelRequest]: %v", err)
	}

	return RepositoryGetAllTunnelRequests()
}

// RepositoryUpdateTunnelRequest - updates an existing record via payload (json)
func RepositoryUpdateTunnelRequest(tunnelRequest *dbcontext.TunnelRequest) (dbcontext.TunnelRequests, error) {

	err := tunnelRequest.Update()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryUpdateTunnelRequest]: %v", err)
	}

	return RepositoryGetAllTunnelRequests()
}
