package costmanagement

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/csv"
	"github.com/Go/azuremonitor/common/httpclient"
	"net/http"
	"strings"
)

type LogicAppWorkFlowResponse struct {
	Responses []Responses `json:"responses"`
}

type LogicAppWorkFlow struct {
	Resource azure.Resource `json:"resource"`
	WorkflowExecutionsAvg     float64 `json:"workflow_executions"`
	WorkflowActionExecutionsAvg     float64 `json:"workflow_action_executions"`
	NativeOperationExecutionsAvg         float64 `json:"native_operation_executions"`
	StandardConnectorExecutionsAvg     float64 `json:"standard_connector_executions"`
	StorageConsumptionExecutionsAvg float64 `json:"storage_consumption_executions"`
	Responses []Responses `json:"responses"`
}

type LogicAppWorkFlows []LogicAppWorkFlow

var (
	mapLogicAppWorkFlow = make(map[string]LogicAppWorkFlow)
	LogicApp_Workflows = LogicAppWorkFlows{}
)


func (lg *LogicAppWorkFlow) ExecuteRequest(r httpclient.IRequest) {

	//1-Filters Storage Disk only
	requests := lg.getRequests()
	requests.Execute()

	//2-Serializes All Storage Disks and Sets Metrics
	LogicApp_Workflows = lg.parseRequests(requests)

}

func (lg *LogicAppWorkFlow) GetUrl() string {

	url := azure.QueryUrl
	return url
}
func (lg *LogicAppWorkFlow) GetMethod() string {
	return httpclient.Methods.POST
}
func (lg *LogicAppWorkFlow) GetPayload() string {
	payload := azure.LogicAppUsagePayload
	payload = strings.ReplaceAll(payload, "{{startdate}}", StartDate)
	payload = strings.ReplaceAll(payload, "{{enddate}}", EndDate)
	payload = strings.ReplaceAll(payload, "{{subscriptionid}}", configuration.AccessToken.SubscriptionID)
	payload = strings.ReplaceAll(payload, "{{resourcegroup}}", lg.Resource.ResourceGroup)
	payload = strings.ReplaceAll(payload, "{{resourceid}}",lg.Resource.ResourceID )
	return payload
}
func (lg *LogicAppWorkFlow) GetHeader() http.Header {
	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (lg *LogicAppWorkFlow) Print() {

	if len(LogicApp_Workflows) > 0 {
		fmt.Printf("Usage Report Logic App Workflow:\n")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Location Prefix,Consumption Type,Meter,Cost," +
			"Workflow Executions Avg,Workflow Action Executions Avg,Native Operation Executions Avg,Standard Connector Executions Avg,Storage Consumption Executions Avg")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range LogicApp_Workflows {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,$%f,%f,%f,%f,%f,%f\n",item.Resource.ResourceGroup, item.Resource.ResourceID, item.Resource.Service,
				item.Resource.ServiceType, item.Resource.Location,item.Resource.LocationPrefix, item.Resource.ChargeType, item.Resource.Meter, item.Resource.Cost,
				item.WorkflowExecutionsAvg, item.WorkflowActionExecutionsAvg,item.NativeOperationExecutionsAvg, item.StandardConnectorExecutionsAvg, item.StorageConsumptionExecutionsAvg)
		}
	} else {
		fmt.Printf("-\n\n\n")
	}
}

