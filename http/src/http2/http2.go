package main

import (
	"net/http"
	"golang.org/x/net/http2"
)

func main() {
	server := new(http.Server)
	h2Config := new(http2.Server)
	http2.ConfigureServer(server, h2Config)
	http.Handle("/", http.FileServer(http.Dir("public")))
	server.Addr = ":3001"
	http.ListenAndServeTLS(":3001", "cert.pem", "key.pem", nil)


}