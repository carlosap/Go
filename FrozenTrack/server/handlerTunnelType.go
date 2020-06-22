package main

import (
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
	"github.com/gorilla/mux"
	"net/http"
)

func tunnelTypeRoutes() Routes {
	return Routes{
		Route{"CreateTunnelType", "POST", "/tunnel/types", CreateTunnelType, true},
		Route{"UpdateTunnelType", "PUT", "/tunnel/types", UpdateTunnelType, true},
		Route{"FetchTunnelTypes", "GET", "/tunnel/types", FetchTunnelTypes, true},
		Route{"FetchTunnelType", "GET", "/tunnel/types/{id}", FetchTunnelType, true},
		Route{"DeleteTunnelType", "DELETE", "/tunnel/types/{id}", DestroyTunnelType, true},
	}
}

// FetchTunnelTypes - returns the latest collection
func FetchTunnelTypes(w http.ResponseWriter, r *http.Request) {

	states, _ := RepositoryGetAllTunnelTypes()

	if len(states) > 0 {
		Json(w, http.StatusOK, states)
		return
	}

	JsonError(w, http.StatusInternalServerError, "failed to fetch tunnel types")
}

// FetchTunnelType - returns a single node by id
func FetchTunnelType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	state, err := RepositoryFindTunnelTypeById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to fetch tunnel type with  id "+id)
		return
	}

	Json(w, http.StatusOK, state)
}

// DestroyTunnelType - deletes a node by ID and returns the latest collection
func DestroyTunnelType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	tunneltype, err := RepositoryDestroyTunnelTypeById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to delete tunnel type with  id "+id)
		return
	}

	Json(w, http.StatusOK, tunneltype)
}

// CreateTunnelType - creates a node and returns the latest collection
func CreateTunnelType(w http.ResponseWriter, r *http.Request) {
	tunnelType := &dbcontext.TunnelType{}

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [CreateTunnelType]")
		return
	}

	name := payload["Name"].(string)
	if len(name) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Name")
		return
	}
	tunnelType.Name = &name

	tunnelTypes, err := RepositoryCreateTunnelType(tunnelType)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new tunnel type")
		return
	}

	Json(w, http.StatusOK, tunnelTypes)
}

// UpdateTunnelType - updates a node and returns the latest collection
func UpdateTunnelType(w http.ResponseWriter, r *http.Request) {

	tunnelType := &dbcontext.TunnelType{}
	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body")
		return
	}

	id := payload["TunnelTypeID"].(string)
	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelTypeID")
		return
	}
	tunnelType.TunnelTypeID = id

	name := payload["Name"].(string)
	if len(name) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Name")
		return
	}
	tunnelType.Name = &name

	tunnelTypes, err := RepositoryUpdateTunnelType(tunnelType)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to update tunnel types")
		return
	}

	Json(w, http.StatusOK, tunnelTypes)
}
