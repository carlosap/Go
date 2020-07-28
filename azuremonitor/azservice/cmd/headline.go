package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

type HeadLineNews struct {
	MsgType      string `json:"msgtype"`
	Type         string `json:"_type"`
	ReadLink     string `json:"readLink"`
	QueryContext struct {
		OriginalQuery string `json:"originalQuery"`
	} `json:"queryContext"`
	TotalEstimatedMatches int `json:"totalEstimatedMatches"`
	Value                 []struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
		Image struct {
			Thumbnail struct {
				ContentURL string `json:"contentUrl"`
				Width      int    `json:"width"`
				Height     int    `json:"height"`
			} `json:"thumbnail"`
		} `json:"image"`
		Description string `json:"description"`
		About       []struct {
			ReadLink string `json:"readLink"`
			Name     string `json:"name"`
		} `json:"about,omitempty"`
		Provider []struct {
			Type  string `json:"_type"`
			Name  string `json:"name"`
			Image struct {
				Thumbnail struct {
					ContentURL string `json:"contentUrl"`
				} `json:"thumbnail"`
			} `json:"image"`
		} `json:"provider"`
		DatePublished time.Time `json:"datePublished"`
		Category      string    `json:"category"`
		Headline      bool      `json:"headline"`
		AmpURL        string    `json:"ampUrl,omitempty"`
		Mentions      []struct {
			Name string `json:"name"`
		} `json:"mentions,omitempty"`
	} `json:"value"`
}

func init() {

	newsApi, err := setNewsCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(newsApi)
}

func setNewsCommand() (*cobra.Command, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.News.DescriptionLine1,
		cl.AppConfig.News.DescriptionLine2,
		cl.AppConfig.News.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.News.Command,
		Short: cl.AppConfig.News.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		n := &HeadLineNews{}
		ip, err := IpInfo.getIpInfo(false)
		if err != nil {
			return err
		}

		n, err = n.getHeadLineNews(ip)
		if err != nil {
			return err
		}

		//clearTerminal()
		n.Print()

		return nil
	}
	return cmd, nil
}

func (n *HeadLineNews) getHeadLineNews(ip *IpapiResponse) (*HeadLineNews, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	n, err = cl.getHeadLineNews(ip)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *HeadLineNews) Print() {
	fmt.Printf(
		`
Current Headline:

`)

	for i := 0; i < len(n.Value); i++ {
		//dummy data to just display
		//if i > 5 {break}

		el := n.Value[i]
		fmt.Printf("-----------------------%s------------------------------------\n", el.Provider[0].Name)
		if len(el.Category) > 0 {
			fmt.Printf("Category: %s\n", el.Category)
		}

		fmt.Printf("%s - %s\n", el.DatePublished, el.Name)
		fmt.Printf("%s\n\n", el.Description)

	}

}

func (n *HeadLineNews) doHeadLineNewsRequest(url string, apiKey string, lat string, long string, ip string) error {
	opts := make(map[string]string)
	location := fmt.Sprintf("lat:%s;long:%s;re:1800", lat, long)
	opts["Ocp-Apim-Subscription-Key"] = apiKey
	opts["X-Search-Location"] = location
	opts["X-Search-ClientIP"] = ip
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	bytes, err := executeRequest(req, opts)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &n)
	if err != nil {
		fmt.Println("error: ", err)
		return nil //err
	}

	fmt.Println(n)
	return nil
}

//-------------------IpInfo Client Response ----------------------------------

var broadcastNews = make(chan *HeadLineNews)

func doBroadcastNews() {
	for {
		val := <-broadcastNews
		for client := range clients {
			strnews, err := json.Marshal(val)
			if err != nil {
				client.Close()
				delete(clients, client)
			}

			err = client.WriteMessage(websocket.TextMessage, strnews)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	var news HeadLineNews
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		http.Error(w, "Bad request", http.StatusTeapot)
		return
	}
	defer r.Body.Close()
	go writeNews(&news)
}

func writeNews(n *HeadLineNews) {
	broadcastNews <- n
}

func (n *HeadLineNews) clientNewsResponse() error {

	n.MsgType = "news"
	url := "http://localhost:5000/news"
	strNews, _ := json.Marshal(n)
	body := strings.NewReader(string(strNews))
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error: failed create new request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error: failed do news request:: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
