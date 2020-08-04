package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

type AppConfig struct {
	IP struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"ip"`
	AccessToken struct {
	Command          string `json:"command"`
	CommandComments  string `json:"command_comments"`
	DescriptionLine1 string `json:"description_line1"`
	DescriptionLine2 string `json:"description_line2"`
	DescriptionLine3 string `json:"description_line3"`
	Key              string `json:"key"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	GrantType        string `json:"grant_type"`
	ClientID         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	Scope            string `json:"scope"`
	SubscriptionID   string `json:"subscription_id"`
	TenantID         string `json:"tenant_id"`
} `json:"access_token"`
	Resources struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		GrantType        string `json:"grant_type"`
		ClientID         string `json:"client_id"`
		ClientSecret     string `json:"client_secret"`
		Scope            string `json:"scope"`
		SubscriptionID   string `json:"subscription_id"`
		TenantID         string `json:"tenant_id"`
	} `json:"resources"`
	SubscriptionInfo struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
		SubscriptionID   string `json:"subscription_id"`
	} `json:"subscriptioninfo"`
	RecommendationList struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"recommendationlist"`
	Recommendation struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"recommendation"`
	ResourceGroups struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"resourcegroups"`
	Weather struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"weather"`
	Forecast struct {
		Command          string `json:"command"`
		CommandComments  string `json:"command_comments"`
		DescriptionLine1 string `json:"description_line1"`
		DescriptionLine2 string `json:"description_line2"`
		DescriptionLine3 string `json:"description_line3"`
		Key              string `json:"key"`
		Name             string `json:"name"`
		URL              string `json:"url"`
	} `json:"forecast"`
}

// ReadConfig reads the file of the filename
func ReadConfig(filename string) (*AppConfig, error) {

	var ac AppConfig
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read error")
	}
	err = json.Unmarshal([]byte(file), &ac)

	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	//fmt.Println("appConfig:", ac)
	return &ac, nil
}
