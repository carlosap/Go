package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"strings"

	"github.com/Go/server/util/networking/httputil"
	"github.com/Go/server/util"
	"github.com/Go/server/util/logging"
)

const (
	texasWebRootEnv  = "TEXAS_WEBROOT"
	texasPortEnv     = "TEXAS_PORT"
	texasSetupDirEnv = "TEXAS_SETUP_DIR"
)

var (
	listenAddress string
	texasWebRoot  string
	texasSetupDir string
	Version string
)

func init() {
	getEnvironmentals()
}

func main() {
	
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

	server.RunServer()
}

func getEnvironmentals() {
	texasWebRoot = util.GetSetEnv(texasWebRootEnv, "/opt/texas/public/build")
	texasSetupDir = util.GetSetEnv(texasSetupDirEnv, "/opt/texas/setup")
}

func RunServer() {
	m := getServer()
	listenAddress = util.GetSetEnv(texasPortEnv, ":8686")
	if !strings.Contains(listenAddress, ":") {
		listenAddress = fmt.Sprintf(":%s", listenAddress)
	}
	http.Handle("/", m)
	logging.Info("Starting Texas Server on port %s.", listenAddress)
	logging.Errorf("%+v", http.ListenAndServe(listenAddress, nil))
}

func getServer() http.Handler {
	m := httputil.SetupMartini(httputil.MartiniLogging, texasWebRoot, "index.html")
	registerStepEndpoints(m)
	return m
}