package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"fmt"
	"os"
	"strings"
)

type RecommendationList struct {
	Value []Value `json:"value"`
}
type Properties struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type SupportedValues struct {
	RecommendationCategory string       `json:"recommendationCategory"`
	RecommendationImpact   string       `json:"recommendationImpact"`
	SupportedResourceType  string       `json:"supportedResourceType"`
	ID                     string       `json:"id"`
	DisplayName            string       `json:"displayName"`
	Properties             []Properties `json:"properties"`
}
type Properties struct {
	DisplayName         string            `json:"displayName"`
	DependsOn           []string          `json:"dependsOn"`
	ApplicableScenarios []string          `json:"applicableScenarios"`
	SupportedValues     []SupportedValues `json:"supportedValues"`
}
type SupportedValues struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}
type Properties struct {
	DisplayName         string            `json:"displayName"`
	ApplicableScenarios []string          `json:"applicableScenarios"`
	SupportedValues     []SupportedValues `json:"supportedValues"`
}
type Properties struct {
	DisplayName     string            `json:"displayName"`
	SupportedValues []SupportedValues `json:"supportedValues"`
}
type Properties struct {
	DisplayName         string            `json:"displayName"`
	ApplicableScenarios []string          `json:"applicableScenarios"`
	SupportedValues     []SupportedValues `json:"supportedValues"`
}
type Value struct {
	Properties Properties `json:"properties,omitempty"`
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	Properties Properties `json:"properties,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Properties Properties `json:"properties,omitempty"`
	Properties Properties `json:"properties,omitempty"`
}

func init() {

	r, err := setRecommendationListCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setRecommendationListCommand() (*cobra.Command, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.RecommendationList.DescriptionLine1,
		cl.AppConfig.RecommendationList.DescriptionLine2,
		cl.AppConfig.RecommendationList.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.RecommendationList.Command,
		Short: cl.AppConfig.RecommendationList.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		r := &RecommendationList{}
		r, err := r.getAzureRecommendationList()
		if err != nil {
			return err
		}

		clearTerminal()
		r.Print()
		return nil
	}
	return cmd, nil
}

func (r *RecommendationList) getAzureRecommendationList() (*RecommendationList, error) {
	var at = &AccessToken{}
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	at, err = at.getAccessToken()
	if err != nil {
		return nil, err
	}

	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	client := &http.Client {}
	req, _ := http.NewRequest("GET", cl.AppConfig.RecommendationList.URL, nil)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//fmt.Println(string(body))

	err = json.Unmarshal(body, r)
	if err != nil {
		fmt.Println("recommendation list unmarshal body response: ", err)
	}

	return r, nil
}

func (r *RecommendationList) Print() {
	fmt.Println("Azure Recommendation List:")
	fmt.Println("----------------------------------------")

	for i := 0; i < len(r.Value); i++ {
		recommendaiton := r.Value[i]
		for x := 0; x < len(recommendaiton.Properties.)
	}
}
