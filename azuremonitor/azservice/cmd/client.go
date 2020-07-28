package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Go/azuremonitor/azservice/config"
)

type Client struct {
	AppConfig *config.AppConfig
}

// New - add a constructor
func (c *Client) New() error {
	var err error
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.AppConfig = config
	return nil
}

func (c *Client) getForeCastByLocation(ip *IpapiResponse) (*ForeCast, error) {
	f := &ForeCast{}
	lat := fmt.Sprintf("%f", ip.Latitude)
	lon := fmt.Sprintf("%f", ip.Longitude)
	url := fmt.Sprintf("%s/forecast?lat=%s&lon=%s&units=imperial&appid=%s", c.AppConfig.Weather.URL, lat, lon, c.AppConfig.Weather.Key)
	err := f.doForeCastRequest(url)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (c *Client) getWeather(ip *IpapiResponse) (*Weather, error) {
	w := &Weather{}
	lat := fmt.Sprintf("%f", ip.Latitude)
	lon := fmt.Sprintf("%f", ip.Longitude)
	url := fmt.Sprintf("%s/weather?lat=%s&lon=%s&units=imperial&appid=%s", c.AppConfig.Weather.URL, lat, lon, c.AppConfig.Weather.Key)
	err := w.doWeatherRequest(url)
	if err != nil {
		return w, err
	}

	return w, nil
}


// loads the application configurations
func loadConfig(filename string) (*config.AppConfig, error) {

	ac, err := config.ReadConfig(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file")
	}
	return ac, nil
}

func executeRequest(req *http.Request, options map[string]string) ([]byte, error) {

	setHeaderOptions(options, req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("body: %s\n", string(body))
	return body, nil
}

func setHeaderOptions(options map[string]string, req *http.Request) {
	if options != nil && len(options) > 0 {
		for key, value := range options {
			if len(key) > 0 && len(value) > 0 {
				req.Header.Set(key, value)
			}
		}
	}
}
