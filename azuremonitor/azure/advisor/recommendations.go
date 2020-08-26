package advisor

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/httpclient"
	c "github.com/Go/azuremonitor/config"
	"net/http"
	"strings"
	"time"
)


type ShortDescription struct {
	Problem  string `json:"problem"`
	Solution string `json:"solution"`
}

type ExtendedProperties struct {
	Location            string `json:"location"`
	VMSize              string `json:"vmSize"`
	TargetResourceCount string `json:"targetResourceCount"`
	Term                string `json:"term"`
	SavingsPercentage   string `json:"savingsPercentage"`
	ReservationType     string `json:"reservationType"`
	SavingsAmount       string `json:"savingsAmount"`
	AnnualSavingsAmount string `json:"annualSavingsAmount"`
	SavingsCurrency     string `json:"savingsCurrency"`
	Scope               string `json:"scope"`
}

type ResourceMetadata struct {
	ResourceID string `json:"resourceId"`
}
type Recommendations struct {
	Value []RecommendationValue `json:"value"`
}

type RecommendationValue struct {
	Properties Properties `json:"properties,omitempty"`
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Name       string     `json:"name"`
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
	Name                string            `json:"name"`
	Value               string            `json:"value"`
	DisplayName         string            `json:"displayName"`
	DependsOn           []string          `json:"dependsOn"`
	ApplicableScenarios []string          `json:"applicableScenarios"`
	SupportedValues     []SupportedValues `json:"supportedValues"`

	//adding additional fields to support recommendaitons by subscription
	Category             string             `json:"category"`
	Impact               string             `json:"impact"`
	ImpactedField        string             `json:"impactedField"`
	ImpactedValue        string             `json:"impactedValue"`
	LastUpdated          time.Time          `json:"lastUpdated"`
	RecommendationTypeID string             `json:"recommendationTypeId"`
	ShortDescription     ShortDescription   `json:"shortDescription"`
	ExtendedProperties   ExtendedProperties `json:"extendedProperties"`
	ResourceMetadata     ResourceMetadata   `json:"resourceMetadata"`
}

var (
	configuration    c.CmdConfig
)

func init(){
	configuration, _ = c.GetCmdConfig()
}

func (rlist *Recommendations) ExecuteRequest(r httpclient.IRequest) {

	request := httpclient.Request{
		"recommendations",
		r.GetUrl(),
		r.GetMethod(),
		r.GetPayload(),
		r.GetHeader(),
		true,
	}
	_ = request.Execute()
	body := request.GetResponse()
	err := json.Unmarshal(body, rlist)
	if err != nil {
		fmt.Println("unmarshal body response: ", err)
	}
}

func (rlist *Recommendations) GetUrl() string {

	return strings.Replace(configuration.Recommendation.URL,
		"{{subscriptionID}}",
		configuration.AccessToken.SubscriptionID, 1)
}
func (rlist *Recommendations) GetMethod() string {
	return httpclient.Methods.GET
}
func (rlist *Recommendations) GetPayload() string {
	return ""
}
func (rlist *Recommendations) GetHeader() http.Header {

	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (rlist *Recommendations) Print() {

	fmt.Println("Subscription Recommendations:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i := 0; i < len(rlist.Value); i++ {
		printRecommendation(i, rlist.Value[i])
	}

}

func printRecommendation(index int, recommendaiton RecommendationValue) {
	fmt.Printf("%d - Recommendation Type: %s\nImpact: [%s]\n", index+1, recommendaiton.Properties.Category, recommendaiton.Properties.Impact)
	fmt.Printf("Resource ID: %s\n", recommendaiton.Properties.ResourceMetadata.ResourceID)
	fmt.Printf(
		`
ID                       %s
Type                     %s
Description:             %s
Problem:                 %s
Recommendation:          %s
Additional Notes:        %v

`, recommendaiton.Properties.RecommendationTypeID,
		recommendaiton.Properties.ImpactedField,
		recommendaiton.Properties.ImpactedValue,
		recommendaiton.Properties.ShortDescription.Problem,
		recommendaiton.Properties.ShortDescription.Solution,
		recommendaiton.Properties.ExtendedProperties,
	)

	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
}
