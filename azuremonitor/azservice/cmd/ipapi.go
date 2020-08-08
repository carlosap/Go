package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// IpapiResponse - marshal data response from JSON
type IpapiResponse struct {
	MsgType string  `json:"msgtype"`
	IP            string      `json:"ip"`
	Type          string      `json:"type"`
	ContinentCode string      `json:"continent_code"`
	ContinentName string      `json:"continent_name"`
	CountryCode   string      `json:"country_code"`
	CountryName   string      `json:"country_name"`
	RegionCode    string      `json:"region_code"`
	RegionName    string      `json:"region_name"`
	City          string      `json:"city"`
	Zip           interface{} `json:"zip"`
	Latitude      float64     `json:"latitude"`
	Longitude     float64     `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEu                    bool   `json:"is_eu"`
	} `json:"location"`
}

var cacheFile = "ipinfo.json"
var IpInfo IpapiResponse

func init() {
	IpInfo = IpapiResponse{}
	ipapi, err := setIpApiCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(ipapi)
}

func setIpApiCommand() (*cobra.Command, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.IP.DescriptionLine1,
		cl.AppConfig.IP.DescriptionLine2,
		cl.AppConfig.IP.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.IP.Command,
		Short: cl.AppConfig.IP.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		IpInfo, err := IpInfo.getIpInfo(false)
		if err != nil {
			return err
		}

		clearTerminal()
		IpInfo.Print()
		return nil
	}
	return cmd, nil
}

func (ipInfo *IpapiResponse) getIpInfo(isCache bool) (*IpapiResponse, error) {

	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	if isCache {
		ipInfo = getIpInfoFromCache()
		if len(ipInfo.IP) <= 0 {
			ipInfo, err = cl.getIPLocalization()
			if err != nil {
				return nil, err
			}
		}
	} else {
		clearCache(cacheFile)
		ipInfo, err = cl.getIPLocalization()
		if err != nil {
			return nil, err
		}
	}

	if len(ipInfo.IP) <= 0 {
		return nil, fmt.Errorf("No IP Address was captured from IPAPI")
	}

	_ = saveCache(cacheFile, ipInfo)

	return ipInfo, nil
}

func (ipInfo *IpapiResponse) Print() {

	fmt.Printf(
		`
IP Information:
--------------------------------------
IP Address:               %s
Type:                     %s
Continent [code]:         %s
Continent Name:           %s
Country [Code]:           %s
Country Name:             %s
Region [Code]:            %s
Region Name:              %s
City:                     %s
Zip:                      %v
Latitude:                 %f
Longitude:                %f
Country Capital:          %s
Country Languages:        %v
Country Flag Url:         %s

`,
		ipInfo.IP,
		ipInfo.Type,
		ipInfo.ContinentCode,
		ipInfo.ContinentName,
		ipInfo.CountryCode,
		ipInfo.CountryName,
		ipInfo.RegionCode,
		ipInfo.RegionName,
		ipInfo.City,
		ipInfo.Zip,
		ipInfo.Latitude,
		ipInfo.Longitude,
		ipInfo.Location.Capital,
		ipInfo.Location.Languages,
		ipInfo.Location.CountryFlag,
	)
}

func (c *Client) getIPLocalization() (*IpapiResponse, error) {
	ipInfo := &IpapiResponse{}
	ips, err := getIP()
	if err != nil {
		fmt.Println(err)
		return ipInfo, err
	}

	for i := 0; i < len(ips); i++ {
		url := fmt.Sprintf("%s/%s?access_key=%s", c.AppConfig.IP.URL, ips[i], c.AppConfig.IP.Key)
		//fmt.Println("url: ", url)
		err = ipInfo.doDiscoveryRequest(url)
		//fmt.Printf("IP Address Type: %s - %T - sized[%d]\n", retval.Type, retval.Type, len(retval.Type))
		if len(ipInfo.Type) > 0 {
			break
		}
	}

	return ipInfo, nil
}

func (ipInfo *IpapiResponse) doDiscoveryRequest(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	bytes, err := executeRequest(req, nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &ipInfo)
	if err != nil {
		return nil //err
	}

	return nil
}



func getIpInfoFromCache() *IpapiResponse {
	data := &IpapiResponse{}
	file, _ := ioutil.ReadFile(cacheFile)

	_ = json.Unmarshal([]byte(file), data)
	return data
}



//-------------------IpInfo Client Response ----------------------------------

var broadcastIpInfo = make(chan *IpapiResponse)

func doBroadcastIpInfo() {
	for {
		val := <-broadcastIpInfo
		for client := range clients {
			strip, err := json.Marshal(val)
			if err != nil {
				client.Close()
				delete(clients, client)
			}

			err = client.WriteMessage(websocket.TextMessage, strip)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func ipInfoHandler(w http.ResponseWriter, r *http.Request) {
	var ip IpapiResponse
	if err := json.NewDecoder(r.Body).Decode(&ip); err != nil {
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writeIpInfo(&ip)
}

func writeIpInfo(ip *IpapiResponse) {
	broadcastIpInfo <- ip
}

func (ip *IpapiResponse) clientIpInfoResponse() error {

	ip.MsgType = "ipinfo"
	url := "http://localhost:5000/ip"
	strIp, _ := json.Marshal(ip)
	body := strings.NewReader(string(strIp))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error: failed create new request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error: failed do ip request:: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
