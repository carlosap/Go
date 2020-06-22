package main

import (
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
	"github.com/gorilla/mux"
	"net/http"
)

func tunnelRequestRoutes() Routes {
	return Routes{
		Route{"CreateTunnelRequest", "POST", "/tunnel/requests", CreateTunnelRequest, true},
		Route{"UpdateTunnelRequest", "PUT", "/tunnel/requests", UpdateTunnelRequest, true},
		Route{"FetchTunnelRequests", "GET", "/tunnel/requests", FetchTunnelRequests, true},
		Route{"FetchTunnelRequest", "GET", "/tunnel/requests/{id}", FetchTunnelRequest, true},
		Route{"DeleteTunnelRequest", "DELETE", "/tunnel/requests/{id}", DestroyTunnelRequest, true},
	}
}

// FetchTunnelRequests - returns the latest collection
func FetchTunnelRequests(w http.ResponseWriter, r *http.Request) {

	requests, _ := RepositoryGetAllTunnelRequests()

	if len(requests) > 0 {
		Json(w, http.StatusOK, requests)
		return
	}

	JsonError(w, http.StatusInternalServerError, "failed to fetch tunnel requests")
}

// FetchTunnelRequest - returns a single node by id
func FetchTunnelRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	request, err := RepositoryFindTunnelRequestById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to fetch node with  request id "+id)
		return
	}

	Json(w, http.StatusOK, request)
}

// DestroyTunnelRequest - deletes a node by ID and returns the latest collection
func DestroyTunnelRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	requests, err := RepositoryDestroyTunnelRequestById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to delete tunnel with  id "+id)
		return
	}

	Json(w, http.StatusOK, requests)
}

// CreateTunnelRequest - creates a node and returns the latest collection
func CreateTunnelRequest(w http.ResponseWriter, r *http.Request) {
	request := &dbcontext.TunnelRequest{}

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [CreateNode]")
		return
	}

	userid := payload["UserID"].(string)
	if len(userid) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse UserID")
		return
	}
	request.UserID = userid

	hops := int(payload["Hops"].(float64))
	request.Hops = &hops

	poe := payload["Poe"].(string)
	if len(poe) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Poe")
		return
	}
	request.Poe = &poe

	pop := payload["Pop"].(string)
	if len(pop) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Pop")
		return
	}
	request.Pop = &pop

	active := payload["Actived"].(bool)
	request.Actived = &active

	requests, err := RepositoryCreateTunnelRequest(request)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new request")
		return
	}

	Json(w, http.StatusOK, requests)
}

// UpdateTunnelRequest - updates a node and returns the latest collection
func UpdateTunnelRequest(w http.ResponseWriter, r *http.Request) {

	request := &dbcontext.TunnelRequest{}
	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body")
		return
	}

	id := payload["TunnelRequestID"].(string)
	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelRequestID")
		return
	}
	request.TunnelRequestID = id

	userid := payload["UserID"].(string)
	if len(userid) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse UserID")
		return
	}
	request.UserID = userid

	hops := int(payload["Hops"].(float64))
	request.Hops = &hops

	poe := payload["Poe"].(string)
	if len(poe) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Poe")
		return
	}
	request.Poe = &poe

	pop := payload["Pop"].(string)
	if len(pop) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Pop")
		return
	}
	request.Pop = &pop

	requests, err := RepositoryUpdateTunnelRequest(request)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new tunnel request")
		return
	}

	Json(w, http.StatusOK, requests)
}
