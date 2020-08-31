package costmanagement

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/batch"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/csv"
	"github.com/Go/azuremonitor/common/httpclient"
	c "github.com/Go/azuremonitor/config"
	"net/http"
	"strings"
)

type ResourceGroupCost struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	ResourceGroupName string      `json:"resourcegroupname"`
	Type              string      `json:"type"`
	Location          interface{} `json:"location"`
	Sku               interface{} `json:"sku"`
	ETag              interface{} `json:"eTag"`
	Properties        struct {
		NextLink interface{} `json:"nextLink"`
		Columns  []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"columns"`
		Rows [][]interface{} `json:"rows"`
	} `json:"properties"`
}
type ResourceGroupCosts []ResourceGroupCost

var (
	configuration    c.CmdConfig
	StartDate        string
	EndDate          string
	CsvRguReportName string
	CsvRgcReportName string
	IgnoreZeroCost bool
	SaveCsv bool
	Resources = azure.Resources{}
)



func init(){
	configuration, _ = c.GetCmdConfig()
}
func (rgc *ResourceGroupCost) ExecuteRequest(r httpclient.IRequest) {

	requests := rgc.getRequests()
	requests.Execute()
	rgc.parseRequests(requests)

}
func (rgc *ResourceGroupCost) GetUrl() string {

	url := strings.Replace(configuration.ResourceGroupCost.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	url = strings.Replace(url, "{{resourceGroup}}", rgc.Name, 1)
	return url
}
func (rgc *ResourceGroupCost) GetMethod() string {
	return httpclient.Methods.POST
}
func (rgc *ResourceGroupCost) GetPayload() string {

	if StartDate == "" || EndDate == "" {
		fmt.Println("StartDate and EndDate are Required in the payload. -", rgc.Name)
		return ""
	}

	url := fmt.Sprintf("{\"type\": \"ActualCost\",\"dataSet\": {\"granularity\": \"None\","+
		"\"aggregation\": {\"totalCost\": {\"name\": \"Cost\",\"function\": \"Sum\"},"+
		"\"totalCostUSD\": {\"name\": \"CostUSD\",\"function\": \"Sum\"}},"+
		"\"grouping\": [{\"type\": \"Dimension\",\"name\": \"ResourceId\"},"+
		" {\"type\": \"Dimension\",\"name\": \"ResourceType\"}, {\"type\": \"Dimension\",\"name\": \"ResourceLocation\"}, "+
		"{\"type\": \"Dimension\",\"name\": \"ChargeType\"}, {\"type\": \"Dimension\",\"name\": \"ResourceGroupName\"}, "+
		"{\"type\": \"Dimension\",\"name\": \"PublisherType\"}, {\"type\": \"Dimension\",\"name\": \"ServiceName\"}, "+
		"{\"type\": \"Dimension\",\"name\": \"Meter\"}],\"include\": [\"Tags\"]},\"timeframe\": \"Custom\","+
		"\"timePeriod\": {"+
		"\"from\": \"%sT00:00:00+00:00\","+
		"\"to\": \"%sT23:59:59+00:00\"}}",
		StartDate,
		EndDate,
	)
	return url
}
func (rgc *ResourceGroupCost) GetHeader() http.Header {
	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (rgc *ResourceGroupCost) Print() {
	if len(Resources) > 0 {
		fmt.Println("Consumption Report:")
		fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Location Prefix,Consumption Type,Meter,Cost")
		fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range Resources {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,$%s\n", item.ResourceGroup, item.ResourceID, item.Service, item.ServiceType, item.Location,item.LocationPrefix,item.ChargeType, item.Meter, item.Cost)
		}
	} else {
		fmt.Printf("-")
	}
}
func (rgc *ResourceGroupCost) WriteCSV(filepath string) {

	if len(Resources) > 0 {
		var matrix [][]string
		rec := []string{"Resource Group", "ResourceID", "Service Name", "Resource Type", "Resource Location","Location Prefix", "Consumption Type", "Meter", "Cost"}
		matrix = append(matrix, rec)
		for _, item := range Resources {
			//fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,$%s\n", item.ResourceGroup, item.ResourceID, item.Service, item.ServiceType,item.Location,item.LocationPrefix, item.ChargeType,item.Meter, item.Cost)
			var rec []string
			rec = append(rec, item.ResourceGroup)
			rec = append(rec, item.ResourceID)
			rec = append(rec, item.Service)
			rec = append(rec, item.ServiceType)
			rec = append(rec, item.Location)
			rec = append(rec, item.LocationPrefix)
			rec = append(rec, item.ChargeType)
			rec = append(rec, item.Meter)
			rec = append(rec, item.Cost)
			matrix = append(matrix, rec)
		}
		csv.SaveMatrixToFile(filepath, matrix)
	}
}
func (rgc *ResourceGroupCost) getRequests() httpclient.Requests {
	requests := httpclient.Requests{}
	rgl := batch.ResourceGroupList{}
	rgl.ExecuteRequest(&rgl)

	for _, item := range rgl.ToList() {
		rgc.Name = item
		rgc.ResourceGroupName = item
		request := httpclient.Request{}
		request.Name = item
		request.Header = rgc.GetHeader()
		request.Payload = rgc.GetPayload()
		request.Url = rgc.GetUrl()
		request.Method = rgc.GetMethod()
		request.IsCache = true
		requests = append(requests, request)
	}
	return requests
}
func (rgc *ResourceGroupCost) parseRequests(requests httpclient.Requests) {
	for _, item := range requests {
		bData := item.GetResponse()
		if len(bData) > 0 {
			_ = json.Unmarshal(bData, rgc)
			rgc.ResourceGroupName = item.Name
			rgc.addResource()
		}
	}
}
func (rgc *ResourceGroupCost) addResource() {

	for i := 0; i < len(rgc.Properties.Rows); i++ {
		row := rgc.Properties.Rows[i]
		if len(row) > 0 {
			costUSD := fmt.Sprintf("%v", row[1])
			resourceId := fmt.Sprintf("%v", row[2])
			resourceType := fmt.Sprintf("%v", row[3])
			resourceLocation := fmt.Sprintf("%v", row[4])
			chargeType := fmt.Sprintf("%v", row[5])
			serviceName := fmt.Sprintf("%v", row[8])
			meter := fmt.Sprintf("%v", row[9])

			//format cost
			if len(costUSD) > 5 {
				costUSD = costUSD[0:5]
			}

			if IgnoreZeroCost {
				if costUSD == "0" {
					continue
				}
			}

			//remove path
			if strings.Contains(resourceType, "/") {
				pArray := strings.Split(resourceType, "/")
				resourceType = pArray[len(pArray)-1]
			}

			if strings.Contains(resourceId, "/") {
				pArray := strings.Split(resourceId, "/")
				resourceId = pArray[len(pArray)-1]
			}


			resource := azure.Resource{
				ResourceGroup: rgc.ResourceGroupName,
				ResourceID: resourceId,
				Service: serviceName,
				ServiceType: resourceType,
				Location: resourceLocation,
				LocationPrefix: resourceLocation,
				ChargeType: chargeType,
				Meter: meter,
				Cost: costUSD,
			}
			Resources = append(Resources, resource)
		}
	}
}


