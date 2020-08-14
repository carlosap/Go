package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/db/dbcontext"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/Go/azuremonitor/db/cache"
)

type ResourceGroupCost struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	ResourceGroupName       string      `json:"resourcegroupname"`
	Type       string      `json:"type"`
	Location   interface{} `json:"location"`
	Sku        interface{} `json:"sku"`
	ETag       interface{} `json:"eTag"`
	Properties struct {
		NextLink interface{} `json:"nextLink"`
		Columns  []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"columns"`
		Rows [][]interface{} `json:"rows"`
	} `json:"properties"`
}

var (
	layoutISO = "2006-01-02"
	startDate string
	endDate   string
)

func init() {
	r, err := setResourceGroupCostCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	now := time.Now()
	month := now.AddDate(0, 0, -29)
	//make sure we support a syntax like this  .\azservice.exe get-rgc --from 2020-07-01 --to 2020-07-30
	rootCmd.PersistentFlags().StringVar(&startDate, "from", month.Format("2006-01-02"), "start date of report (i.e. YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&endDate, "to", now.Format("2006-01-02"), "end date of report (i.e. YYYY-MM-DD)")
	rootCmd.AddCommand(r)
}

func setResourceGroupCostCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		cmdConfig.ResourceGroupCost.DescriptionLine1,
		cmdConfig.ResourceGroupCost.DescriptionLine2,
		cmdConfig.ResourceGroupCost.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cmdConfig.ResourceGroupCost.Command,
		Short: cmdConfig.ResourceGroupCost.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		var r = &ResourceGroupCost{}
		rgList := ResourceGroupList{}
		rgList, err := rgList.getResourceGroups()
		if err != nil {
			return err
		}

		clearTerminal()
		if len(rgList) > 0 {

			r.PrintHeader()
			for i := 0; i < len(rgList); i++ {
				rgName := rgList[i]
				r, err = r.getResourceGroupCost(rgName, startDate, endDate)
				if err != nil {
					return err
				}
				r.Print()
			}

			for i := 0; i < len(rgList); i++ {
				rgName := rgList[i]
				r, err = r.getResourceGroupCost(rgName, startDate, endDate)
				if err != nil {
					return err
				}
				r.PrintUsage()
			}

		}

		return nil
	}
	return cmd, nil
}

func (r *ResourceGroupCost) getResourceGroupCost(resourceGroupName string, startD string, endD string) (*ResourceGroupCost, error) {

	if resourceGroupName == "" || startD == "" || endD == "" {
		fmt.Println("error: resource group cost function requires resource group")
	}

	r.ResourceGroupName = resourceGroupName

	//Cache lookup
	c := &cache.Cache{}
	cKey := fmt.Sprintf("%s_%s_GetResourceGroupCost_%s_%s", cmdConfig.AccessToken.SubscriptionID, resourceGroupName, startD, endD)
	cHashVal := c.Get(cKey)
	if len(cHashVal) <= 0 {
		//Execute Request
		r, err := r.executeRequest(resourceGroupName, startD, endD, cKey, cmdConfig.AccessToken.SubscriptionID)
		if err != nil {
			return r, err
		}

	} else {
		//Load From Cache
		err := LoadFromCache(cKey, r)
		if err != nil {
			r, err := r.executeRequest(resourceGroupName, startD, endD, cKey, cmdConfig.AccessToken.SubscriptionID)
			if err != nil {
				return r, err
			}
		}
	}

	return r, nil
}

func (r *ResourceGroupCost) executeRequest(resourceGroupName string, startD string, endD string, cKey string, subscriptionId string) (*ResourceGroupCost, error) {
	var at = &AccessToken{}

	at, err := at.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := strings.Replace(cmdConfig.ResourceGroupCost.URL, "{{subscriptionID}}", subscriptionId, 1)
	url = strings.Replace(url, "{{resourceGroup}}", resourceGroupName, 1)

	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	payload := strings.NewReader(fmt.Sprintf("{\"type\": \"ActualCost\",\"dataSet\": {\"granularity\": \"None\","+
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
		startD,
		endD,
	))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, r)
	if err != nil {
		return r, fmt.Errorf("recommendation list unmarshal body response: ", err)
	}

	//cached it
	err = saveCache(cKey, r)
	if err != nil {
		return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
	}

	return r, nil
}

func (r ResourceGroupCost) PrintHeader() {
	fmt.Println("Consumption Report:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Consumption Type,Meter,Cost")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")

}


func (r ResourceGroupCost) PrintUsage() {

	printResourceGroupUsage(r)
}

func (r ResourceGroupCost) Print() {

	printResourceGroupCost(r)
}

func printResourceGroupCost(r ResourceGroupCost) {
	fmt.Printf("%s\n",r.ResourceGroupName)
	for i := 0; i < len(r.Properties.Rows); i++ {
		row := r.Properties.Rows[i]
		if len(row) > 0 {
			//casting interface to string
			costUSD := fmt.Sprintf("%v", row[1])
			resourceId := fmt.Sprintf("%v", row[2])
			resourceType := fmt.Sprintf("%v", row[3])
			resourceLocation := fmt.Sprintf("%v", row[4])
			chargeType := fmt.Sprintf("%v", row[5])
			//resourceGroupName := fmt.Sprintf("%v", row[6])
			//publisherType := fmt.Sprintf("%v", row[7])
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

			fmt.Printf("\t%s,%s,%s,%s,%s,%s,$%s\n",resourceId,serviceName,resourceType, resourceLocation, chargeType, meter,costUSD)
		}
	}
}

func printResourceGroupUsage(r ResourceGroupCost) {
	fmt.Println("")
	for i := 0; i < len(r.Properties.Rows); i++ {
		row := r.Properties.Rows[i]
		//fmt.Printf("%v\n", row)
		if len(row) > 0 {
			//casting interface to string
			costUSD := fmt.Sprintf("%v", row[1])
			resourceId := fmt.Sprintf("%v", row[2])
			resourceType := fmt.Sprintf("%v", row[3])
			resourceLocation := fmt.Sprintf("%v", row[4])
			chargeType := fmt.Sprintf("%v", row[5])
			resourceGroupName := fmt.Sprintf("%v", row[6])
			//publisherType := fmt.Sprintf("%v", row[7])
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
