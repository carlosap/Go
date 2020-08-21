package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

type Resource struct {
	Values []Value `json:"value"`
}

type Value struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Kind      string `json:"kind,omitempty"`
	Location  string `json:"location"`
	ManagedBy string `json:"managedBy,omitempty"`
	Sku       Sku    `json:"sku,omitempty"`
	Tags      Tags   `json:"tags,omitempty"`
	Plan      Plan   `json:"plan,omitempty"`
}
type Plan struct {
	Name          string `json:"name"`
	PromotionCode string `json:"promotionCode"`
	Product       string `json:"product"`
	Publisher     string `json:"publisher"`
}
type Sku struct {
	Name string `json:"name"`
	Tier string `json:"tier"`
}

type Tags struct {
	MsResourceUsage string `json:"ms-resource-usage"`
}

func init() {

	r, err := setResourcesCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourcesCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.Resources.DescriptionLine1,
		configuration.Resources.DescriptionLine2,
		configuration.Resources.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.Resources.Command,
		Short: configuration.Resources.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		r := &Resource{}

		clearTerminal()
		request := Request{
			Name:      "resources",
			Url:       r.getUrl(),
			Method:    Methods.GET,
			Payload:   "",
			Header:    r.getHeader(),
			IsCache:   false,
			ValueType: r,
		}
		errors := request.Execute()
		IfErrorsPrintThem(errors)

		body := request.GetResponse()
		_ = json.Unmarshal(body, r)
		r.Print()
		return nil
	}
	return cmd, nil
}

func (r *Resource) getHeader() http.Header {
	var at = &AccessToken{}
	var header = http.Header{}
	at, err := at.getAccessToken()
	if err != nil {
		return nil
	}
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (r *Resource) getUrl() string {
	url := strings.Replace(configuration.Resources.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	return url
}
func (r *Resource) Print() {
	fmt.Println("Resource Report:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Name,Type,Kind,Location,ManageBy,Sku Name, Sku Tier,Tags,Plan Name, Plan Promotion Code, Plan Product, Plan Publisher")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i :=0; i< len(r.Values); i++ {
		var resourceType, resourceManageby string
		item := r.Values[i]

		//remove path
		if strings.Contains(item.Type, "/") {
			pArray := strings.Split(item.Type, "/")
			resourceType = pArray[len(pArray)-1]
		}

		if strings.Contains(item.ManagedBy, "/") {
			pArray := strings.Split(item.ManagedBy, "/")
			resourceManageby = pArray[len(pArray)-1]
		}


		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",item.Name, resourceType, item.Kind,item.Location,resourceManageby,
			item.Sku.Name, item.Sku.Tier,item.Tags.MsResourceUsage, item.Plan.Name,
			item.Plan.PromotionCode,item.Plan.Product, item.Plan.Publisher)
	}

}
