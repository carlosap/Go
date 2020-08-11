package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ResourceUsageStorageAccount struct {
	Responses []struct {
		Content struct {
			Cost     int       `json:"cost"`
			Timespan time.Time `json:"timespan"`
			Interval string    `json:"interval"`
			Value    []struct {
				ID   string `json:"id"`
				Type string `json:"type"`
				Name struct {
					Value          string `json:"value"`
					LocalizedValue string `json:"localizedValue"`
				} `json:"name"`
				DisplayDescription string `json:"displayDescription"`
				Unit               string `json:"unit"`
				Timeseries         []struct {
					Metadatavalues []interface{} `json:"metadatavalues"`
					Data           []struct {
						TimeStamp time.Time `json:"timeStamp"`
						Total     float64   `json:"total"`
					} `json:"data"`
				} `json:"timeseries"`
				ErrorCode string `json:"errorCode"`
			} `json:"value"`
		} `json:"content"`
	} `json:"responses"`
}

func (r *ResourceUsageStorageAccount) getStorageAccountByResourceId(storageAccount string, startD string,endD string) (*ResourceUsageStorageAccount, error) {
	//Validate
	if storageAccount == "" || startD == "" || endD == "" {
		return nil, fmt.Errorf("resource id name is required")
	}

	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	//Cache lookup
	c := &Cache{}
	cKey := fmt.Sprintf("%s_%s_GetStorageAccountByResourceId_%s_%s",cl.AppConfig.AccessToken.SubscriptionID, storageAccount, startD, endD)
	cHashVal := c.Get(cKey)
	if len(cHashVal) <= 0 {
		//Execute Request
		r, err := r.executeRequest(storageAccount, startD, endD, cKey, cl.AppConfig.AccessToken.SubscriptionID)
		if err != nil {
			return r, err
		}

	} else {
		//Load From Cache
		err := LoadFromCache(cKey, r)
		if err != nil {
			r, err := r.executeRequest(storageAccount, startD, endD, cKey,cl.AppConfig.AccessToken.SubscriptionID)
			if err != nil {
				return r, err
			}
		}
		//fmt.Println(r)
	}

	return r, nil
}

func (r *ResourceUsageStorageAccount) executeRequest(storageAccount string, startD string, endD string, cKey string, subscriptionId string) (*ResourceUsageStorageAccount, error) {

	var at = &AccessToken{}
	at, err := at.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := "https://management.azure.com/batch?api-version=2015-11-01"
	token := fmt.Sprintf("Bearer %s", at.AccessToken)

	fmt.Printf("Request: %s - %s - [%s/%s]\n", subscriptionId, storageAccount, startD, endD)
	payload := strings.NewReader(fmt.Sprintf("{\"requests\": [{\"httpMethod\": \"GET\",\"name\": \"f16292f0-f5b1-4162-bd72-1ff2bc6391cb\"," +
		"\"requestHeaderDetails\": {\"commandName\": \"fx.\"},\"url\": \"https://management.azure.com/subscriptions/" +
		"%s/resourceGroups/cloud-shell-storage-eastus/providers/Microsoft.Storage/storageAccounts/" +
		"%s/providers/microsoft.Insights/metrics?" +
		"interval=FULL&metricnames=Egress" +
		"&aggregation=total&metricNamespace=microsoft.storage/storageaccounts&" +
		"validatedimensions=false&api-version=2019-07-01\"}]}", subscriptionId, storageAccount))

	client := &http.Client {}
	req, _ := http.NewRequest("POST",url, payload)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	err = json.Unmarshal(body,r)
	if err != nil {
		return r, fmt.Errorf("recommendation list unmarshal body response: ", err)
	}

	//cached it
	err = saveCache(cKey, r)
	if err != nil {
		return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
	}

	fmt.Printf("%v\n", r)
	return r, nil
}

func (r ResourceUsageStorageAccount) Print() {


}


