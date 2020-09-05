package costmanagement

import (
	"encoding/json"
	"fmt"
	str "github.com/Go/azuremonitor/common/strings"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/azure/subscription"
	"github.com/Go/azuremonitor/common/csv"
	"github.com/Go/azuremonitor/common/httpclient"
	"net/http"
	"strings"
)

type StorageAccountResponse struct {
	Responses []Responses `json:"responses"`
}

type StorageAccount struct {
	Resource azure.Resource `json:"resource"`
	EgressAvg     float64 `json:"egress_avg"`
	IngressAvg     float64 `json:"ingress_avg"`
	TransactionTotal         float64 `json:"transaction_total"`

	BlobEgressAvg float64 `json:"blob_egress_avg"`
	BlobIngressAvg float64 `json:"blob_ingress_avg"`
	BlobTransactionTotal float64 `json:"blob_transaction_total"`
	BlobCountAvg float64 `json:"blob_count_avg"`

	FileCountAvg float64 `json:"file_count_avg"`
	FileEgressAvg float64 `json:"file_egress_avg"`
	FileIngressAvg float64 `json:"file_ingress_avg"`
	FileTransactionTotal float64 `json:"file_transaction_total"`

	QueueCountAvg float64 `json:"queue_count_avg"`
	QueueEgressAvg float64 `json:"queue_egress_avg"`
	QueueIngress float64 `json:"queue_ingress_avg"`
	QueueTransactionTotal float64 `json:"queue_transaction_total"`

	TableCountAvg float64 `json:"table_count_avg"`
	TableEntityCountAvg float64 `json:"table_entity_count_avg"`
	TableEgressAvg float64 `json:"table_egress_avg"`
	TableIngressAvg float64 `json:"table_ingress_avg"`
	TableTransactionsTotal float64 `json:"table_transactions_total"`

	Responses []Responses `json:"responses"`
}

type StorageAccounts []StorageAccount

var (
	mapStorageAccount = make(map[string]StorageAccount)
	Storage_StorageAccounts = StorageAccounts{}
)


func (st *StorageAccount) ExecuteRequest(r httpclient.IRequest) {

	//1-Filters Storage Disk only
	requests := st.getRequests()
	requests.Execute()

	//2-Serializes All Storage Disks and Sets Metrics
	Storage_StorageAccounts = st.parseRequests(requests)

}

func (st *StorageAccount) GetUrl() string {

	url := azure.QueryUrl
	return url
}
func (st *StorageAccount) GetMethod() string {
	return httpclient.Methods.POST
}
func (st *StorageAccount) GetPayload() string {


	payload := azure.StorageStorageAccountPayload
	payload = strings.ReplaceAll(payload, "{{startdate}}", StartDate)
	payload = strings.ReplaceAll(payload, "{{enddate}}", EndDate)
	payload = strings.ReplaceAll(payload, "{{subscriptionid}}", configuration.AccessToken.SubscriptionID)
	payload = strings.ReplaceAll(payload, "{{resourcegroup}}", st.Resource.ResourceGroupName)
	payload = strings.ReplaceAll(payload, "{{resourceid}}", st.Resource.ResourceID )
	return payload
}
func (st *StorageAccount) GetHeader() http.Header {
	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (st *StorageAccount) Print() {

	if len(Storage_StorageAccounts) > 0 {
		fmt.Printf("Usage Report Storage Disk:\n")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
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
			"PreTaxCostUSD," +
			"StorageAccountEgressAvg," +
			"StorageAccountIngressAvg," +
			"StorageAccountTransactionTotal," +
			"BlobEgressAvg," +
			"BlobIngressAvg," +
			"BlobTransactionTotal," +
			"BlobCountAvg," +
			"FileCountAvg," +
			"FileEgressAvg," +
			"FileIngressAvg," +
			"FileTransactionTotal," +
			"QueueCountAvg," +
			"QueueEgressAvg," +
			"QueueIngress," +
			"QueueTransactionTotal," +
			"TableCountAvg," +
			"TableEntityCountAvg," +
			"TableEgressAvg," +
			"TableIngressAvg," +
			"TableTransactionsTotal")

		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range Storage_StorageAccounts {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%f,$%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f\n",
				item.Resource.ResourceGroupName,
				item.Resource.ResourceID,
				item.Resource.ResourceType,
				item.Resource.ResourceLocation,
				item.Resource.ChargeType,
				item.Resource.ServiceName,
				item.Resource.Meter,
				item.Resource.MeterCategory,
				item.Resource.MeterSubCategory,
				item.Resource.ServiceFamily,
				item.Resource.UnitOfMeasure,
				item.Resource.CostAllocationRuleName,
				item.Resource.Product,
				item.Resource.Frequency,
				item.Resource.PricingModel,
				item.Resource.Currency,
				item.Resource.UsageQuantity,
				item.Resource.PreTaxCostUSD,

				item.EgressAvg,
				item.IngressAvg,
				item.TransactionTotal,

				item.BlobEgressAvg,
				item.BlobIngressAvg,
				item.BlobTransactionTotal,
				item.BlobCountAvg,

				item.FileCountAvg,
				item.FileEgressAvg,
				item.FileIngressAvg,
				item.FileTransactionTotal,

				item.QueueCountAvg,
				item.QueueEgressAvg,
				item.QueueIngress,
				item.QueueTransactionTotal,

				item.TableCountAvg,
				item.TableEntityCountAvg,
				item.TableEgressAvg,
				item.TableIngressAvg,
				item.TableTransactionsTotal)
		}
	} else {
		fmt.Printf("-\n\n\n")
	}
}

