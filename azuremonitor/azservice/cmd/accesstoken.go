package cmd

import (
	"encoding/json"
	"fmt"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

type AccessToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}

func init() {
	at, err := setAccessTokenCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(at)
}

func setAccessTokenCommand() (*cobra.Command, error) {

	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.AccessToken.DescriptionLine1,
		configuration.AccessToken.DescriptionLine2,
		configuration.AccessToken.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.AccessToken.Command,
		Short: configuration.AccessToken.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		at := &AccessToken{}
		at, err := at.getAccessToken()
		if err != nil {
			return err
		}

		clearTerminal()
		at.Print()
		return nil
	}
	return cmd, nil
}

func (at *AccessToken) getAccessToken() (*AccessToken, error) {

	url := strings.Replace(configuration.AccessToken.URL, "{{tenantID}}", configuration.AccessToken.TenantID, 1)
	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	strPayload := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&scope=%s",
		configuration.AccessToken.GrantType,
		configuration.AccessToken.ClientID,
		configuration.AccessToken.ClientSecret,
		configuration.AccessToken.Scope)

	request := Request{
		"AccessToken",
		url,
		Methods.POST,
		strPayload,
		header,
	}
	_ = request.Execute()
	body := request.GetResponse()
	err := json.Unmarshal(body, at)
	if err != nil {
		fmt.Println("unmarshal body response: ", err)
	}

	return at, err
}

func (at *AccessToken) Print() {

	fmt.Printf(
		`
Access Token:
--------------------------------------
%s
`, at.AccessToken)

}
