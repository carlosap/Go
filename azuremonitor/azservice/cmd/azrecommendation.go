package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
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

func (r *RecommendationList) getAzureRecommendation() (*RecommendationList, error) {
	var at = &AccessToken{}

	at, err := at.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := strings.Replace(configuration.Recommendation.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, r)
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
