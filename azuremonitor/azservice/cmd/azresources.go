package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
	"fmt"
)

type Resource struct {
	Values []Value `json:"value"`
}

type Value struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Kind     string `json:"kind,omitempty"`
	Location string `json:"location"`
	ManagedBy string `json:"managedBy,omitempty"`
	Sku  Sku  `json:"sku,omitempty"`
	Tags Tags `json:"tags,omitempty"`
	Plan Plan `json:"plan,omitempty"`
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

type Tags     struct {
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
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.Resources.DescriptionLine1,
		cl.AppConfig.Resources.DescriptionLine2,
		cl.AppConfig.Resources.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.Resources.Command,
		Short: cl.AppConfig.Resources.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		r := &Resource{}
		r, err := r.getResources()
		if err != nil {
			return err
		}

		clearTerminal()
		r.Print()
		return nil
	}
	return cmd, nil
}

func (r *Resource) getResources() (*Resource, error) {
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

	url := strings.Replace(cl.AppConfig.Resources.URL, "{{subscriptionID}}", cl.AppConfig.AccessToken.SubscriptionID, 1)
	strheaderToken := fmt.Sprintf("Bearer %s", at.AccessToken)
	client := &http.Client {}
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("request : ", err)
	}

	req.Header.Add("Authorization", strheaderToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	//fmt.Println(string(body))

	err = json.Unmarshal(body, r)
	if err != nil {
		fmt.Println("resources unmarshal body response: ", err)
	}

	return r, nil
}

func (r *Resource) Print() {

	fmt.Printf(
		`
Azure Resources:
--------------------------------------
%v
`,r.Values,)

}

