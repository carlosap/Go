package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/azure/batch"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/csv"
	"github.com/Go/azuremonitor/common/filesystem"
	"github.com/Go/azuremonitor/common/httpclient"
	"github.com/Go/azuremonitor/common/terminal"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

type ResourceGroupUsage struct {
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

func init() {
	r, err := setResourceGroupUsageCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupUsageCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.ResourceGroupUsage.DescriptionLine1,
		configuration.ResourceGroupUsage.DescriptionLine2,
		configuration.ResourceGroupUsage.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.ResourceGroupUsage.Command,
		Short: configuration.ResourceGroupUsage.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		var r = &ResourceGroupUsage{}
		rgl := batch.ResourceGroupList{}
		rgl.ExecuteRequest(&rgl)

		requests := r.getRequests(rgl)
		_= requests.Execute()
		//IfErrorsPrintThem(errors)

		if saveCsv {
			filesystem.RemoveFile(csvRguReportName)
			r.PrintHeader()
		}

		for _, item := range requests {
			if len(item.GetResponse()) > 0 {
				_ = json.Unmarshal(item.GetResponse(), r)
				r.ResourceGroupName = item.Name
				r.Print()
				r.writeCSV()
			}
		}

		return nil
	}
	return cmd, nil
}

func (r *ResourceGroupUsage) getRequests(rsgroups []string) httpclient.Requests {
	requests := httpclient.Requests{}
	header := r.getHeader()
	for i := 0; i < len(rsgroups); i++ {
		rgName := rsgroups[i]
		request := httpclient.Request{}
		request.Name = rgName
		request.Header = header
		request.Payload = r.getPayload()
		request.Url = r.getUrl(rgName)
		request.Method = httpclient.Methods.POST
		request.IsCache = true
		requests = append(requests, request)

	}
	return requests
}
func (r *ResourceGroupUsage) getUrl(resourceGroupName string) string {
	url := strings.Replace(configuration.ResourceGroupUsage.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	url = strings.Replace(url, "{{resourceGroup}}", resourceGroupName, 1)
	return url
}
func (r *ResourceGroupUsage) getPayload() string {
	return fmt.Sprintf("{\"type\": \"ActualCost\",\"dataSet\": {\"granularity\": \"None\","+
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
		startDate,
		endDate,
	)
}
func (r *ResourceGroupUsage) getHeader() http.Header {
	at := &oauth2.AccessToken{}
	at.ExecuteRequest(at)
	var header = http.Header{}

	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}
func (r ResourceGroupUsage) PrintHeader() {
	fmt.Printf("Usage Report:\n")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("ResourceID,Resource Group,Service Name,Cost,Resource Type,Resource Location,Consumption Type,Meter,CPU Utilization Avg,Available Memory,Logical Disk Latency,Disk IOPs,Disk Bytes/sec,Network Sent Rate, Network Received Rate")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
	if saveCsv {
		var matrix [][]string
		rec := []string{"ResourceID", "Resource Group", "Service Name", "Cost", "Resource Type", "Resource Location", "Consumption Type", "Meter", "CPU Utilization Avg", "Available Memory", "Logical Disk Latency", "Disk IOPs", "Disk Bytes/sec", "Network Sent Rate", "Network Received Rate"}
		matrix = append(matrix, rec)
		csv.SaveMatrixToFile(csvRguReportName, matrix)
	}
}

func (r ResourceGroupUsage) Print() {
	for i := 0; i < len(r.Properties.Rows); i++ {
		row := r.Properties.Rows[i]
		if len(row) > 0 {
			costUSD := fmt.Sprintf("%v", row[1])
			resourceId := fmt.Sprintf("%v", row[2])
			resourceType := fmt.Sprintf("%v", row[3])
			resourceLocation := fmt.Sprintf("%v", row[4])
			chargeType := fmt.Sprintf("%v", row[5])
			resourceGroupName := fmt.Sprintf("%v", row[6])
			serviceName := fmt.Sprintf("%v", row[8])
			meter := fmt.Sprintf("%v", row[9])

			//format cost
			if len(costUSD) > 5 {
				costUSD = costUSD[0:5]
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

			//Additional requests
			if serviceName == "virtual machines" && resourceType == "virtualmachines" && len(costUSD) > 0 && chargeType == "usage" {

				var vm = &ResourceUsageVirtualMachine{}
				vm, err := vm.getVmUsage(resourceGroupName, resourceId)
				if err != nil {
					fmt.Printf("Error: failed to retrieve vm resouce usage %v\n", err)
				}

				fmt.Printf("%s,%s,%s,$%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", resourceId, resourceGroupName, serviceName, costUSD, resourceType, resourceLocation, chargeType, meter, vm.CpuUtilization, vm.MemoryAvailable, vm.DiskLatency, vm.DiskIOPs, vm.DiskBytes, vm.NetworkSentRate, vm.NetworkSentRate)

			}
		}
	}
}

func (r ResourceGroupUsage) writeCSV() {

	if saveCsv {
		var matrix [][]string
		for i := 0; i < len(r.Properties.Rows); i++ {
			row := r.Properties.Rows[i]
			if len(row) > 0 {
				costUSD := fmt.Sprintf("%v", row[1])
				resourceId := fmt.Sprintf("%v", row[2])
				resourceType := fmt.Sprintf("%v", row[3])
				resourceLocation := fmt.Sprintf("%v", row[4])
				chargeType := fmt.Sprintf("%v", row[5])
				resourceGroupName := fmt.Sprintf("%v", row[6])
				serviceName := fmt.Sprintf("%v", row[8])
				meter := fmt.Sprintf("%v", row[9])

				//format cost
				if len(costUSD) > 5 {
					costUSD = costUSD[0:5]
				}

				if ignoreZeroCost {
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

				//Additional requests
				if serviceName == "virtual machines" && resourceType == "virtualmachines" && len(costUSD) > 0 && chargeType == "usage" {

					var vm = &ResourceUsageVirtualMachine{}
					vm, err := vm.getVmUsage(resourceGroupName, resourceId)
					if err != nil {
						fmt.Printf("Error: failed to retrieve vm resouce usage %v\n", err)
					}

					var rec []string
					rec = append(rec, resourceId)
					rec = append(rec, resourceGroupName)
					rec = append(rec, serviceName)
					rec = append(rec, costUSD)
					rec = append(rec, resourceType)
					rec = append(rec, resourceLocation)
					rec = append(rec, chargeType)
					rec = append(rec, meter)
					rec = append(rec, vm.CpuUtilization)
					rec = append(rec, vm.MemoryAvailable)
					rec = append(rec, vm.DiskLatency)
					rec = append(rec, vm.DiskIOPs)
					rec = append(rec, vm.DiskBytes)
					rec = append(rec, vm.NetworkSentRate)
					rec = append(rec, vm.NetworkSentRate)
					matrix = append(matrix, rec)
				}
			}
		}
		csv.SaveMatrixToFile(csvRguReportName, matrix)
	}
}

//TODO:::::
//if serviceName == "storage" && resourceType == "storageaccounts" && chargeType == "usage" {
//	var stacc = &StorageAccountAvailability{}
//	stacc, err := stacc.getStorageAccountAvailability(resourceGroupName, resourceId, "2020-08-01", "2020-08-07")
//	if err != nil {
//		fmt.Printf("Error: failed to retrieve Availability resouce usage %v\n", err)
//	}
//
//	//var transaction = &StorageAccountTransaction{}
//	//transaction, err := transaction.getStorageAccountTransaction(resourceGroupName, resourceId, "2020-08-01", "2020-08-07")
//	//if err != nil {
//	//	fmt.Printf("Error: failed to retrieve Transaction resouce usage %v\n", err)
//	//}
//
//	fmt.Printf("Resource Group Consumption: %s-%s\n", serviceName, resourceType)
//	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
//	fmt.Println("ResourceID,Resource Group,Service Name,Cost,Resource Type,Resource Location,Consumption Type,Meter,Availability, Total Transactions,E2E Latency, Server Lantency, Failures,Capacity")
//	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
//	fmt.Printf("%s,%s,%s,$%s,%s,%s,%s,%s,%g%%,%g\n",
//		resourceId, resourceGroupName, serviceName, costUSD,
//		resourceType, resourceLocation, chargeType, meter,
//		stacc.getAvailability(),300.0) //transaction.getTransactions())
//}

//var vmContext = &dbcontext.Virtualmachine{}
//vmContext.Resourceid = &resourceId
//vmContext.Resourcegroup = &resourceGroupName
//vmContext.Servicename = &serviceName
//vmContext.Cost = &costUSD
//vmContext.Resourcetype = &resourceType
//vmContext.Meter = &meter
//vmContext.Cpuutilization = &vm.CpuUtilization
//vmContext.Availablememory = &vm.MemoryAvailable
//vmContext.Disklatency = &vm.DiskLatency
//vmContext.Diskiops = &vm.DiskIOPs
//vmContext.Diskbytespersec = &vm.DiskBytes
//vmContext.Networksentrate = &vm.NetworkSentRate
//vmContext.Networkreceivedrate = &vm.NetworkReceivedRate
//vmContext.Resourcelocation = &resourceLocation
//vmContext.Consumptiontype = &chargeType
//vmContext.Reportstartdate = &startDate
//vmContext.Reportenddate = &endDate
//var dataDictionary map[string]interface{}
//d, _ := json.Marshal(&r)
//_ = json.Unmarshal(d, &dataDictionary)
//vmContext.Data = dataDictionary
//
//err = vmContext.Insert()
//if err != nil {
//	fmt.Printf("Error: while inserting vm record %v", err)
//}
