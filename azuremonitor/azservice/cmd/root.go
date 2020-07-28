package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	externalip "github.com/glendc/go-external-ip"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

var clear map[string]func() //create a map for storing clear funcs
func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "localization services",
	Short: "Elysium Localization Services",
	Long:  ``,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// returns internal ip and public ip
// you can also get this information https://myexternalip.com/raw
func getIP() ([]string, error) {
	var ips []string
	extIp := externalip.DefaultConsensus(nil, nil)
	ipTemp, _ := extIp.ExternalIP()
	if len(ipTemp.String()) > 0 {
		ips = append(ips, ipTemp.String())
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				ips = append(ips, ip.String())
				//fmt.Printf("ip net: %s\n", ip.String())
			case *net.IPAddr:
				ip = v.IP
				//fmt.Printf("ip address: %s\n", ip.String())
				ips = append(ips, ip.String())
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ips, nil
		}
	}
	return ips, errors.New("no network connection detected")
}

func startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/", homePageHandle).Methods("GET")
	router.HandleFunc("/messanger", msgHandler).Methods("POST")
	router.HandleFunc("/weather", weatherHandler).Methods("POST")
	router.HandleFunc("/forecast", forecastHandler).Methods("POST")
	router.HandleFunc("/ip", ipInfoHandler).Methods("POST")
	router.HandleFunc("/news", newsHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandle)
	go doBroadcast()
	go doBroadcastWeather()
	go doBroadcastForecast()
	go doBroadcastIpInfo()
	go doBroadcastNews()
	log.Fatal(http.ListenAndServe(":5000", router))
}

func homePageHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
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

func clearTerminal() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("wrong platform")
	}
}
