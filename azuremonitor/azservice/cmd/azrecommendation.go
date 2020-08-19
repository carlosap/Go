package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
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

func init() {

	r, err := setRecommendationCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setRecommendationCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.Recommendation.DescriptionLine1,
		configuration.Recommendation.DescriptionLine2,
		configuration.Recommendation.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.Recommendation.Command,
		Short: configuration.Recommendation.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		r := &RecommendationList{}
		r, err := r.getAzureRecommendation()
		if err != nil {
			return err
		}

		clearTerminal()
		r.PrintRecommendations()
		return nil
	}
	return cmd, nil
}

func (r *RecommendationList) getURL() string {

return 	strings.Replace(configuration.Recommendation.URL,
	"{{subscriptionID}}",
	configuration.AccessToken.SubscriptionID, 1)

}

func (r *RecommendationList) getHeader() http.Header {
	var at = &AccessToken{}
	var header = http.Header{}
	at, _ = at.getAccessToken()
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")

	return header
}

func (r *RecommendationList) getAzureRecommendation() (*RecommendationList, error) {

	request := Request{
		"RecommendationList",
		r.getURL(),
		Methods.GET,
		"",
		r.getHeader(),
		false,
		r,
	}
	_ = request.Execute()
	body := request.GetResponse()
	err := json.Unmarshal(body, r)
	if err != nil {
		fmt.Println("recommendation list unmarshal body response: ", err)
	}

	return r, nil
}

func (r *RecommendationList) PrintRecommendations() {
	fmt.Println("Subscription Recommendations:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i := 0; i < len(r.Value); i++ {
		printRecommendation(i, r.Value[i])
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
