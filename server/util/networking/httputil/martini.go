package httputil

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Novetta/kerbproxy/clientAuthentication"
	"github.com/Novetta/kerbproxy/kerbtypes"
	"github.com/go-martini/martini"

	mgzip "github.com/martini-contrib/gzip"
)

var httplog = log.New(os.Stdout, "", 0)

//PaginateHandler parses request uri's for our standard pagination parameters
func PaginateHandler(r *http.Request, c martini.Context) {
	r.ParseForm()
	pageStateString := r.FormValue("pageState")
	startNumString := r.FormValue("start")
	maxResultString := r.FormValue("limit")
	sortString := r.FormValue("sort")
	filterString := r.FormValue("filter")

	pageState, err := base64.StdEncoding.DecodeString(pageStateString)
	if len(pageStateString) > 0 && err != nil {
		log.Print("Error parsing Page State: ", err)
		pageState = make([]byte, 0, 0)
	}
	startNum, err := strconv.Atoi(startNumString)
	if len(startNumString) > 0 && err != nil {
		log.Print("Error parsing Start Number: ", err)
		startNum = 0
	}
	maxResults, err := strconv.Atoi(maxResultString)
	if len(maxResultString) > 0 && err != nil {
		log.Print("Error parsing Max Results: ", err)
		maxResults = 0
	}
	s := SortMethod{}
	sList := []SortMethod{s}
	err = json.Unmarshal([]byte(sortString), &sList)
	if err != nil {
		sList = make([]SortMethod, 0, 0)
	}
	fList := []SortMethod{}
	err = json.Unmarshal([]byte(filterString), &fList)
	if err != nil {
		fList = make([]SortMethod, 0, 0)
	}
	pag := &Paginate{PageState: pageState, StartNum: startNum, MaxResults: maxResults, SortList: sList, FilterList: fList}
	c.Map(pag)
}

//SortMethod handles sorting direction provided by the request
type SortMethod struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
	Value     string `json:"value"`
}

//Paginate handles paginate information provided by the request
type Paginate struct {
	PageState  []byte       `json:"pageState"`
	StartNum   int          `json:"startnum"`
	MaxResults int          `json:"maxResults"`
	SortList   []SortMethod `json:"sortList"`
	FilterList []SortMethod `json:"filterList"`
}

//QueryResultMetaData contains the page state and record count of the response
type QueryResultMetaData struct {
	PageState   string `json:"pageState"`
	RecordCount int    `json:"recordCount"`
	QueryTime   string `json:"queryTime"`
}

// PagedData is the response used for a page state query
type PagedData struct {
	MetaData QueryResultMetaData `json:"metaData"`
	Data     interface{}         `json:"data"`
}

//SetupMartini creates our standard martini instance.  It automatically handles logging of requests
//and sets security relevant headers for all responses
func SetupMartini(logging martini.Handler, webroot string, indexFile string) *martini.ClassicMartini {
	if indexFile == "" {
		indexFile = "index.html"
	}
	m := martini.New()
	m.Handlers(
		authHandler.AuthenticateHandler,
		logging,
		martini.Recovery,
	)
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(JEncoder{
			Writer: w,
		}, (*JsonEncoder)(nil))
	})
	m.Use(mgzip.All())
	m.Use(func(w http.ResponseWriter) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Pragma", "private")
		w.Header().Set("X-UA-Compatible", "IE=edge")
	})
	if len(webroot) > 0 {
		m.Use(martini.Static(webroot, martini.StaticOptions{IndexFile: indexFile, SkipLogging: true}))
	}
	r := martini.NewRouter()
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	classic := &martini.ClassicMartini{
		Martini: m,
		Router:  r,
	}

	return classic
}

//BeforeFunc allows us to modify responses before the data is finally written to the network
func BeforeFunc(w martini.ResponseWriter) {
	log.Printf("Running before")
	if len(w.Header().Get("Content-Type")) == 0 && w.Size() > 0 {
		log.Printf("Setting Content-Type header")
		w.Header().Set("Content-Type", "application/json")
	} else {
		log.Printf("Not setting header, %d, %d", len(w.Header().Get("Content-Type")), w.Size())
	}
}

//OldLogging handles logging to a specific logger
func OldLogging(r *http.Request, l *log.Logger, user *kerbtypes.User, w http.ResponseWriter, c martini.Context) {
	c.Next()
	remoteIP := r.RemoteAddr
	if f := r.Header.Get("X-Forwarded-For"); len(f) > 0 {
		remoteIP = strings.Split(f, ", ")[0]
	}
	l.Printf("%s:%s:%s:%d", remoteIP, user, r.RequestURI, w.(martini.ResponseWriter).Status())
}

//MartiniLogging handles logging requests to stdout
func MartiniLogging(r *http.Request, user *kerbtypes.User, w http.ResponseWriter, c martini.Context) {
	t := time.Now()
	c.Next()

	remoteIP := r.RemoteAddr
	if f := r.Header.Get("X-Forwarded-For"); len(f) > 0 {
		remoteIP = strings.Split(f, ", ")[0]
	}
	agent := r.UserAgent()
	if agent == "" {
		agent = "unknown"
	}
	httplog.Printf("%v ~ %s ~ %s ~ %s ~ %s ~ %d ~ %s ~ %s", time.Now().UTC().Format(time.RFC3339), remoteIP, user, r.Method, r.RequestURI, w.(martini.ResponseWriter).Status(), time.Since(t), agent)
}
