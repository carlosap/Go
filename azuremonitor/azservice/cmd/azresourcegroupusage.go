package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/db/cache"
	"github.com/Go/azuremonitor/db/dbcontext"
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
		var r = &ResourceGroupUsage{}
		rgList := ResourceGroupList{}
		rgList, err := rgList.getResourceGroups()
		if err != nil {
			return err
		}

		clearTerminal()
		if len(rgList) > 0 {
			for i := 0; i < len(rgList); i++ {
				rgName := rgList[i]
				r, err = r.getResourceGroupUsage(rgName)
				if err != nil {
					return err
				}
				r.Print()
			}
		}
		return nil
	}
	return cmd, nil
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
func (r *ResourceGroupUsage) getHeader() (http.Header, error) {
	var at = &AccessToken{}
	var header = http.Header{}
	at, err := at.getAccessToken()
	if err != nil {
		return nil, err
	}
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header, err
}

func (r *ResourceGroupUsage) getResourceGroupUsage(resourceGroupName string) (*ResourceGroupUsage, error) {

	if resourceGroupName == "" {
		fmt.Println("error: resource group cost function requires resource group")
	}

	r.ResourceGroupName = resourceGroupName
	url := r.getUrl(resourceGroupName)
	payload := r.getPayload()
	header, _ := r.getHeader()
	c := &cache.Cache{}
	cKey := fmt.Sprintf("%s_%s_GetResourceGroupUsage_%s_%s", configuration.AccessToken.SubscriptionID, resourceGroupName, startDate, endDate)
	cHashVal := c.Get(cKey)
	if len(cHashVal) <= 0 {
		request := Request{
			"ResourceGroups",
			url,
			Methods.POST,
			payload,
			header,
		}
		_ = request.Execute()
		body := request.GetResponse()
		err := json.Unmarshal(body, r)
		if err != nil {
			fmt.Println("GetResourceGroupUsage unmarshal body response: ", err)
		}
		err = saveCache(cKey, r)
		if err != nil {
			return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
		}

	} else {
		//Load From Cache
		err := LoadFromCache(cKey, r)
		if err != nil {
			request := Request{
				"ResourceGroups",
				url,
				Methods.POST,
				payload,
				header,
			}
			_ = request.Execute()
			body := request.GetResponse()
			err := json.Unmarshal(body, r)
			if err != nil {
				fmt.Println("GetResourceGroupUsage unmarshal body response: ", err)
			}
			err = saveCache(cKey, r)
			if err != nil {
				return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
			}
		}
	}

	return r, nil
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
				var vmContext = &dbcontext.Virtualmachine{}
				var vm = &ResourceUsageVirtualMachine{}
				vm, err := vm.getVirtualMachineByResourceId(resourceGroupName, resourceId)
				if err != nil {
					fmt.Printf("Error: failed to retrieve vm resouce usage %v\n", err)
				}

				fmt.Printf("Usage Report:\n")
				fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
				fmt.Println("ResourceID,Resource Group,Service Name,Cost,Resource Type,Resource Location,Consumption Type,Meter,CPU Utilization Avg,Available Memory,Logical Disk Latency,Disk IOPs,Disk Bytes/sec,Network Sent Rate, Network Received Rate")
				fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
				fmt.Printf("%s,%s,%s,$%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n", resourceId, resourceGroupName, serviceName, costUSD, resourceType, resourceLocation, chargeType, meter, vm.CpuUtilization, vm.MemoryAvailable, vm.DiskLatency, vm.DiskIOPs, vm.DiskBytes, vm.NetworkSentRate, vm.NetworkSentRate)
				vmContext.Resourceid = &resourceId
				vmContext.Resourcegroup = &resourceGroupName
				vmContext.Servicename = &serviceName
				vmContext.Cost = &costUSD
				vmContext.Resourcetype = &resourceType
				vmContext.Meter = &meter
				vmContext.Cpuutilization = &vm.CpuUtilization
				vmContext.Availablememory = &vm.MemoryAvailable
				vmContext.Disklatency = &vm.DiskLatency
				vmContext.Diskiops = &vm.DiskIOPs
				vmContext.Diskbytespersec = &vm.DiskBytes
				vmContext.Networksentrate = &vm.NetworkSentRate
				vmContext.Networkreceivedrate = &vm.NetworkReceivedRate
				vmContext.Resourcelocation = &resourceLocation
				vmContext.Consumptiontype = &chargeType
				vmContext.Reportstartdate = &startDate
				vmContext.Reportenddate = &endDate
				var dataDictionary map[string]interface{}
				d, _ := json.Marshal(&r)
				_ = json.Unmarshal(d, &dataDictionary)
				vmContext.Data = dataDictionary

				err = vmContext.Insert()
				if err != nil {
					fmt.Printf("Error: while inserting vm record %v", err)
				}
			}

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

		}
	}
}

