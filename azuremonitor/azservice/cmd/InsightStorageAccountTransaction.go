package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type StorageAccountTransaction struct {
	Cost     int       `json:"cost"`
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
				Total     float64   `json:"total"`
			} `json:"data"`
		} `json:"timeseries"`
	} `json:"value"`
}

func (r *StorageAccountTransaction) getStorageAccountTransaction(resurceGroup string, storageAccount string, startD string,endD string) (*StorageAccountTransaction, error) {
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
		r, err := r.executeRequest(cl.AppConfig.AccessToken.SubscriptionID, resurceGroup, storageAccount, startD, endD, cKey)
		if err != nil {
			return r, err
		}

	} else {
		//Load From Cache
		err := LoadFromCache(cKey, r)
		if err != nil {
			r, err := r.executeRequest(cl.AppConfig.AccessToken.SubscriptionID,resurceGroup, storageAccount, startD, endD, cKey)
			if err != nil {
				return r, err
			}
		}
	}

	return r, nil
}

func (r *StorageAccountTransaction) executeRequest(subscriptionId string, resourceGroup string, storageAccount string, startD string, endD string, cKey string) (*StorageAccountTransaction, error) {

	var at = &AccessToken{}
	at, err := at.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://management.azure.com/subscriptions/" +
		"%s/resourceGroups/" +
		"%s/providers/Microsoft.Storage/storageAccounts/" +
		"%s/providers/microsoft.Insights/metrics?" +
		"timespan=%sT05:43:23.526Z/%sT09:43:23.526Z&interval=FULL" +
		"&metricnames=Transactions" +
		"&aggregation=total" +
		"&metricNamespace=Microsoft.Storage/storageAccounts&validatedimensions=false&api-version=2019-07-01",
		subscriptionId,
		resourceGroup,
		storageAccount,
		startD,
		endD,
		)


	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	payload := strings.NewReader("")

	client := &http.Client {}
	req, _ := http.NewRequest("GET",url, payload)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))

	err = json.Unmarshal(body,r)
	if err != nil {
		return r, fmt.Errorf("recommendation list unmarshal body response: ", err)
	}

	//cached it
	err = saveCache(cKey, r)
	if err != nil {
		return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
	}

	//fmt.Printf("%v\n", r)
	return r, nil
}

func (r *StorageAccountTransaction) getTransactions() float64 {
	var retVal float64
	for i:=0; i < len(r.Value); i++ {
		val := r.Value[i]
		for x:=0; x < len(val.Timeseries); x++ {
			retVal = val.Timeseries[x].Data[0].Total
			return retVal
		}
	}

	return retVal
}