//---------------Other Functions --------------------------------------------------------------
func (st *StorageAccount) getRequests() httpclient.Requests {
	requests := httpclient.Requests{}
	if len(Resources) > 0 {
		for index, resource := range Resources {
			if resource.ServiceName == "storage" && resource.ResourceType == "storageaccounts" && resource.ChargeType == "usage" && resource.PreTaxCostUSD > 0.0 {
				rName := "st_" + resource.ResourceID + "_" + fmt.Sprintf("%d", index)
				st.Resource = resource
				request := httpclient.Request{
					Name:    rName,
					Header:  st.GetHeader(),
					Payload: st.GetPayload(),
					Url:     st.GetUrl(),
					Method:  st.GetMethod(),
					IsCache: false,
				}
				mapStorageAccount[rName] = *st
				requests = append(requests, request)
			}
		}
	}
	return requests
}
func (st *StorageAccount) parseRequests(requests httpclient.Requests) StorageAccounts {
	staccounts := StorageAccounts{}
	var sdResponse BatchResponse
	for _, item := range requests {
		bData := item.GetResponse()
		if len(bData) > 0 {
			err := json.Unmarshal(bData, &sdResponse)
			if err != nil {
				fmt.Printf("error: failed to unmarshal - %v\n\n", err)
			}
			//fmt.Printf("data: %s\n\n", string(bData))
			stRef, hasKey := mapStorageAccount[item.Name]
			if hasKey {
				st.Resource = stRef.Resource
				st.Responses = sdResponse.Responses
				st.setUsageValue()
				staccounts = append(staccounts, *st)
			}
		}
	}
	return staccounts
}
func (st *StorageAccount) setUsageValue() {

	if len(st.Responses) > 0 {
		for _, response := range st.Responses {
			if len(response.Content.Value) > 0 {
				//fmt.Printf("value: %v\n",response.Content.Value)
				namespace := str.GetLastValueFromSeparator(response.Content.Namespace, "/")
				valueName := response.Content.Value[0].Name.Value
				switch {
				case namespace == "storageaccounts" && valueName == "Egress":
					st.EgressAvg = response.Content.Value[0].Timeseries[0].Data[0].Average
				case namespace == "storageaccounts" && valueName == "Ingress":
					st.IngressAvg = response.Content.Value[0].Timeseries[0].Data[0].Average
				case namespace == "storageaccounts" && valueName == "Transactions":
					st.TransactionTotal = response.Content.Value[0].Timeseries[0].Data[0].Total
				case namespace == "blobservices" && valueName == "Egress":
					st.BlobEgressAvg = response.Content.Value[0].Timeseries[0].Data[0].Average
				case namespace == "blobservices" && valueName == "Ingress":
					st.BlobIngressAvg = response.Content.Value[0].Timeseries[0].Data[0].Average
				case namespace == "blobservices" && valueName == "Transactions":
					st.BlobTransactionTotal = response.Content.Value[0].Timeseries[0].Data[0].Total
				case namespace == "blobservices" && valueName == "BlobCount":
					st.BlobCountAvg = response.Content.Value[0].Timeseries[0].Data[0].Average
				case namespace == "fileservices" && valueName == "FileCount":
					st.FileCountAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "fileservices" && valueName == "Egress":
					st.FileEgressAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "fileservices" && valueName == "Ingress":
					st.FileIngressAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "fileservices" && valueName == "Transactions":
					st.FileTransactionTotal = getTotalValue(response.Content.Value[0].Timeseries)
				case namespace == "queueservices" && valueName == "QueueCount":
					st.QueueCountAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "queueservices" && valueName == "Egress":
					st.QueueEgressAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "queueservices" && valueName == "Ingress":
					st.QueueIngress = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "queueservices" && valueName == "Transactions":
					st.QueueTransactionTotal = getTotalValue(response.Content.Value[0].Timeseries)
				case namespace == "tableservices" && valueName == "TableCount":
					st.TableCountAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "tableservices" && valueName == "TableEntityCount":
					st.TableEntityCountAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "tableservices" && valueName == "Egress":
					st.TableEgressAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "tableservices" && valueName == "Ingress":
					st.TableIngressAvg = getAvgValue(response.Content.Value[0].Timeseries)
				case namespace == "tableservices" && valueName == "Transactions":
					st.TableTransactionsTotal = getTotalValue(response.Content.Value[0].Timeseries)

				}

			}
		}
	}
}

