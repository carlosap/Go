package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Novetta/common/networking/httputil"
)

var port = flag.Int("port", 8080, "Port to serve http on")
var remoteURI = flag.String("remoteURI", "", "Remote URI to reflect")
var name = flag.String("name", "", "Name for this instance of this service")

func init() {
	flag.Parse()
	if *remoteURI == "" {
		flag.Usage()
		log.Fatal("Invalid arguments!")
	}

	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix(fmt.Sprintf("Proxy %s ", *name))
}

func main() {
	httputil.HTTPProxy(*remoteURI, *port)
}
