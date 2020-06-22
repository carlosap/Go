package main

import (
	"fmt"
	"github.com/aagon00/FrozenTraceWeb/db/dbcontext"
)

// RepositoryGetAllNodes - returns all records from the table
func RepositoryGetAllNodes() (dbcontext.Nodes, error) {
	n := dbcontext.Node{}
	nodes, err := n.GetAll()
	if err != nil {
		return nodes, fmt.Errorf("error [GetAllNodes]: %v", err)
	}

	return nodes, err
}

// RepositoryFindNodeById - returns a single records queried by id
func RepositoryFindNodeById(id string) (dbcontext.Node, error) {

	n := dbcontext.Node{}
	err := n.GetNodeByID(id)
	if err != nil {
		return n, fmt.Errorf("error [FindNodeById]: %v", err)
	}

	return n, nil
}

// RepositoryDestroyNodeById - deletes a single records by id
func RepositoryDestroyNodeById(id string) (dbcontext.Nodes, error) {

	node, err := RepositoryFindNodeById(id)
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyNodeById %s ", id)
	}

	err = node.Delete()
	if err != nil {
		return nil, fmt.Errorf("error: RepositoryDestroyNodeById -> failed to delete node %v", err)
	}

	return RepositoryGetAllNodes()
}

// RepositoryCreateNode - creates a new record via payload (json)
func RepositoryCreateNode(node *dbcontext.Node) (dbcontext.Nodes, error) {

	err := node.Insert()
	if err != nil {
		return nil, fmt.Errorf("error [FindNodeById]: %v", err)
	}

	return RepositoryGetAllNodes()
}

// RepositoryUpdateNode - updates an existing record via payload (json)
func RepositoryUpdateNode(node *dbcontext.Node) (dbcontext.Nodes, error) {

	err := node.Update()
	if err != nil {
		return nil, fmt.Errorf("error [RepositoryUpdateNode]: %v", err)
	}

	return RepositoryGetAllNodes()
}
