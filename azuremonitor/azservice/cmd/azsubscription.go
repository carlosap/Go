package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/httpclient"
	"github.com/Go/azuremonitor/common/terminal"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

type SubscriptionInfo struct {
	ID                   string               `json:"id"`
	AuthorizationSource  string               `json:"authorizationSource"`
	ManagedByTenants     []interface{}        `json:"managedByTenants"`
	SubscriptionID       string               `json:"subscriptionId"`
	TenantID             string               `json:"tenantId"`
	DisplayName          string               `json:"displayName"`
	State                string               `json:"state"`
	SubscriptionPolicies SubscriptionPolicies `json:"subscriptionPolicies"`
}
type SubscriptionPolicies struct {
	LocationPlacementID string `json:"locationPlacementId"`
	QuotaID             string `json:"quotaId"`
	SpendingLimit       string `json:"spendingLimit"`
}

func init() {

	r, err := setSubscriptionInfoCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}
func setSubscriptionInfoCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.SubscriptionInfo.DescriptionLine1,
		configuration.SubscriptionInfo.DescriptionLine2,
		configuration.SubscriptionInfo.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.SubscriptionInfo.Command,
		Short: configuration.SubscriptionInfo.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		s := &SubscriptionInfo{}
		terminal.Clear()
		request := httpclient.Request{
			Name:    "subscriptionInfo",
			Url:     s.getUrl(),
			Method: httpclient.Methods.GET,
			Payload: "",
			Header:  s.getHeader(),
			IsCache: true,
		}
		errors := request.Execute()
		IfErrorsPrintThem(errors)

		body := request.GetResponse()
		_ = json.Unmarshal(body, s)
		s.Print()
		return nil
	}
	return cmd, nil
}
func (r *SubscriptionInfo) getHeader() http.Header {
	at := &oauth2.AccessToken{}
	at.ExecuteRequest(at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (r *SubscriptionInfo) getUrl() string {
	url := strings.Replace(configuration.SubscriptionInfo.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	return url
}
func (s *SubscriptionInfo) Print() {

	fmt.Printf(
		`
--------------------------------------
Subscription Inforamtion: %s
--------------------------------------
Name:                     %s
Authorization Source:     %s
Manage By Tenants:        %v
Status:                   %s
Policies:                 %v

`, s.SubscriptionID,
		s.DisplayName,
		s.AuthorizationSource,
		s.ManagedByTenants,
		s.State,
		s.SubscriptionPolicies,
	)

}
