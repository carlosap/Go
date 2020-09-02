package costmanagement

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/batch"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/convert"
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

	payload := azure.ActualCostManagementPayload
	payload = strings.ReplaceAll(payload, "{{startdate}}", StartDate)
	payload = strings.ReplaceAll(payload, "{{enddate}}", EndDate)

	return payload
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
		fmt.Println("Resource Group," +
			"ResourceID," +
			"Resource Type," +
			"Resource Location," +
			"Charge Type," +
			"Service Name," +
			"Meter," +
			"Meter Category," +
			"Meter SubCategory," +
			"Service Family," +
			"Unit Of Measure," +
			"Cost Allocation Rule Name," +
			"Product," +
			"Frequency," +
			"Pricing Model," +
			"Currency," +
			"UsageQuantity," +
			"PreTaxCostUSD")
		fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")

		for _, item := range Resources {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%f,$%f\n",
				item.ResourceGroupName,
				item.ResourceID,
				item.ResourceType,
				item.ResourceLocation,
				item.ChargeType,
				item.ServiceName,
				item.Meter,
				item.MeterCategory,
				item.MeterSubCategory,
				item.ServiceFamily,
				item.UnitOfMeasure,
				item.CostAllocationRuleName,
				item.Product,
				item.Frequency,
				item.PricingModel,
				item.Currency,
				item.UsageQuantity,
				item.PreTaxCostUSD)
		}
	} else {
		fmt.Printf("-")
	}
}
func (rgc *ResourceGroupCost) WriteCSV(filepath string) {

	if len(Resources) > 0 {
		var matrix [][]string
		rec := []string{
			"Resource Group",
			"ResourceID",
			"Resource Type",
			"Resource Location",
			"Charge Type",
			"Service Name",
			"Meter",
			"Meter Category",
			"Meter SubCategory",
			"Service Family",
			"Unit Of Measure",
			"Cost Allocation Rule Name",
			"Product",
			"Frequency",
			"Pricing Model",
			"Currency",
			"UsageQuantity",
			"PreTaxCostUSD"}

		matrix = append(matrix, rec)
		for _, item := range Resources {
			var rec []string
			rec = append(rec, item.ResourceGroupName)
			rec = append(rec, item.ResourceID)
			rec = append(rec, item.ResourceType)
			rec = append(rec, item.ResourceLocation)
			rec = append(rec, item.ChargeType)
			rec = append(rec, item.ServiceName)
			rec = append(rec, item.Meter)
			rec = append(rec, item.MeterCategory)
			rec = append(rec, item.MeterSubCategory)
			rec = append(rec, item.ServiceFamily)
			rec = append(rec, item.UnitOfMeasure)
			rec = append(rec, item.CostAllocationRuleName)
			rec = append(rec, item.Product)
			rec = append(rec, item.Frequency)
			rec = append(rec, item.PricingModel)
			rec = append(rec, item.Currency)
			rec = append(rec, fmt.Sprintf("%f",item.UsageQuantity))
			rec = append(rec, fmt.Sprintf("%f", item.PreTaxCostUSD))
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
		//fmt.Printf("Properties: %v\n", rgc.Properties)
		if len(row) > 0 {
			preTaxCost := fmt.Sprintf("%v", row[0])
			usageQuantity := fmt.Sprintf("%v", row[1])
			resourceId := fmt.Sprintf("%v", row[2])
			resourceType := fmt.Sprintf("%v", row[3])
			resourceLocation := fmt.Sprintf("%v", row[4])
			chargeType := fmt.Sprintf("%v", row[5])
			resourceGroupName := fmt.Sprintf("%v", row[6])
			serviceName := fmt.Sprintf("%v", row[7])
			meter := fmt.Sprintf("%v", row[8])
			meterCategory := fmt.Sprintf("%v", row[9])
			meterSubCategory := fmt.Sprintf("%v", row[10])
			serviceFamily := fmt.Sprintf("%v", row[11])
			unitOfMeasure := fmt.Sprintf("%v", row[12])
			costAllocationRuleName := fmt.Sprintf("%v", row[13])
			product := fmt.Sprintf("%v", row[14])
			frequency := fmt.Sprintf("%v", row[15])
			pricingModel := fmt.Sprintf("%v", row[16])
			//tags := fmt.Sprintf("%v", row[17])
			currency := fmt.Sprintf("%v", row[18])


			pCost, _ := convert.StringToFloat(preTaxCost)
			if IgnoreZeroCost {
				if pCost <= 0.0 {
					continue
				}
			}

			uQuantity, _ :=convert.StringToFloat(usageQuantity)
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
				ResourceGroupName: resourceGroupName,
				ResourceID: resourceId,
				ResourceType: resourceType,
				ResourceLocation: resourceLocation,
				ChargeType: chargeType,
				ServiceName: serviceName,
				Meter: meter,
				MeterCategory: meterCategory,
				MeterSubCategory: meterSubCategory,
				ServiceFamily: serviceFamily,
				UnitOfMeasure: unitOfMeasure,
				CostAllocationRuleName: costAllocationRuleName,
				Product: product,
				Frequency: frequency,
				PricingModel: pricingModel,
				//Tags: tags
				Currency: currency,
				PreTaxCostUSD: pCost,
				UsageQuantity: uQuantity,
			}

			Resources = append(Resources, resource)
		}
	}
}


