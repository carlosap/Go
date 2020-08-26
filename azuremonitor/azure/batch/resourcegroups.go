package batch

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/httpclient"
	c "github.com/Go/azuremonitor/config"
	"net/http"
)

type ResourceGroups struct {
	Responses []struct {
		Content struct {
			TotalRecords int `json:"totalRecords"`
			Count        int `json:"count"`
			Data         struct {
				Rows [][]interface{} `json:"rows"`
			} `json:"data"`
		} `json:"content"`
	} `json:"responses"`
}

type ResourceGroupList []string
var resourceGroupList ResourceGroupList

var (
	configuration    c.CmdConfig
)

func init(){
	configuration, _ = c.GetCmdConfig()
}


func (rgl *ResourceGroupList) ExecuteRequest(r httpclient.IRequest) {

	rg := ResourceGroups{}
	request := httpclient.Request{
		Name:    "resourcegroups",
		Url:     r.GetUrl(),
		Method:  r.GetMethod(),
		Payload: r.GetPayload(),
		Header:  r.GetHeader(),
		IsCache: true,
	}

	_ = request.Execute()
	body := request.GetResponse()
	err := json.Unmarshal(body, &rg)
	if err != nil {
		fmt.Println("unmarshal body response: ", err)
	}

	setResourceGroup(rg)
}


func (rgl *ResourceGroupList) GetUrl() string {
	return configuration.ResourceGroups.URL
}
func (rgl *ResourceGroupList) GetMethod() string {
	return httpclient.Methods.POST
}
func (rgl *ResourceGroupList) GetPayload() string {
	//SpartanAppSolutions

	payload := fmt.Sprintf("{\"requests\": [{\"content\": {\"subscriptions\": [\"%s\"],"+
		"\"query\": \"(resourcecontainers|where type in~ ('microsoft.resources/subscriptions/resourcegroups'))"+
		"|where type =~ 'microsoft.resources/subscriptions/resourcegroups'| "+
		"extend status = case((properties.provisioningState =~ 'accepted'), "+
		"'Creating',(properties.provisioningState =~ 'deleted'), "+
		"'Deleted',(properties.provisioningState =~ 'deleting'), "+
		"'Deleting',(properties.provisioningState =~ 'failed'), "+
		"'Failed',(properties.provisioningState =~ 'movingresources'), "+
		"'Moving Resources',properties.provisioningState)| "+
		"project id, name, type, location, subscriptionId, resourceGroup, kind, tags, status\\r\\n|"+
		"extend subscriptionDisplayName=case(subscriptionId == '%s','SpartanAppSolutions',subscriptionId)|"+
		"extend locationDisplayName=case(" +
		azure.LocationNames +
		"location)" +
		"|where (type !~ ('microsoft.confluent/organizations'))|"+
		"where (type !~ ('microsoft.securitydetonation/chambers'))|where (type !~ ('microsoft.intelligentitdigitaltwin/digitaltwins'))|"+
		"where (type !~ ('microsoft.connectedcache/cachenodes'))|where (type !~ ('microsoft.serviceshub/connectors'))|"+
		"where not((type =~ ('microsoft.sql/servers/databases')) and ((kind in~ ('system','v2.0,system','v12.0,system','v12.0,user,datawarehouse,gen2,analytics'))))|"+
		"where not((type =~ ('microsoft.sql/servers')) and ((kind =~ ('v12.0,analytics'))))|"+
		"project name,subscriptionDisplayName,locationDisplayName,id,type,kind,location,subscriptionId,resourceGroup,tags|"+
		"sort by tolower(tostring(name)) asc\",\"options\": {\"$top\": 100,\"$skip\": 0,\"$skipToken\": \"\"}},"+
		"\"httpMethod\": \"POST\",\"name\": \"34cc625b-b20a-423a-9563-33faf337b033\","+
		"\"requestHeaderDetails\": {\"commandName\": \"HubsExtension.BrowseResourceGroups.microsoft.resources/subscriptions/resourcegroups.InitialLoad\"},"+
		"\"url\": \"https://management.azure.com/providers/Microsoft.ResourceGraph/resources?api-version=2018-09-01-preview\"}]}",
		configuration.AccessToken.SubscriptionID,
		configuration.AccessToken.SubscriptionID,
	)
	return payload
}
func (rgl *ResourceGroupList) GetHeader() http.Header {

	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)

	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (rgl *ResourceGroupList) Print() {

	fmt.Println("Resource Groups:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i := 0; i < len(resourceGroupList); i++ {
		fmt.Println(resourceGroupList[i])
	}
}

func (rgl *ResourceGroupList) ToList() []string {
	return resourceGroupList
}

func setResourceGroup(rg ResourceGroups) {

	for i := 0; i < len(rg.Responses); i++ {
		rp := rg.Responses[i]
		if rp.Content.TotalRecords > 0 {
			for x := 0; x < rp.Content.TotalRecords; x++ {
				row := rp.Content.Data.Rows[x]
				if len(row) > 0 {
					str := fmt.Sprintf("%v", row[0])
					resourceGroupList = append(resourceGroupList, str)
				}
			}
		}
	}
}
