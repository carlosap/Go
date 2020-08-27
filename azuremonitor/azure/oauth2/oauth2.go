package oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/common/httpclient"
	c "github.com/Go/azuremonitor/config"
	"net/http"
	"strings"
)

var (
	configuration    c.CmdConfig
)

func init(){
	configuration, _ = c.GetCmdConfig()
}

type AccessToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
}


func (at *AccessToken) ExecuteRequest(r httpclient.IRequest) {

	request := httpclient.Request{
		"AccessToken",
		r.GetUrl(),
		r.GetMethod(),
		r.GetPayload(),
		r.GetHeader(),
		true,
	}
	_ = request.Execute()
	body := request.GetResponse()
	err := json.Unmarshal(body, at)
	if err != nil {
		fmt.Println("unmarshal body response: ", err)
	}

}
func (at *AccessToken) GetUrl() string {

	url := strings.Replace(configuration.AccessToken.URL, "{{tenantID}}", configuration.AccessToken.TenantID, 1)
	return url
}
func (at *AccessToken) GetMethod() string {
	return httpclient.Methods.POST
}
func (at *AccessToken) GetPayload() string {
	strPayload := fmt.Sprintf("grant_type=%s&client_id=%s&client_secret=%s&scope=%s",
		configuration.AccessToken.GrantType,
		configuration.AccessToken.ClientID,
		configuration.AccessToken.ClientSecret,
		configuration.AccessToken.Scope)

	return strPayload
}
func (at *AccessToken) GetHeader() http.Header {

	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	return header
}
func (at *AccessToken) Print() {

	fmt.Printf(
		`
Access Token:
--------------------------------------
%s
`, at.AccessToken)

}
