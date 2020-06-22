package main

import (
	"log"
	"net/http"
)

var routes Routes

func init() {
	routes = append(routes, nodeRoutes()...)
	routes = append(routes, tunnelRoutes()...)
	routes = append(routes, tunnelActionRoutes()...)
	routes = append(routes, tunnelStateRoutes()...)
	routes = append(routes, tunnelTypeRoutes()...)
	routes = append(routes, tunnelRequestRoutes()...)
}

func main() {
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":5555", router))
}
