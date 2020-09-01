package subscription

import (
"encoding/json"
"fmt"
	"github.com/Go/azuremonitor/azure/oauth2"
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

type ResourceSubscription struct {
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

func (resource *ResourceSubscription) ExecuteRequest(r httpclient.IRequest) {

	request := httpclient.Request{
		"resourcesubscription",
		r.GetUrl(),
		r.GetMethod(),
		r.GetPayload(),
		r.GetHeader(),
		true,
	}
	_ = request.Execute()
	body := request.GetResponse()
	err := json.Unmarshal(body, resource)
	if err != nil {
		fmt.Println("unmarshal body response: ", err)
	}
}
func (resource *ResourceSubscription) GetUrl() string {
	url := strings.Replace(configuration.Resources.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	return url
}
func (resource *ResourceSubscription) GetMethod() string {
	return httpclient.Methods.GET
}
func (resource *ResourceSubscription) GetPayload() string {
	return ""
}
func (resource *ResourceSubscription) GetHeader() http.Header {

	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (resource *ResourceSubscription) Print() {

	fmt.Println("ResourceSubscription Report:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("ResourceId,Type,Kind,Location,Sku Name, Sku Tier,Tags,Plan Name, Plan Promotion Code, Plan Product, Plan Publisher, Manage By ResourceID")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i := 0; i < len(resource.Values); i++ {
		var resourceType, resourceManageby string
		item := resource.Values[i]

		//remove path
		if strings.Contains(item.Type, "/") {
			pArray := strings.Split(item.Type, "/")
			resourceType = pArray[len(pArray)-1]
		}

		if strings.Contains(item.ManagedBy, "/") {
			pArray := strings.Split(item.ManagedBy, "/")
			resourceManageby = pArray[len(pArray)-1]
		}

		fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", item.Name, resourceType, item.Kind, item.Location,
			item.Sku.Name, item.Sku.Tier, item.Tags.MsResourceUsage, item.Plan.Name,
			item.Plan.PromotionCode, item.Plan.Product, item.Plan.Publisher, resourceManageby)
	}
}
func (resource *ResourceSubscription) GetManageByResourceId(resourceid string) string {
	var retVal string
	for i := 0; i < len(resource.Values); i++ {
		item := resource.Values[i]

		if resourceid == item.Name {
			if strings.Contains(item.ManagedBy, "/") {
				pArray := strings.Split(item.ManagedBy, "/")
				retVal = pArray[len(pArray)-1]
			} else {
				retVal = item.ManagedBy
			}
			break
		}
	}

	return retVal
}
