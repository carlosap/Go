package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type MessageToken struct {
	Msg string `json:"msg"`
}
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan *MessageToken)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	HandshakeTimeout:  0,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	WriteBufferPool:   nil,
	Subprotocols:      nil,
	Error:             nil,
	EnableCompression: false,
}

func startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", homePageHandle).Methods("GET")
	router.HandleFunc("/messanger", msgHandler).Methods("POST")
	//router.HandleFunc("/weather", weatherHandler).Methods("POST")
	//router.HandleFunc("/forecast", forecastHandler).Methods("POST")
	//router.HandleFunc("/ip", ipInfoHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandle)
	go doBroadcast()
	//go doBroadcastWeather()
	//go doBroadcastForecast()
	//go doBroadcastIpInfo()
	log.Fatal(http.ListenAndServe(":5000", router))
}

func homePageHandle(w http.ResponseWriter, r *http.Request) {
	//funciton there you read the index.html
	//take that out string
	fmt.Fprintf(w, "<strong>Home page</strong>")
}

func writer(msg *MessageToken) {
	broadcast <- msg
}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	clients[ws] = true
	//echo := &MessageToken{Msg: "starting echo communication...."}
	//go writer(echo)
}

func doBroadcast() {
	for {
		val := <-broadcast
		// send to every client that is currently connected
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(val.Msg))
			if err != nil {
				log.Printf("Websocket error: %s", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func msgHandler(w http.ResponseWriter, r *http.Request) {
	var msg MessageToken
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writer(&msg)
}
