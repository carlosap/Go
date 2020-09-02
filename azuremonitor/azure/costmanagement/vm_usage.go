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

type VirtualMachine struct {
	Resource azure.Resource `json:"resource"`
	CpuUtilization      float64 `json:"cpu_utilization"`
	DiskReads     float64 `json:"memory_reads"`
	DiskWrites         float64 `json:"disk_writes"`
	NetworkSentRate     float64 `json:"network_sent_rate"`
	NetworkReceivedRate float64 `json:"network_received_rate"`
	OsDisk  OsDisk `json:"osdisk"`
	Responses []Responses `json:"responses"`
}

type VirtualMachines []VirtualMachine

var (

	mapVirtualMachines = make(map[string]VirtualMachine)
	Virtual_Machines = VirtualMachines{}
)


func (vm *VirtualMachine) ExecuteRequest(r httpclient.IRequest) {

	//1-Filters Virtual Machines only
	requests := vm.getRequests()
	requests.Execute()

	//2-Serializes All VMs and Sets Metrics
	Virtual_Machines = vm.parseRequests(requests)
}

func (vm *VirtualMachine) GetUrl() string {

	url := azure.QueryUrl
	return url
}
func (vm *VirtualMachine) GetMethod() string {
	return httpclient.Methods.POST
}
func (vm *VirtualMachine) GetPayload() string {

	payload := azure.VmUsagePayload
	payload = strings.ReplaceAll(payload, "{{startdate}}", StartDate)
	payload = strings.ReplaceAll(payload, "{{enddate}}", EndDate)
	payload = strings.ReplaceAll(payload, "{{subscriptionid}}", configuration.AccessToken.SubscriptionID)
	payload = strings.ReplaceAll(payload, "{{resourcegroup}}", vm.Resource.ResourceGroupName)
	payload = strings.ReplaceAll(payload, "{{resourceid}}", vm.Resource.ResourceID)
	return payload
}

func (vm *VirtualMachine) GetHeader() http.Header {
	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (vm *VirtualMachine) Print() {
	if len(Virtual_Machines) > 0 {
		fmt.Printf("VM Usage Report:\n")
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
			"OS Type," +
			"Disk Name," +
			"Storage Account Type" +
			"Percentage CPU Avg," +
			"Bytes Read," +
			"Bytes Written," +
			"Incoming Traffic (Network Received)," +
			"Outgoing Traffic (Network Sent)")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range Virtual_Machines {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%f,$%f,%s,%s,%s,%f,%f,%f,%f,%f\n",
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

				item.OsDisk.OsType,
				item.OsDisk.Name,
				item.OsDisk.ManagedDisk.StorageAccountType,
				item.CpuUtilization,
				item.DiskReads,
				item.DiskWrites,
				item.NetworkReceivedRate,
				item.NetworkSentRate)
		}
	} else {
		fmt.Printf("-")
	}
}

//---------------Other Functions --------------------------------------------------------------
func (vm *VirtualMachine) getRequests() httpclient.Requests {
	requests := httpclient.Requests{}
	if len(Resources) > 0 {
		for _, resource := range Resources {
			if resource.Service == "virtual machines" && resource.ServiceType == "virtualmachines" &&
				resource.Cost > 0.0 && resource.ChargeType == "usage" {

				rName := "vm_" + resource.ResourceID
				vm.Resource = resource
				request := httpclient.Request{
					Name:    rName,
					Header:  vm.GetHeader(),
					Payload: vm.GetPayload(),
					Url:     vm.GetUrl(),
					Method:  vm.GetMethod(),
					IsCache: false,
				}
				mapVirtualMachines[rName] = *vm
				requests = append(requests, request)
			}
		}
	}
	return requests
}
func (vm *VirtualMachine) parseRequests(requests httpclient.Requests) VirtualMachines {
	vms := VirtualMachines{}
	var vmResponse BatchResponse
	for _, item := range requests {
		bData := item.GetResponse()
		if len(bData) > 0 {
			_ = json.Unmarshal(bData, &vmResponse)
			vmRef, hasKey := mapVirtualMachines[item.Name]
			if hasKey {
				vm.Resource = vmRef.Resource
				vm.Responses = vmResponse.Responses
				vm.setUsageValue()
				vms = append(vms, *vm)
			}
		}
	}
	return vms
}
func (vm *VirtualMachine) setUsageValue() {

	if len(vm.Responses) > 0 {
		vm.OsDisk = vm.Responses[0].Content.Properties.StorageProfile.OsDisk
		for _, response := range vm.Responses {

			if len(response.Content.Value) > 0 {
				for _, valueItem := range response.Content.Value {
					switch valueItem.Name.Value {
					case "Percentage CPU":
						vm.CpuUtilization = valueItem.Timeseries[0].Data[0].Average
					case "Disk Read Bytes":
						vm.DiskReads = valueItem.Timeseries[0].Data[0].Total
					case "Disk Write Bytes":
						vm.DiskWrites = valueItem.Timeseries[0].Data[0].Total
					case "Network In Total":
						vm.NetworkReceivedRate = valueItem.Timeseries[0].Data[0].Total
					case "Network Out Total":
						vm.NetworkSentRate = valueItem.Timeseries[0].Data[0].Total
					}
				}
			}
		}
	}
}

func (vm *VirtualMachine) WriteCSV(filepath string) {

	if len(Virtual_Machines) > 0 {
		var matrix [][]string
		rec := []string{"Resource Group","ResourceID","Service Name","Resource Type","Resource Location","Location Prefix","Consumption Type","Meter","Cost",
			"Percentage CPU Avg","Bytes read from disk during monitoring period","Bytes written to disk during monitoring period","Incoming Traffic","Outgoing Traffic"}
		matrix = append(matrix, rec)
		for _, item := range Virtual_Machines {
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

			rec = append(rec, fmt.Sprintf("%f",item.CpuUtilization))
			rec = append(rec, fmt.Sprintf("%f",item.DiskReads))
			rec = append(rec, fmt.Sprintf("%f",item.DiskWrites))
			rec = append(rec, fmt.Sprintf("%f",item.NetworkReceivedRate))
			rec = append(rec, fmt.Sprintf("%f",item.NetworkSentRate))
			matrix = append(matrix, rec)
		}
		csv.SaveMatrixToFile(filepath, matrix)
	}
}
