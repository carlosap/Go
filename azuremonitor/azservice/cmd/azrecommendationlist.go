package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type RecommendationList struct {
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

func init() {

	r, err := setRecommendationListCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setRecommendationListCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.RecommendationList.DescriptionLine1,
		configuration.RecommendationList.DescriptionLine2,
		configuration.RecommendationList.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.RecommendationList.Command,
		Short: configuration.RecommendationList.CommandComments,
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
	request := Request{
		"RecommendationList_RL",
		configuration.RecommendationList.URL,
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

func (r *RecommendationList) Print() {
	fmt.Println("Azure Recommendation List:")
	fmt.Println("----------------------------------------")

	for i := 0; i < len(r.Value); i++ {

		recommendaiton := r.Value[i]
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
		fmt.Printf("[%d]ID [%s] - [%s]\n",x+1, v.ID, v.DisplayName)
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
