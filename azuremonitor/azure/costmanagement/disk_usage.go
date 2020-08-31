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
	"time"
)

type StorageDiskResponse struct {
	Responses []Responses `json:"responses"`
}


type StorageDisk struct {
	Resource azure.Resource `json:"resource"`
	CpuUtilization      float64 `json:"cpu_utilization"`
	DiskReads     float64 `json:"memory_reads"`
	DiskWrites         float64 `json:"disk_writes"`
	NetworkSentRate     float64 `json:"network_sent_rate"`
	NetworkReceivedRate float64 `json:"network_received_rate"`
	Responses []Responses `json:"responses"`
}

type StorageDisks []StorageDisk

var (
	mapStorageDisks = make(map[string]VirtualMachine)
	Storage_Disks = StorageDisk{}
)





func (sd *StorageDisk) ExecuteRequest(r httpclient.IRequest) {

	//1-Three Node of Resources
	rg := ResourceGroupCost{}
	rg.ExecuteRequest(&rg)

	//2-Filters Storage Disk only
	requests := sd.getRequests()
	requests.Execute()

	//3-Serializes All Storage Disks and Sets Metrics
	Virtual_Machines = sd.parseRequests(requests)

	//4-Storage Disks can be used through any output requirements
}

func (sd *StorageDisk) GetUrl() string {

	url := azure.QueryUrl
	return url
}
func (sd *StorageDisk) GetMethod() string {
	return httpclient.Methods.POST
}
func (sd *StorageDisk) GetPayload() string {

	payload := azure.VmUsagePayload
	payload = strings.ReplaceAll(payload, "{{startdate}}", StartDate)
	payload = strings.ReplaceAll(payload, "{{enddate}}", EndDate)
	payload = strings.ReplaceAll(payload, "{{subscriptionid}}", configuration.AccessToken.SubscriptionID)
	payload = strings.ReplaceAll(payload, "{{resourcegroup}}", sd.Resource.ResourceGroup)
	payload = strings.ReplaceAll(payload, "{{resourceid}}", sd.Resource.ResourceID)
	return payload
}
func (sd *StorageDisk) GetHeader() http.Header {
	at := oauth2.AccessToken{}
	at.ExecuteRequest(&at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (sd *StorageDisk) Print() {
	if len(Virtual_Machines) > 0 {
		fmt.Printf("Usage Report:\n")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Location Prefix,Consumption Type,Meter,Cost,Percentage CPU Avg,Bytes read from disk during monitoring period,Bytes written to disk during monitoring period,Incoming Traffic,Outgoing Traffic")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range Virtual_Machines {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,$%s,%f,%f,%f,%f,%f,%f,%f\n",item.Resource.ResourceGroup, item.Resource.ResourceID, item.Resource.Service,
				item.Resource.ServiceType, item.Resource.Location,item.Resource.LocationPrefix, item.Resource.ChargeType, item.Resource.Meter, item.Resource.Cost,
				item.CpuUtilization, item.DiskReads,item.DiskWrites, item.NetworkReceivedRate, item.NetworkSentRate,
				item.NetworkSentRate, item.NetworkReceivedRate)
		}
	} else {
		fmt.Printf("-")
	}
}

//---------------Other Functions --------------------------------------------------------------
func (sd *StorageDisk) getRequests() httpclient.Requests {
	requests := httpclient.Requests{}
	if len(Resources) > 0 {
		for _, resource := range Resources {
			if resource.Service == "virtual machines" && resource.ServiceType == "virtualmachines" &&
				len(resource.Cost) > 0 && resource.ChargeType == "usage" {

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
func (sd *StorageDisk) parseRequests(requests httpclient.Requests) VirtualMachines {
	vms := VirtualMachines{}
	var vmResponse BatchResponse
	for _, item := range requests {
		bData := item.GetResponse()
		if len(bData) > 0 {
			_ = json.Unmarshal(bData, &vmResponse)
			vmRef, hasKey := mapVirtualMachines[item.Name]
			if hasKey {
				sd.Resource = vmRef.Resource
				sd.Responses = vmResponse.Responses
				sd.setUsageValue()
				vms = append(vms, *vm)
			}
		}
	}
	return vms
}
func (sd *StorageDisk) setUsageValue() {

	if len(sd.Responses) > 0 {
		for _, response := range sd.Responses {
			if len(response.Content.Value) > 0 {
				for _, valueItem := range response.Content.Value {
					switch valueItem.Name.Value {
					case "Percentage CPU":
						sd.CpuUtilization = valueItem.Timeseries[0].Data[0].Average
					case "Disk Read Bytes":
						sd.DiskReads = valueItem.Timeseries[0].Data[0].Total
					case "Disk Write Bytes":
						sd.DiskWrites = valueItem.Timeseries[0].Data[0].Total
					case "Network In Total":
						sd.NetworkReceivedRate = valueItem.Timeseries[0].Data[0].Total
					case "Network Out Total":
						sd.NetworkSentRate = valueItem.Timeseries[0].Data[0].Total
					}
				}
			}
		}
	}
}
func (sd *StorageDisk) WriteCSV(filepath string) {

	if len(Virtual_Machines) > 0 {
		var matrix [][]string
		rec := []string{"Resource Group","ResourceID","Service Name","Resource Type","Resource Location","Location Prefix","Consumption Type","Meter","Cost",
			"Percentage CPU Avg","Bytes read from disk during monitoring period","Bytes written to disk during monitoring period","Incoming Traffic","Outgoing Traffic"}
		matrix = append(matrix, rec)
		for _, item := range Virtual_Machines {
			//fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,$%s,%s,%s,%s,%s,%s,%s,%s\n", item.ResourceGroup, item.ResourceID, item.Service, item.ServiceType, item.Location,item.Meter, item.Cost)
			var rec []string
			rec = append(rec, item.Resource.ResourceGroup)
			rec = append(rec, item.Resource.ResourceID)
			rec = append(rec, item.Resource.Service)
			rec = append(rec, item.Resource.ServiceType)
			rec = append(rec, item.Resource.Location)
			rec = append(rec, item.Resource.LocationPrefix)
			rec = append(rec, item.Resource.ChargeType)
			rec = append(rec, item.Resource.Meter)
			rec = append(rec, item.Resource.Cost)

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
