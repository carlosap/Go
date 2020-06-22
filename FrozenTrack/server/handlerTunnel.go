package main

import (
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
	"github.com/gorilla/mux"
	"net/http"
)

func tunnelRoutes() Routes {
	return Routes{
		Route{"CreateTunnel", "POST", "/tunnels", CreateTunnel, true},
		Route{"UpdateTunnel", "PUT", "/tunnels", UpdateTunnel, true},
		Route{"FetchTunnel", "GET", "/tunnels", FetchTunnels, true},
		Route{"FetchTunnel", "GET", "/tunnels/{id}", FetchTunnel, true},
		Route{"DeleteTunnel", "DELETE", "/tunnels/{id}", DestroyTunnel, true},
	}
}

// FetchTunnels - returns the latest collection
func FetchTunnels(w http.ResponseWriter, r *http.Request) {

	nodes, _ := RepositoryGetAllTunnels()

	if len(nodes) > 0 {
		Json(w, http.StatusOK, nodes)
		return
	}

	JsonError(w, http.StatusInternalServerError, "failed to fetch tunnels")
}

// FetchTunnel - returns a single node by id
func FetchTunnel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	tunnel, err := RepositoryFindTunnelById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to fetch node with  id "+id)
		return
	}

	Json(w, http.StatusOK, tunnel)
}

// DestroyTunnel - deletes a node by ID and returns the latest collection
func DestroyTunnel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	tunnels, err := RepositoryDestroyTunnelById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to delete tunnel with  id "+id)
		return
	}

	Json(w, http.StatusOK, tunnels)
}

// CreateTunnel - creates a node and returns the latest collection
func CreateTunnel(w http.ResponseWriter, r *http.Request) {
	tunnel := &dbcontext.Tunnel{}

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [CreateNode]")
		return
	}

	stateid := payload["TunnelStateID"].(string)
	if len(stateid) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelStateID")
		return
	}
	tunnel.TunnelState = stateid

	actionId := payload["TunnelActionID"].(string)
	if len(actionId) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelActionID")
		return
	}
	tunnel.TunnelActionID = actionId

	typeID := payload["TunnelTypeID"].(string)
	if len(typeID) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelTypeID")
		return
	}
	tunnel.TunnelTypeID = typeID

	hops := int(payload["Hops"].(float64))
	tunnel.Hops = &hops

	poe := payload["Poe"].(string)
	if len(poe) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Poe")
		return
	}

	tunnel.Poe = &poe

	pop := payload["Pop"].(string)
	if len(pop) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Pop")
		return
	}

	tunnel.Pop = &pop

	ip := payload["IpAddress"].(string)
	if len(ip) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse IpAddress")
		return
	}

	tunnel.IpAddress = &ip

	tunnels, err := RepositoryCreateTunnel(tunnel)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new node")
		return
	}

	Json(w, http.StatusOK, tunnels)
}

// UpdateTunnel - updates a node and returns the latest collection
func UpdateTunnel(w http.ResponseWriter, r *http.Request) {

	tunnel := &dbcontext.Tunnel{}
	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body")
		return
	}

	id := payload["TunnelID"].(string)
	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelID")
		return
	}
	tunnel.TunnelID = id

	stateid := payload["TunnelStateID"].(string)
	if len(stateid) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelStateID")
		return
	}
	tunnel.TunnelState = stateid

	actionId := payload["TunnelActionID"].(string)
	if len(actionId) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelActionID")
		return
	}

	tunnel.TunnelActionID = actionId

	typeID := payload["TunnelTypeID"].(string)
	if len(typeID) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelTypeID")
		return
	}

	tunnel.TunnelTypeID = typeID

	hops := int(payload["Hops"].(float64))
	tunnel.Hops = &hops

	poe := payload["Poe"].(string)
	if len(poe) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Poe")
		return
	}

	tunnel.Poe = &poe

	pop := payload["Pop"].(string)
	if len(pop) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Pop")
		return
	}

	tunnel.Pop = &pop

	ip := payload["IpAddress"].(string)
	if len(ip) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse IpAddress")
		return
	}

	tunnel.IpAddress = &ip

	tunnels, err := RepositoryUpdateTunnel(tunnel)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new node")
		return
	}

	Json(w, http.StatusOK, tunnels)
}
