package costmanagement

import (
	"fmt"
	"github.com/Go/azuremonitor/azure"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/convert"
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


func (vm *VirtualMachine) ExecuteRequest(r httpclient.IRequest) {
	rg := ResourceGroupCost{}
	rg.ExecuteRequest(&rg)
	if len(Resources) > 0 {
		fmt.Printf("The resources are : %v\n", Resources)
	}
}

func (vm *VirtualMachine) GetUrl() string {
	url := azure.QueryUrl
	url = strings.ReplaceAll(url, "{{subscriptionid}}", configuration.AccessToken.SubscriptionID)
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
	if len(Resources) > 0 {
		fmt.Println("Consumption Report:")
		fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Consumption Type,Meter,Cost")
		fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
		for _, item := range Resources {
			fmt.Printf("%s,%s,%s,%s,%s,%s,$%s\n", item.ResourceGroup, item.ResourceID, item.Service, item.ServiceType, item.Location,item.Meter, item.Cost)
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
