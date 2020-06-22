package main

import (
	"fmt"
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
)

// RepositoryGetAllTunnels - returns all records from the table
func RepositoryGetAllTunnels() (dbcontext.Tunnels, error) {
	t := dbcontext.Tunnel{}
	tunnels, err := t.GetAll()
	if err != nil {
		return tunnels, fmt.Errorf("error [RepositoryGetAllTunnels]: %v", err)
	}

	return tunnels, err
}

// RepositoryFindTunnelById - returns a single records queried by id
func RepositoryFindTunnelById(id string) (dbcontext.Tunnel, error) {

	t := dbcontext.Tunnel{}
	tunnel, err := t.GetTunnelByID(id)
	if err != nil {
		return t, fmt.Errorf("error [RepositoryFindTunnelById]: %v", err)
	}

	return tunnel, nil
}

// RepositoryDestroyTunnelById - deletes a single records by id
func RepositoryDestroyTunnelById(id string) (dbcontext.Tunnels, error) {

	tunnel, err := RepositoryFindTunnelById(id)
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyNodeById %s ", id)
	}

	err = tunnel.Delete()
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyNodeById -> failed to delete node %v", err)
	}

	return RepositoryGetAllTunnels()
}

// RepositoryCreateTunnel - creates a new record via payload (json)
func RepositoryCreateTunnel(tunnel *dbcontext.Tunnel) (dbcontext.Tunnels, error) {

	err := tunnel.Insert()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryCreateTunnel]: %v", err)
	}

	return RepositoryGetAllTunnels()
}

// RepositoryUpdateTunnel - updates an existing record via payload (json)
func RepositoryUpdateTunnel(tunnel *dbcontext.Tunnel) (dbcontext.Tunnels, error) {

	err := tunnel.Update()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryUpdateTunnel]: %v", err)
	}

	return RepositoryGetAllTunnels()
}
