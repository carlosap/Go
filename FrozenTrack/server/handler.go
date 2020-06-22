package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

var PollerRunning = false

// must is helper function to save redundant error expressions/validations
func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Json helper function to return JSON content
func Json(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	must(json.NewEncoder(w).Encode(data))
}

func JsonError(w http.ResponseWriter, statusCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	must(json.NewEncoder(w).Encode(JErr{Code: statusCode, Text: errorMsg}))
}

// JError struct to return errors with code and custom text
type JErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// ReadBody to read body of request
func ReadBody(r *http.Request) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read http body")
	}

	values := make(map[string]interface{})
	if err = json.Unmarshal(body, &values); err != nil {
		return nil, errors.Wrap(err, "failed to parse json http body")
	}

	return values, nil
}

/*
************** CPEREZ BACK FROM OLD FILES- WE WILL KEEP SOME OF THESE CODES UNTIL WE SURE WE KNOW WHAT WE NEED **************
 */

//	/* endpoints for UI to connect to  */
//	/* these are test endpoints */
//	router.HandleFunc("/addN", addNode)
//	router.HandleFunc("/load", db.loadConfig).Methods("OPTIONS", "GET", "POST")
//	router.HandleFunc("/save", saveConfig)
//
//	router.HandleFunc("/destroy", closeConnection)
//	router.HandleFunc("/networkinfo", networkInfo)
//	router.HandleFunc("/internettype", internetType)
//	router.HandleFunc("/status", checkipStatus)
//	router.HandleFunc("/downloadPackage/{package}", downloadPackage)
//

//func init_UI() http.Handler {
//	static := http.FileServer(http.Dir(getDirectory()))
//	fmt.Println("server directory:" + getDirectory())
//	fmt.Println("Server Running on port 5555")
//	return static
//}

/*
	this function gets the regions and zones array of objects
*/

//func getAvailableLocations(w http.ResponseWriter, r *http.Request) {
//	locations := getRegionsAndZones()
//	jsonify, err := json.Marshal(locations)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	/* Set the content type that your are sending to the response writer */
//	w.Header().Set("Content-Type", "application/json")
//
//	// will receive json format of the regions and zones array on front end
//	w.Write(jsonify)
//}

//func setUserName(w http.ResponseWriter, r *http.Request) {
//	un := r.Header.Get("Username") //convert to r.body()
//	setUsername(un)
//	getAvailableLocations(w, r)
//}
//
//func closeConnection(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("destory button clicked\n shutting down server...")
//	os.Exit(0)
//}

//func networkInfo(w http.ResponseWriter, r *http.Request) { // Queues browserleaks to pop up so user can check PoP
//	var cmd5 = exec.Command("firefox", "-private-window", "https://www.browserleaks.com/ip")
//	go cmd5.Run()
//	if !PollerRunning {
//		go PollStatus()
//		PollerRunning = true
//	}
//}
//
//func internetType(w http.ResponseWriter, r *http.Request) { // Checks the kind of internet type and if it is working correctly
//	ifStatusW := CheckInterfaceStatus("wlan0")
//	ifStatusE := CheckInterfaceStatus("eth0")
//
//	if !ifStatusW && !ifStatusE {
//		fmt.Println("Disruption in connection detected, shutting the pipeline down!")
//		closeConnection(w, r)
//	}
//}

//func checkipStatus(w http.ResponseWriter, r *http.Request) {
//	var cmd1 = exec.Command("ipsec", "status")
//	out1, err1 := cmd1.Output()
//	if err1 != nil {
//		fmt.Println(err1)
//	}
//	if strings.Contains(string(out1), "(0 up, 0 connecting)") || strings.Contains(string(out1), "(0 up, 1 connecting)") {
//		fmt.Println("Disruption in connection detected, shutting the pipeline down!")
//		closeConnection(w, r)
//	}
//
//	cmd2 := exec.Command("curl", "ipinfo.io/ip")
//	out2, err2 := cmd2.Output()
//	if err2 != nil {
//		fmt.Println(err2)
//	}
//	ipCheck := false
//	for _, pop := range Selections.Option {
//		if strings.Contains(string(out2), pop.PoPIP) {
//			ipCheck = true
//			break
//		}
//	}
//	if !ipCheck {
//		fmt.Println("Disruption in connection detected, shutting the pipeline down!")
//		closeConnection(w, r)
//	}
//}

//type FileUpload struct {
//	QQuuid          string `json:"qquuid"`
//	QQfilename      string `json:"qqfilename"`
//	QQtotalfilesize string `json:"qqtotalfilesize"`
//	Success         bool   `json:"success"`
//}

//func addTunnel(w http.ResponseWriter, r *http.Request) {
//
//	fmt.Println("add tunnel button clicked")
//}
//
//func addNode(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("add node button clicked")
//}

//func (db Database) loadConfig(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "OPTIONS" {
//		fmt.Println("OPTIONS")
//		w.Header().Set("Content-Type", "application/json")
//	} else {
//		fmt.Println("POST")
//
//		file := FileUpload{
//			QQuuid:          r.FormValue("qquuid"),
//			QQfilename:      r.FormValue("qqfilename"),
//			QQtotalfilesize: r.FormValue("qqtotalfilesize"),
//			Success:         false,
//		}
//
//		fileptr, handler, err := r.FormFile("qqfile")
//		if err != nil {
//			fmt.Println("Error getting file")
//			return
//		}
//		defer fileptr.Close()
//
//		tempFile, err := ioutil.TempFile("configs", "configTest.toml")
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		defer tempFile.Close()
//
//		fileBytes, err := ioutil.ReadAll(fileptr)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		tempFile.Write(fileBytes)
//		fmt.Println("FILE UPLOADED")
//
//		file.Success = true
//
//		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
//		fmt.Printf("File Size: %+v\n", handler.Size)
//		fmt.Printf("MIME Header: %+v\n", handler.Header)
//
//		jsonFile, err := json.Marshal(file)
//		if err != nil {
//			panic(err)
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		w.Write(jsonFile)
//	}
//}

//func downloadPackage(w http.ResponseWriter, r *http.Request) {
//	fileName := mux.Vars(r)
//	//log.Println("GET REQUEST MADE: REQUESTING FILE -> ", fileName["package"])
//
//	filePath := "Packages/" + fileName["package"]
//	filePtr, err := os.Open(filePath)
//	defer filePtr.Close()
//	if err != nil {
//		http.Error(w, "File not found", 404)
//	}
//
//	FileHeader := make([]byte, 512)
//	//Copy the headers into the FileHeader buffer
//	filePtr.Read(FileHeader)
//	//Get content type of file
//	FileContentType := http.DetectContentType(FileHeader)
//
//	//Get the file size
//	FileStat, _ := filePtr.Stat()                      //Get info from file
//	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
//
//	//Send the headers
//	w.Header().Set("Content-Disposition", "attachment; filename="+fileName["package"])
//	w.Header().Set("Content-Type", FileContentType)
//	w.Header().Set("Content-Length", FileSize)
//
//	//Send the file
//	//We read 512 bytes from the file already, so we reset the offset back to 0
//	filePtr.Seek(0, 0)
//	io.Copy(w, filePtr) //'Copy' the file to the client
//	return
//}

//func saveConfig(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("save config button clicked")
//}