//---------------Other Functions --------------------------------------------------------------
func (lg *LogicAppWorkFlow) getRequests() httpclient.Requests {
	requests := httpclient.Requests{}
	if len(Resources) > 0 {
		for index, resource := range Resources {
			if resource.Service == "logic apps" && resource.ServiceType == "workflows" && resource.ChargeType == "usage" && resource.Cost > 0.0 {
				rName := "lg_" + resource.ResourceID + "_" + fmt.Sprintf("%d", index)
				lg.Resource = resource
				request := httpclient.Request{
					Name:    rName,
					Header:  lg.GetHeader(),
					Payload: lg.GetPayload(),
					Url:     lg.GetUrl(),
					Method:  lg.GetMethod(),
					IsCache: false,
				}
				mapLogicAppWorkFlow[rName] = *lg
				requests = append(requests, request)
			}
		}
	}
	return requests
}
func (lg *LogicAppWorkFlow) parseRequests(requests httpclient.Requests) LogicAppWorkFlows {
	lgs := LogicAppWorkFlows{}
	var lgResponse BatchResponse
	for _, item := range requests {
		bData := item.GetResponse()
		if len(bData) > 0 {
			err := json.Unmarshal(bData, &lgResponse)
			if err != nil {
				fmt.Printf("error: failed to unmarshal - %v\n\n", err)
			}
			lgRef, hasKey := mapLogicAppWorkFlow[item.Name]
			if hasKey {
				lg.Resource = lgRef.Resource
				lg.Responses = lgResponse.Responses
				lg.setUsageValue()
				lgs = append(lgs, *lg)
			}
		}
	}
	return lgs
}
func (lg *LogicAppWorkFlow) setUsageValue() {

	if len(lg.Responses) > 0 {
		for _, response := range lg.Responses {
			if len(response.Content.Value) > 0 {
				//fmt.Printf("value: %v\n",response.Content.Value)
				for _, valueItem := range response.Content.Value {
					switch valueItem.Name.Value {
					case "TotalBillableExecutions":
						lg.WorkflowExecutionsAvg = valueItem.Timeseries[0].Data[0].Average
					case "BillableActionExecutions":
						lg.WorkflowActionExecutionsAvg = valueItem.Timeseries[0].Data[0].Average
					case "BillingUsageNativeOperation":
						lg.NativeOperationExecutionsAvg = valueItem.Timeseries[0].Data[0].Average
					case "BillingUsageStandardConnector":
						lg.StandardConnectorExecutionsAvg = valueItem.Timeseries[0].Data[0].Average
					case "BillingUsageStorageConsumption":
						lg.StorageConsumptionExecutionsAvg = valueItem.Timeseries[0].Data[0].Average
					}
				}
			}
		}
	}
}
func (lg *LogicAppWorkFlow) WriteCSV(filepath string) {

	if len(LogicApp_Workflows) > 0 {
		var matrix [][]string
		rec := []string{"Resource Group","ResourceID","Service Name","Resource Type","Resource Location","Location Prefix","Consumption Type","Meter","Cost",
			"Workflow Executions Avg","Workflow Action Executions Avg","Native Operation Executions Avg","Standard Connector Executions Avg","Storage Consumption Executions Avg"}
		matrix = append(matrix, rec)
		for _, item := range LogicApp_Workflows {
			var rec []string
			rec = append(rec, item.Resource.ResourceGroup)
			rec = append(rec, item.Resource.ResourceID)
			rec = append(rec, item.Resource.Service)
			rec = append(rec, item.Resource.ServiceType)
			rec = append(rec, item.Resource.Location)
			rec = append(rec, item.Resource.LocationPrefix)
			rec = append(rec, item.Resource.ChargeType)
			rec = append(rec, item.Resource.Meter)
			rec = append(rec, fmt.Sprintf("%f",item.Resource.Cost))

			rec = append(rec, fmt.Sprintf("%f",item.WorkflowExecutionsAvg))
			rec = append(rec, fmt.Sprintf("%f",item.WorkflowActionExecutionsAvg))
			rec = append(rec, fmt.Sprintf("%f",item.NativeOperationExecutionsAvg))
			rec = append(rec, fmt.Sprintf("%f",item.StandardConnectorExecutionsAvg))
			rec = append(rec, fmt.Sprintf("%f",item.StorageConsumptionExecutionsAvg))
			matrix = append(matrix, rec)
		}
		csv.SaveMatrixToFile(filepath, matrix)
	}
}
