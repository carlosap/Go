package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"strings"
)

type Weather struct {
	MsgType string  `json:"msgtype"`
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   int     `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func init() {

	weatherApi, err := setWeatherCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(weatherApi)
}

func setWeatherCommand() (*cobra.Command, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.Weather.DescriptionLine1,
		cl.AppConfig.Weather.DescriptionLine2,
		cl.AppConfig.Weather.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.Weather.Command,
		Short: cl.AppConfig.Weather.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		w := &Weather{}
		ip, err := IpInfo.getIpInfo(false)
		if err != nil {
			return err
		}

		w, err = w.getWeather(ip)
		if err != nil {
			return err
		}

		clearTerminal()
		w.Print()

		return nil
	}
	return cmd, nil
}

func (w *Weather) getWeather(ip *IpapiResponse) (*Weather, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	w, err = cl.getWeather(ip)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Weather) Print() {
	fmt.Printf(
		`
Current Weather:
--------------------------------------
`)
	temp := fmt.Sprintf("%f", w.Main.Temp)
	pressure := fmt.Sprintf("%d", w.Main.Pressure)
	humidity := fmt.Sprintf("%d", w.Main.Humidity)
	speed := fmt.Sprintf("%f", w.Wind.Speed)
	speedDeg := fmt.Sprintf("%d", w.Wind.Deg)
	description := fmt.Sprintf("%s. wind speed %s at %s deg", w.Weather[0].Description, speed, speedDeg)
	fmt.Printf("Temp: %s\t Pressure: %s\t Humidity: %s\t Description: %s\n",
		temp,
		pressure,
		humidity,
		description)

}

func (w *Weather) doWeatherRequest(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	bytes, err := executeRequest(req, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &w)
	if err != nil {
		return nil //err
	}

	return nil
}

//-------------------Weather Client Response ----------------------------------
var broadcastWeather = make(chan *Weather)

func doBroadcastWeather() {
	for {
		val := <-broadcastWeather
		for client := range clients {
			strWeather, err := json.Marshal(val)
			if err != nil {
				client.Close()
				delete(clients, client)
			}

			err = client.WriteMessage(websocket.TextMessage, strWeather)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	var weather Weather
	if err := json.NewDecoder(r.Body).Decode(&weather); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writeWeather(&weather)
}

func writeWeather(w *Weather) {
	broadcastWeather <- w
}

func (w *Weather) clientWeatherResponse() error {

	w.MsgType = "weather"
	url := "http://localhost:5000/weather"
	strWeather, _ := json.Marshal(w)
	body := strings.NewReader(string(strWeather))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error: failed create new request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error: failed do weather request:: %v", err)
	}
	defer resp.Body.Close()
	return nil
}
