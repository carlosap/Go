package main

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/Go/server"
	"github.com/Go/server/util/logging"
)

var (
	//Version is the build number of the app
	Version string
)

func main() {
	//Print version information
	printVersion := flag.Bool("v", false, "Print Version")
	if !flag.Parsed() {

		flag.Parse()
	} else {
		logging.Fatalf("Flags parsed unexpectedly")
	}
	if *printVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	fmt.Println("final")
}