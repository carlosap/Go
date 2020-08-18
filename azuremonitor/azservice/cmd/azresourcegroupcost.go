package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Go/azuremonitor/db/cache"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
	"time"
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

func init() {

	now := time.Now()
	month := now.AddDate(0, 0, -29)
	rootCmd.PersistentFlags().StringVar(&startDate, "from", month.Format(layoutISO), "start date of report (i.e. YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&endDate, "to", now.Format(layoutISO), "end date of report (i.e. YYYY-MM-DD)")

	r, err := setResourceGroupCostCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupCostCommand() (*cobra.Command, error) {

	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.ResourceGroupCost.DescriptionLine1,
		configuration.ResourceGroupCost.DescriptionLine2,
		configuration.ResourceGroupCost.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.ResourceGroupCost.Command,
		Short: configuration.ResourceGroupCost.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		var r = &ResourceGroupCost{}
		rgList := ResourceGroupList{}

		rgList, err := rgList.getResourceGroups()
		if err != nil {
			return err
		}

		clearTerminal()
		requests := r.getRequests(rgList)
		_ = requests.Execute()

		for _, item := range requests {
			fmt.Printf("The item: %s\n", item.Name)
			if len(item.GetResponse()) > 0 {
				err := json.Unmarshal(item.GetResponse(), r)
				if err != nil {
					fmt.Println("GetResourceGroupCost unmarshal body response: ", err)
				}
				r.Print()
			}
			fmt.Println()
		}

		//fmt.Printf("the end of the line : %v\n", requests)
		//if len(rgList) > 0 {
		//
		//	r.PrintHeader()
		//	for i := 0; i < len(rgList); i++ {
		//		rgName := rgList[i]
		//		r, err = r.getResourceGroupCost(rgName)
		//		if err != nil {
		//			return err
		//		}
		//		r.Print()
		//	}
		//}

		return nil
	}
	return cmd, nil
}

func (r *ResourceGroupCost) getRequests(rsgroups []string) Requests {
	requests := Requests{}
	header, _ := r.getHeader()
	for i := 0; i < len(rsgroups); i++ {
		rgName := rsgroups[i]
		request := Request{}
		request.Name = rgName
		request.Header = header
		request.Payload = r.getPayload()
		request.Url = r.getUrl(rgName)
		request.Method = Methods.POST
		requests = append(requests, request)
	}
	return requests
}

func (r *ResourceGroupCost) getUrl(resourceGroupName string) string {
	url := strings.Replace(configuration.ResourceGroupCost.URL, "{{subscriptionID}}", configuration.AccessToken.SubscriptionID, 1)
	url = strings.Replace(url, "{{resourceGroup}}", resourceGroupName, 1)
	return url
}
func (r *ResourceGroupCost) getPayload() string {
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
func (r *ResourceGroupCost) getHeader() (http.Header, error) {
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

func (r *ResourceGroupCost) getResourceGroupCost(resourceGroupName string) (*ResourceGroupCost, error) {

	if resourceGroupName == "" {
		fmt.Println("error: resource group cost function requires resource group")
	}

	r.ResourceGroupName = resourceGroupName
	url := r.getUrl(resourceGroupName)
	payload := r.getPayload()
	header, _ := r.getHeader()
	//Cache lookup
	c := &cache.Cache{}
	cKey := fmt.Sprintf("%s_%s_GetResourceGroupCost_%s_%s", configuration.AccessToken.SubscriptionID, resourceGroupName, startDate, endDate)
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
			fmt.Println("GetResourceGroupCost unmarshal body response: ", err)
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
				fmt.Println("GetResourceGroupCost unmarshal body response: ", err)
			}
			err = saveCache(cKey, r)
			if err != nil {
				return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
			}
		}
	}
	return r, nil
}

func (r ResourceGroupCost) PrintHeader() {
	fmt.Println("Consumption Report:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Consumption Type,Meter,Cost")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
}

func (r ResourceGroupCost) Print() {
	fmt.Printf("%s\n", r.ResourceGroupName)
	for i := 0; i < len(r.Properties.Rows); i++ {
		row := r.Properties.Rows[i]
		if len(row) > 0 {
			//casting interface to string
			costUSD := fmt.Sprintf("%v", row[1])
			resourceId := fmt.Sprintf("%v", row[2])
			resourceType := fmt.Sprintf("%v", row[3])
			resourceLocation := fmt.Sprintf("%v", row[4])
			chargeType := fmt.Sprintf("%v", row[5])
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

			fmt.Printf("\t%s,%s,%s,%s,%s,%s,$%s\n", resourceId, serviceName, resourceType, resourceLocation, chargeType, meter, costUSD)
		}
	}
}
