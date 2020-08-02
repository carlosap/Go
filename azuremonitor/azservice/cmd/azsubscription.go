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
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.SubscriptionInfo.DescriptionLine1,
		cl.AppConfig.SubscriptionInfo.DescriptionLine2,
		cl.AppConfig.SubscriptionInfo.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.SubscriptionInfo.Command,
		Short: cl.AppConfig.SubscriptionInfo.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		s := &SubscriptionInfo{}
		s, err := s.getSubscriptionInfo()
		if err != nil {
			return err
		}

		clearTerminal()
		s.Print()
		return nil
	}
	return cmd, nil
}

func (s *SubscriptionInfo) getSubscriptionInfo() (*SubscriptionInfo, error) {
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

	url := strings.Replace(cl.AppConfig.SubscriptionInfo.URL, "{{subscriptionID}}", cl.AppConfig.AccessToken.SubscriptionID, 1)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	client := &http.Client {}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, s)
	if err != nil {
		fmt.Println("subscription info unmarshal body response: ", err)
	}

	return s, nil
}

func (s *SubscriptionInfo) Print() {

	fmt.Printf(
		`
Azure Subscription Inforamtion: %s
--------------------------------------
Name:                     %s
Authorization Source:     %s
Manage By Tenants:        %v
Status:                   %s
Policies:                 %v

`,s.SubscriptionID,
s.DisplayName,
s.AuthorizationSource,
s.ManagedByTenants,
s.State,
s.SubscriptionPolicies,
)

}