func getAvgValue(timeseries []Timeseries) float64 {
	var retVal float64
	if len(timeseries) > 0 &&	len(timeseries[0].Data) > 0 {
		retVal = timeseries[0].Data[0].Average
	}
	return retVal
}

func getTotalValue(timeseries []Timeseries) float64 {
	var retVal float64
	if len(timeseries) > 0 &&	len(timeseries[0].Data) > 0 {
		retVal = timeseries[0].Data[0].Total
	}
	return retVal
}

func (st *StorageAccount) WriteCSV(filepath string) {

	if len(Storage_Disks) > 0 {
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
			"PreTaxCostUSD",
			"Disk Read Bytes/sec Avg",
			"Disk Write Bytes/sec Avg",
			"Disk Read Operations/Sec Avg",
			"Disk Write Operations/Sec Avg",
			"Disk Queue Depth"}
		matrix = append(matrix, rec)
		for _, item := range Storage_Disks {
			var rec []string
			rec = append(rec, item.Resource.ResourceGroupName)
			rec = append(rec, item.Resource.ResourceID)
			rec = append(rec, item.Resource.ResourceType)
			rec = append(rec, item.Resource.ResourceLocation)
			rec = append(rec, item.Resource.ChargeType)
			rec = append(rec, item.Resource.ServiceName)
			rec = append(rec, item.Resource.Meter)
			rec = append(rec, item.Resource.MeterCategory)
			rec = append(rec, item.Resource.MeterSubCategory)
			rec = append(rec, item.Resource.ServiceFamily)
			rec = append(rec, item.Resource.UnitOfMeasure)
			rec = append(rec, item.Resource.CostAllocationRuleName)
			rec = append(rec, item.Resource.Product)
			rec = append(rec, item.Resource.Frequency)
			rec = append(rec, item.Resource.PricingModel)
			rec = append(rec, item.Resource.Currency)
			rec = append(rec, fmt.Sprintf("%f", item.Resource.UsageQuantity))
			rec = append(rec, fmt.Sprintf("%f", item.Resource.PreTaxCostUSD))
			rec = append(rec, fmt.Sprintf("%f",item.DiskReads))
			rec = append(rec, fmt.Sprintf("%f",item.DiskWrite))
			rec = append(rec, fmt.Sprintf("%f",item.DiskReadOperations))
			rec = append(rec, fmt.Sprintf("%f",item.DiskWriteOperations))
			rec = append(rec, fmt.Sprintf("%f",item.QueueDepth))
			matrix = append(matrix, rec)
		}
		csv.SaveMatrixToFile(filepath, matrix)
	}
}
