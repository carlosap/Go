package main

import (
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
	"github.com/gorilla/mux"
	"net/http"
)

func nodeRoutes() Routes {
	return Routes{
		Route{"Index", "GET", "/", Index, true},
		Route{"CreateNode", "POST", "/nodes", CreateNode, true},
		Route{"UpdateNode", "PUT", "/nodes", UpdateNode, true},
		Route{"FetchNodes", "GET", "/nodes", FetchNodes, true},
		Route{"FetchNode", "GET", "/nodes/{id}", FetchNode, true},
		Route{"DeleteNode", "DELETE", "/nodes/{id}", DestroyNode, true},
	}
}

// FetchNodes - returns the latest collection
func FetchNodes(w http.ResponseWriter, r *http.Request) {

	nodes, _ := RepositoryGetAllNodes()

	if len(nodes) > 0 {
		Json(w, http.StatusOK, nodes)
		return
	}

	JsonError(w, http.StatusInternalServerError, "failed to fetch nodes")
}

// FetchNode - returns a single node by id
func FetchNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params id")
		return
	}

	node, err := RepositoryFindNodeById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to fetch node with  id "+id)
		return
	}

	Json(w, http.StatusOK, node)
}

// DestroyNode - deletes a node by ID and returns the latest collection
func DestroyNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string

	id, _ = vars["id"]

	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse query params [id]")
		return
	}

	nodes, err := RepositoryDestroyNodeById(id)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to delete node with  id "+id)
		return
	}

	Json(w, http.StatusOK, nodes)
}

// CreateNode - creates a node and returns the latest collection
func CreateNode(w http.ResponseWriter, r *http.Request) {
	node := &dbcontext.Node{}

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [CreateNode]")
		return
	}

	tid := payload["TunnelID"].(string)
	if len(tid) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelID")
		return
	}
	node.TunnelID = &tid

	ip := payload["IpAddress"].(string)
	if len(ip) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse IpAddress")
		return
	}

	node.IpAddress = &ip

	region := payload["Region"].(string)
	if len(region) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Region")
		return
	}

	node.Region = &region

	provider := payload["Provider"].(string)
	if len(provider) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse tunnel id")
		return
	}

	node.Provider = &provider

	nodes, err := RepositoryCreateNode(node)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new node")
		return
	}

	Json(w, http.StatusOK, nodes)
}

// UpdateNode - updates a node and returns the latest collection
func UpdateNode(w http.ResponseWriter, r *http.Request) {

	payload, err := ReadBody(r)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to parse body [UpdateNode]")
		return
	}

	id := payload["NodeID"].(string)
	if len(id) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse NodeID")
		return
	}

	node, err := RepositoryFindNodeById(id)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error: the id you provided was not found in our system. there is nothing to update at this time")
		return
	}

	tid := payload["TunnelID"].(string)
	if len(tid) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse TunnelID")
		return
	}
	node.TunnelID = &tid

	ip := payload["IpAddress"].(string)
	if len(ip) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse IpAddress")
		return
	}

	node.IpAddress = &ip

	region := payload["Region"].(string)
	if len(region) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse Region")
		return
	}

	node.Region = &region

	provider := payload["Provider"].(string)
	if len(provider) == 0 {
		JsonError(w, http.StatusInternalServerError, "failed to parse tunnel id")
		return
	}

	node.Provider = &provider

	nodes, err := RepositoryUpdateNode(&node)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "failed to create a new node")
		return
	}

	Json(w, http.StatusOK, nodes)
}
