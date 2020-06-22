package main

import (
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
	"github.com/gorilla/mux"
	"net/http"
)

func tunnelActionRoutes() Routes {
	return Routes{
		Route{"CreateTunnelAction", "POST", "/tunnel/actions", CreateTunnelAction, true},
		Route{"UpdateTunnelAction", "PUT", "/tunnel/actions", UpdateTunnelAction, true},
		Route{"FetchTunnelActions", "GET", "/tunnel/actions", FetchTunnelActions, true},
		Route{"FetchTunnelAction", "GET", "/tunnel/actions/{id}", FetchTunnelAction, true},
		Route{"DeleteTunnelAction", "DELETE", "/tunnel/actions/{id}", DestroyTunnelAction, true},
	}
}

// FetchTunnelsAction - returns the latest collection
func FetchTunnelActions(w http.ResponseWriter, r *http.Request) {

	actions, _ := RepositoryGetAllTunnelActions()

	if len(actions) > 0 {
		Json(w, http.StatusOK, actions)
		return
	}

	JsonError(w, http.StatusInternalServerError, "failed to fetch tunnel actions")
}

// FetchTunnelAction - returns a single node by id
func FetchTunnelAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	action, err := RepositoryFindTunnelActionById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to fetch node with  id "+id)
		return
	}

	Json(w, http.StatusOK, action)
}

// DestroyTunnelAction - deletes a node by ID and returns the latest collection
func DestroyTunnelAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	actions, err := RepositoryDestroyTunnelActionById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to delete tunnel with  id "+id)
		return
	}

	Json(w, http.StatusOK, actions)
}

// CreateTunnelAction - creates a node and returns the latest collection
func CreateTunnelAction(w http.ResponseWriter, r *http.Request) {
	action := &dbcontext.TunnelAction{}

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [CreateNode]")
		return
	}

	name := payload["Name"].(string)
	if len(name) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Name")
		return
	}
	action.Name = &name

	tunnelActions, err := RepositoryCreateTunnelAction(action)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new tunnel action")
		return
	}

	Json(w, http.StatusOK, tunnelActions)
}

// UpdateTunnelAction - updates a node and returns the latest collection
func UpdateTunnelAction(w http.ResponseWriter, r *http.Request) {

	action := &dbcontext.TunnelAction{}
	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body")
		return
	}

	id := payload["TunnelActionID"].(string)
	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelActionID")
		return
	}
	action.TunnelActionID = id

	name := payload["Name"].(string)
	if len(name) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Name")
		return
	}
	action.Name = &name

	tunnelActions, err := RepositoryUpdateTunnelAction(action)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to update tunnel action")
		return
	}

	Json(w, http.StatusOK, tunnelActions)
}
