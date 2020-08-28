package costmanagement

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/convert"
	"github.com/Go/azuremonitor/common/csv"
	"github.com/Go/azuremonitor/common/httpclient"
	"net/http"
	"strings"
)

type VirtualMachine struct {
	Resource azure.Resource `json:"resource"`
	CpuUtilization      string `json:"cpu_utilization"`
	MemoryAvailable     string `json:"memory_available"`
	DiskLatency         string `json:"disk_latency"`
	DiskIOPs            string `json:"disk_iops"`
	DiskBytes           string `json:"disk_bytes"`
	NetworkSentRate     string `json:"network_sent_rate"`
	NetworkReceivedRate string `json:"network_received_rate"`
	Tables              []struct {
		Name    string `json:"name"`
		Columns []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"columns"`
		Rows [][]interface{} `json:"rows"`
	} `json:"tables"`
}
type VirtualMachines []VirtualMachine
var (

	mapVirtualMachines = make(map[string]VirtualMachine)
	Virtual_Machines = VirtualMachines{}
)


func (vm *VirtualMachine) ExecuteRequest(r httpclient.IRequest) {

	//1-Three Node of Resources
	rg := ResourceGroupCost{}
	rg.ExecuteRequest(&rg)

	//2-Filters Virtual Machines only
	requests := vm.getRequests()
	requests.Execute()

	//3-Serializes All VMs and Sets Metrics
	Virtual_Machines = vm.parseRequests(requests)

	//4-Virtual Machine can be used through any output requirements
}

func (vm *VirtualMachine) GetUrl() string {
	url := azure.QueryUrl
	url = strings.ReplaceAll(url, "{{subscriptionid}}", configuration.AccessToken.SubscriptionID)
	url = strings.ReplaceAll(url, "{{locationprefix}}", vm.Resource.LocationPrefix)
	//fmt.Printf("the name: %s url is : %s\n", vm.Resource.ResourceID, url)
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
	payload = strings.ReplaceAll(payload, "{{resourcegroup}}", vm.Resource.ResourceGroup)
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
		fmt.Printf("Usage Report:\n")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Location Prefix,Consumption Type,Meter,Cost,CPU Utilization Avg,Available Memory,Logical Disk Latency,Disk IOPs,Disk Bytes/sec,Network Sent Rate, Network Received Rate")
		fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range Virtual_Machines {
			fmt.Printf("%s,%s,%s,%s,%s,%s,%s,%s,$%s,%s,%s,%s,%s,%s,%s,%s\n",item.Resource.ResourceGroup, item.Resource.ResourceID, item.Resource.Service,
				item.Resource.ServiceType, item.Resource.Location,item.Resource.LocationPrefix, item.Resource.ChargeType, item.Resource.Meter,
				item.Resource.Cost, item.CpuUtilization, item.MemoryAvailable,item.DiskLatency, item.DiskIOPs, item.DiskBytes,
				item.NetworkSentRate, item.NetworkReceivedRate)
		}
	} else {
		fmt.Printf("-")
	}
}

func (vm *VirtualMachine) setUsageValue() {

	for i := 0; i < len(vm.Tables); i++ {
		for x := 0; x < len(vm.Tables[i].Rows); x++ {
			row := vm.Tables[i].Rows[x]
			strTile := fmt.Sprintf("%v", row[0])

			//cpu
			if strings.Contains(strTile, "rocessor Time") {
				_, vm.CpuUtilization = getCpuUtilization(row)
			}

			switch strTile {
			case "Available MBytes":
				_, _, vm.MemoryAvailable = getVmAvailableMemory(row)
			case "Avg. Disk sec/Transfer":
				_, _, vm.DiskLatency = getLogicalDiskLatency(row)
			case "Disk Bytes/sec":
				_, _, vm.DiskBytes = getDiskBytesPerSeconds(row)
			case "Disk Transfers/sec":
				_, vm.DiskIOPs = getLogicalDiskIOPs(row)
			case "Bytes Sent/sec":
				_, _, vm.NetworkSentRate = getBytesSentRate(row)
			case "Bytes Received/sec":
				_, _, vm.NetworkReceivedRate = getBytesReceivedRate(row)
			}
		}
	}
}
func (vm *VirtualMachine) getRequests() httpclient.Requests {
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
					IsCache: true,
				}
				//fmt.Printf("saving request with name: %s\n", rName)
				mapVirtualMachines[rName] = *vm
				requests = append(requests, request)
			}
		}
	}
	return requests
}
func (vm *VirtualMachine) parseRequests(requests httpclient.Requests) VirtualMachines {
	vms := VirtualMachines{}
	for _, item := range requests {
		bData := item.GetResponse()
		if len(bData) > 0 {
			_ = json.Unmarshal(bData, vm)
			vmRef, hasKey := mapVirtualMachines[item.Name]
			if hasKey {
				//fmt.Printf("retrieved vm with name: %s\n", item.Name)
				//fmt.Printf("%v\n", vm.Tables)
				vm.Resource = vmRef.Resource
				vm.setUsageValue()
				vms = append(vms, *vm)
			}
		}
	}
	return vms
}
func (vm *VirtualMachine) WriteCSV(filepath string) {

	if len(Virtual_Machines) > 0 {
		var matrix [][]string
		rec := []string{"Resource Group","ResourceID","Service Name","Resource Type","Resource Location","Location Prefix","Consumption Type","Meter","Cost",
			"CPU Utilization Avg","Available Memory","Logical Disk Latency","Disk IOPs","Disk Bytes/sec","Network Sent Rate","Network Received Rate"}
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

			rec = append(rec, item.CpuUtilization)
			rec = append(rec, item.MemoryAvailable)
			rec = append(rec, item.DiskLatency)
			rec = append(rec, item.DiskIOPs)
			rec = append(rec, item.DiskBytes)
			rec = append(rec, item.NetworkSentRate)
			rec = append(rec, item.NetworkReceivedRate)
			matrix = append(matrix, rec)
		}
		csv.SaveMatrixToFile(filepath, matrix)
	}
}
//-------------------------Helper Functions Related to VM Parser-------------------------------------
// interface raw is in Kilo Bytes - need to convert to MegaBytes
func getVmAvailableMemory(row []interface{}) (float64, float64, string) {
	m := fmt.Sprintf("%v", row[12])
	kbValue, err := convert.StringToFloat(m)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", m, kbValue, err)
	}

	gbValue := kbValue / convert.UnitSymbol["GB"]
	strDisplay := fmt.Sprintf("%v", gbValue)
	strValue := fmt.Sprintf("%sGB", strDisplay[0:3])
	return gbValue, kbValue, strValue
}

func getCpuUtilization(row []interface{}) (float64, string) {
	parsedValue := fmt.Sprintf("%v", row[12])
	value, err := convert.StringToFloat(parsedValue)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", parsedValue, value, err)
	}

	strDisplay := fmt.Sprintf("%v", value)
	strValue := fmt.Sprintf("%s%%", strDisplay[0:4])
	return value, strValue
}

func getLogicalDiskLatency(row []interface{}) (float64, float64, string) {
	//the parsed value is in MS
	parsedValue := fmt.Sprintf("%v", row[12])
	value, err := convert.StringToFloat(parsedValue)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", parsedValue, value, err)
	}
	msValue := value * 1000
	strDisplay := fmt.Sprintf("%v", msValue)
	strValue := fmt.Sprintf("%sms", strDisplay[0:4])
	return msValue, value, strValue
}

func getLogicalDiskIOPs(row []interface{}) (float64, string) {
	//the parsed value is in MS
	parsedValue := fmt.Sprintf("%v", row[12])
	value, err := convert.StringToFloat(parsedValue)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", parsedValue, value, err)
	}

	strDisplay := fmt.Sprintf("%v", value)
	strValue := fmt.Sprintf("%s", strDisplay[0:4])
	return value, strValue
}

func getDiskBytesPerSeconds(row []interface{}) (float64, float64, string) {

	parsedValue := fmt.Sprintf("%v", row[12])
	value, err := convert.StringToFloat(parsedValue)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", parsedValue, value, err)
	}

	gbValue := value / convert.UnitSymbol["GB"]
	strDisplay := fmt.Sprintf("%v", value)
	strValue := fmt.Sprintf("%sGB", strDisplay[0:4])
	return gbValue, value, strValue
}

func getBytesSentRate(row []interface{}) (float64, float64, string) {

	parsedValue := fmt.Sprintf("%v", row[12])
	value, err := convert.StringToFloat(parsedValue)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", parsedValue, value, err)
	}

	kbValue := value / convert.UnitSymbol["KB"]
	strDisplay := fmt.Sprintf("%v", kbValue)
	strValue := fmt.Sprintf("%sKB", strDisplay[0:4])
	return kbValue, value, strValue
}

func getBytesReceivedRate(row []interface{}) (float64, float64, string) {

	parsedValue := fmt.Sprintf("%v", row[12])
	value, err := convert.StringToFloat(parsedValue)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", parsedValue, value, err)
	}

	kbValue := value / convert.UnitSymbol["KB"]
	strDisplay := fmt.Sprintf("%v", kbValue)
	strValue := fmt.Sprintf("%sKB", strDisplay[0:4])
	return kbValue, value, strValue
}
