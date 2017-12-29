package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Go/server/config"
	"github.com/Go/server/controllers"
	"github.com/Go/server/dbContext"
	"github.com/Go/server/util/environment"
	"github.com/Go/server/util/logging"
	"github.com/go-martini/martini"
	mgzip "github.com/martini-contrib/gzip"
)

//JsonEncoder is an interface for encoding json http responses
type JsonEncoder interface {
	Encode(v interface{}) ([]byte, error)
	EncodeResponse(status int, v interface{}) error
}

//JEncoder is the type that handles the response encoding
type JEncoder struct {
	Writer http.ResponseWriter
}


var (
	webRoot       string
	setupDir      string
	listenAddress string
)

func init() {
	getEnvironmentals()
	dbContext.InitDBConnection()
}

func main() {
	Run()
}

func getEnvironmentals() {
	//webRoot = environment.GetSet(config.WebRootEnv, "C:\\www")
	webRoot = environment.GetSet(config.WebRootEnv, "/")
	setupDir = environment.GetSet(config.SetupDirEnv, "/")
	//setupDir = environment.GetSet(config.SetupDirEnv, "C:\\www")
}

func getServer() http.Handler {
	m := Setup(logging.MartiniLogging, webRoot, "index.html")
	controllers.RegisterFileEndpoints(m)
	controllers.RegisterBatchEndpoints(m)
	controllers.RegisterStepEndpoints(m)
	controllers.RegisterCommentsEndpoints(m)
	controllers.RegisterMetricsEndpoints(m)
	controllers.RegisterXmlEndpoints(m)
	return m
}

func Run() {
	m := getServer()
	listenAddress = environment.GetSet(config.PortEnv, ":5001")
	if !strings.Contains(listenAddress, ":") {
		listenAddress = fmt.Sprintf(":%s", listenAddress)
	}
	http.Handle("/", m)
	logging.Info("Starting Server on port %s.", listenAddress)
	logging.Errorf("%+v", http.ListenAndServe(listenAddress, nil))
}

//Setup defines middle ware, dependency injections, router
func Setup(logging martini.Handler, webroot string, indexFile string) *martini.ClassicMartini {
	if indexFile == "" {
		indexFile = "index.html"
	}
	m := martini.New()
	m.Handlers(
		logging,
		martini.Recovery,
	)

	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(JEncoder{
			Writer: w,
		}, (*JsonEncoder)(nil))
	})
	// m.Use(func(c martini.Context, w http.ResponseWriter) {
	// 	c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
	// 	//---cperez-->>
	// 	w.Header().Set("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept")
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "PUT,GET,POST,DELETE,OPTIONS")
	// 	//w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// 	//<<---end
	// 	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// })

	m.Use(mgzip.All())
	m.Use(func(w http.ResponseWriter) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Pragma", "private")
		//w.Header().Set("Pragma", "no-cache")
		//w.Header().Set("Expires", "-1")
		w.Header().Set("X-UA-Compatible", "IE=edge")
	})

	if len(webroot) > 0 {
		m.Use(martini.Static(webroot, martini.StaticOptions{IndexFile: indexFile, SkipLogging: true}))
	}

	//Routers::::::
	r := martini.NewRouter()
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	classic := &martini.ClassicMartini{
		Martini: m,
		Router:  r,
	}

	return classic
}



//Encode encodes the given type as json and sets the appropriate http response headers
//for a json response
func (j JEncoder) Encode(v interface{}) ([]byte, error) {
	var data []byte
	var err error
	if v != nil {
		j.Writer.Header().Set("Content-Type", "application/json")
		data, err = json.MarshalIndent(v, "", " ")
	}
	return data, err
}

//EncodeResponse writes the response to the network in addition to Encode above
func (j JEncoder) EncodeResponse(status int, v interface{}) error {
	data, err := j.Encode(v)
	j.Writer.WriteHeader(status)
	j.Writer.Write(data)

	return err
}
