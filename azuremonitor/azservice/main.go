package main

import (
	"fmt"
	c "github.com/Go/azuremonitor/config"
	"io/ioutil"
	"net/http"

	"strings"
	//"github.com/Go/azuremonitor/azservice/cmd"
	//"io/ioutil"
	//"net/http"
	"time"
)

type Request struct {
		Name    string
		Url     string
		Method  string
		Payload string
}
var configuration c.CmdConfig

func main() {
	start := time.Now()
	configuration, _ = c.GetCmdConfig()
	url := strings.Replace(configuration.AccessToken.URL, "{{tenantID}}", configuration.AccessToken.TenantID, 1)
	strPayload := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&scope=%s",configuration.AccessToken.GrantType,configuration.AccessToken.ClientID,configuration.AccessToken.ClientSecret,configuration.AccessToken.Scope)
	fmt.Printf("url: %s\n", url)
	var requests = []Request{
		{"accesstoken", url, "POST", strPayload},
		{"google", "https://www.google.com", "GET", ""},
		{"msn", "https://www.msn.com", "GET", ""},
	}

	ch := make(chan string)
	for _, request := range requests {
		go MakeRequest(request, ch)
	}

	for range requests{
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	//cmd.Execute()
}

func MakeRequest(request Request, ch chan<-string) {
	start := time.Now()
	secs := time.Since(start).Seconds()
	fmt.Printf("request url: %s\n", request.Url)

	payload := strings.NewReader(request.Payload)

	client := &http.Client{}
	req, err := http.NewRequest("POST", request.Url, payload)
	if err != nil {
		fmt.Println(err)
	}

	//need header
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//err = json.Unmarshal(body, at)
	//if err != nil {
	//	fmt.Println("unmarshal body response: ", err)
	//}

	ch <- fmt.Sprintf("%.2f elapsed with response length: %s", secs,string(body))
}
