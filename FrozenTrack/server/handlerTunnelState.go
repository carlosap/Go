package main

import (
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
	"github.com/gorilla/mux"
	"net/http"
)

func tunnelStateRoutes() Routes {
	return Routes{
		Route{"CreateTunnelState", "POST", "/tunnel/states", CreateTunnelState, true},
		Route{"UpdateTunnelState", "PUT", "/tunnel/states", UpdateTunnelState, true},
		Route{"FetchTunnelStates", "GET", "/tunnel/states", FetchTunnelStates, true},
		Route{"FetchTunnelState", "GET", "/tunnel/states/{id}", FetchTunnelState, true},
		Route{"DeleteTunnelState", "DELETE", "/tunnel/states/{id}", DestroyTunnelState, true},
	}
}

// FetchTunnelStates - returns the latest collection
func FetchTunnelStates(w http.ResponseWriter, r *http.Request) {

	states, _ := RepositoryGetAllTunnelStates()

	if len(states) > 0 {
		Json(w, http.StatusOK, states)
		return
	}

	JsonError(w, http.StatusInternalServerError, "failed to fetch tunnel states")
}

// FetchTunnelState - returns a single node by id
func FetchTunnelState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	state, err := RepositoryFindTunnelStateById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to fetch node with  id "+id)
		return
	}

	Json(w, http.StatusOK, state)
}

// DestroyTunnelState - deletes a node by ID and returns the latest collection
func DestroyTunnelState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	states, err := RepositoryDestroyTunnelStateById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to delete tunnel with  id "+id)
		return
	}

	Json(w, http.StatusOK, states)
}

// CreateTunnelState - creates a node and returns the latest collection
func CreateTunnelState(w http.ResponseWriter, r *http.Request) {
	state := &dbcontext.TunnelState{}

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [CreateTunnelState]")
		return
	}

	name := payload["Name"].(string)
	if len(name) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Name")
		return
	}
	state.Name = &name

	tunnelStates, err := RepositoryCreateTunnelState(state)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new tunnel state")
		return
	}

	Json(w, http.StatusOK, tunnelStates)
}

// UpdateTunnelState - updates a node and returns the latest collection
func UpdateTunnelState(w http.ResponseWriter, r *http.Request) {

	state := &dbcontext.TunnelState{}
	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body")
		return
	}

	id := payload["TunnelStateID"].(string)
	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelStateID")
		return
	}
	state.TunnelStateID = id

	name := payload["Name"].(string)
	if len(name) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Name")
		return
	}
	state.Name = &name

	tunnelStates, err := RepositoryUpdateTunnelState(state)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to update tunnel state")
		return
	}

	Json(w, http.StatusOK, tunnelStates)
}
