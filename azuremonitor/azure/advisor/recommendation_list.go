package advisor

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/httpclient"
	"net/http"
)



type RecommendationList struct {
	Value []RecommendationValue `json:"value"`
}

func (rlist *RecommendationList) ExecuteRequest(r httpclient.IRequest) {

	request := httpclient.Request{
		"recommendation_list",
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

func (rlist *RecommendationList) GetUrl() string {

	return configuration.RecommendationList.URL
}

func (rlist *RecommendationList) GetMethod() string {
	return httpclient.Methods.GET
}

func (rlist *RecommendationList) GetPayload() string {
	return ""
}

func (rlist *RecommendationList) GetHeader() http.Header {

	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}

func (rlist *RecommendationList) Print() {

	fmt.Println("Azure Recommendation List:")
	fmt.Println("----------------------------------------")

	for i := 0; i < len(rlist.Value); i++ {

		recommendaiton := rlist.Value[i]
		switch recommendaiton.Properties.DisplayName {
		case "Recommendation Type":
			printRecommendationTypes(recommendaiton)
		case "Category":
			printRecommendationCategory(recommendaiton)
		case "Impact":
			printRecommendationImpact(recommendaiton)
		case "Supported Resource Type":
			printRecommendationResource(recommendaiton)
		case "Level":
			printRecommendationLevel(recommendaiton)
		case "Status":
			printRecommendationStatus(recommendaiton)
		case "Initiated By":
			printRecommendationInitiatedBy(recommendaiton)
		default:
			fmt.Printf("default: a is %s\n", recommendaiton.Properties.DisplayName)
		}
	}
}

func printRecommendationTypes(recommendaiton RecommendationValue) {
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("Category [%s] Impact [%s] Type [%s] - %s\n", v.RecommendationCategory, v.RecommendationImpact, v.SupportedResourceType, v.DisplayName)
	}
}

func printRecommendationCategory(recommendaiton RecommendationValue) {
	fmt.Println("----------------------------------------\n")
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("ID [%s] - [%s]\n", v.ID, v.DisplayName)
	}
}

func printRecommendationImpact(recommendaiton RecommendationValue) {
	fmt.Println("----------------------------------------\n")
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("ID [%s] - [%s]\n", v.ID, v.DisplayName)
	}
}

func printRecommendationResource(recommendaiton RecommendationValue) {
	fmt.Println("----------------------------------------\n")
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("[%d]ID [%s] - [%s]\n", x+1, v.ID, v.DisplayName)
	}
}

func printRecommendationLevel(recommendaiton RecommendationValue) {
	fmt.Println("----------------------------------------\n")
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("ID [%s] - [%s]\n", v.ID, v.DisplayName)
	}
}

func printRecommendationStatus(recommendaiton RecommendationValue) {
	fmt.Println("----------------------------------------\n")
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("ID [%s] - [%s]\n", v.ID, v.DisplayName)
	}
}

func printRecommendationInitiatedBy(recommendaiton RecommendationValue) {
	fmt.Println("----------------------------------------\n")
	fmt.Printf("Name: %s\n", recommendaiton.Properties.DisplayName)
	fmt.Println("----------------------------------------\n")
	for x := 0; x < len(recommendaiton.Properties.SupportedValues); x++ {
		v := recommendaiton.Properties.SupportedValues[x]
		fmt.Printf("ID [%s] - [%s]\n", v.ID, v.DisplayName)
	}
}
