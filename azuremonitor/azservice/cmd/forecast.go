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

type ForeCast struct {
	MsgType string  `json:"msgtype"`
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		} `json:"wind"`
		Rain struct {
			ThreeH float64 `json:"3h"`
		} `json:"rain,omitempty"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

func init() {

	openWeatherApi, err := setForeCastCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(openWeatherApi)
}

func setForeCastCommand() (*cobra.Command, error) {

	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.Forecast.DescriptionLine1,
		cl.AppConfig.Forecast.DescriptionLine2,
		cl.AppConfig.Forecast.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.Forecast.Command,
		Short: cl.AppConfig.Forecast.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		f := &ForeCast{}
		ip, err := IpInfo.getIpInfo(false)
		if err != nil {
			return err
		}

		f, err = f.getForeCast(ip)
		if err != nil {
			return err
		}

		clearTerminal()
		f.Print()

		return nil
	}
	return cmd, nil
}

func (f *ForeCast) getForeCast(ip *IpapiResponse) (*ForeCast, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	f, err = cl.getForeCastByLocation(ip)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *ForeCast) Print() {
	fmt.Printf(
		`
Forcecast:
--------------------------------------
`)

	for i := 0; i < len(f.List); i++ {
		temp := fmt.Sprintf("%f", f.List[i].Main.Temp)
		pressure := fmt.Sprintf("%d", f.List[i].Main.Pressure)
		humidity := fmt.Sprintf("%d", f.List[i].Main.Humidity)
		speed := fmt.Sprintf("%f", f.List[i].Wind.Speed)
		speedDeg := fmt.Sprintf("%d", f.List[i].Wind.Deg)
		description := fmt.Sprintf("%s. wind speed %s at %s deg", f.List[i].Weather[0].Description, speed, speedDeg)

		if i > 7 {
			break
		}
		fmt.Printf("%s - Temp: %s\t Pressure: %s\t Humidity: %s\t Description: %s\n",
			f.List[i].DtTxt,
			temp,
			pressure,
			humidity,
			description)
	}
}

func (f *ForeCast) doForeCastRequest(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	bytes, err := executeRequest(req, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &f)
	if err != nil {
		return nil //err
	}

	return nil
}

//-------------------Forecast Client Response ----------------------------------

var broadcastForecast = make(chan *ForeCast)

func doBroadcastForecast() {
	for {
		val := <-broadcastForecast
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

func forecastHandler(w http.ResponseWriter, r *http.Request) {
	var f ForeCast
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		log.Printf("ERROR: %s", err)
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writeForecast(&f)
}

func writeForecast(f *ForeCast) {
	broadcastForecast <- f
}

func (f *ForeCast) clientForecastResponse() error {

	f.MsgType = "forecast"
	url := "http://localhost:5000/forecast"
	strForecast, _ := json.Marshal(f)
	body := strings.NewReader(string(strForecast))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error: failed create new request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error: failed do forecast request:: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
