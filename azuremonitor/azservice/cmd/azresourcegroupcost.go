package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


type ResourceGroupCost struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
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


func init() {

	r, err := setResourceGroupCostCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupCostCommand() (*cobra.Command, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.ResourceGroupCost.DescriptionLine1,
		cl.AppConfig.ResourceGroupCost.DescriptionLine2,
		cl.AppConfig.ResourceGroupCost.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.ResourceGroupCost.Command,
		Short: cl.AppConfig.ResourceGroupCost.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		var r = &ResourceGroupCost{}
		rgList := ResourceGroupList{}
		rgList, err := rgList.getResourceGroups()
		if err != nil {
			return err
		}

		if len(rgList) > 0 {
			clearTerminal()
			r.PrintHeader()
			for i := 0; i < len(rgList); i++ {
				rgName := rgList[i]
				r, err = r.getResourceGroupCost(rgName)
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

func (r *ResourceGroupCost) getResourceGroupCost(resourceGroupName string) (*ResourceGroupCost, error) {
	var at = &AccessToken{}
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	if len(resourceGroupName) <= 0 {
		return nil, fmt.Errorf("resource group name is required")
	}

	at, err = at.getAccessToken()
	if err != nil {
		return nil, err
	}

	startD := "2020-07-01"
	endD := "2020-07-31"
	url := strings.Replace(cl.AppConfig.ResourceGroupCost.URL, "{{subscriptionID}}",cl.AppConfig.AccessToken.SubscriptionID, 1)
	url = strings.Replace(url, "{{resourceGroup}}",resourceGroupName, 1)
	//fmt.Println(url)

	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	payload := strings.NewReader(fmt.Sprintf("{\"type\": \"ActualCost\",\"dataSet\": {\"granularity\": \"None\"," +
		"\"aggregation\": {\"totalCost\": {\"name\": \"Cost\",\"function\": \"Sum\"}," +
		"\"totalCostUSD\": {\"name\": \"CostUSD\",\"function\": \"Sum\"}}," +
		"\"grouping\": [{\"type\": \"Dimension\",\"name\": \"ResourceId\"}," +
		" {\"type\": \"Dimension\",\"name\": \"ResourceType\"}, {\"type\": \"Dimension\",\"name\": \"ResourceLocation\"}, " +
		"{\"type\": \"Dimension\",\"name\": \"ChargeType\"}, {\"type\": \"Dimension\",\"name\": \"ResourceGroupName\"}, " +
		"{\"type\": \"Dimension\",\"name\": \"PublisherType\"}, {\"type\": \"Dimension\",\"name\": \"ServiceName\"}, " +
		"{\"type\": \"Dimension\",\"name\": \"Meter\"}],\"include\": [\"Tags\"]},\"timeframe\": \"Custom\"," +
		"\"timePeriod\": {" +
		"\"from\": \"%sT00:00:00+00:00\"," +
		"\"to\": \"%sT23:59:59+00:00\"}}",
		startD,
		endD,
	))

	client := &http.Client {}
	req, _ := http.NewRequest("POST",url, payload)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))

	err = json.Unmarshal(body,r)
	if err != nil {
		fmt.Println("recommendation list unmarshal body response: ", err)
	}

	return r, nil
}

func (r ResourceGroupCost) PrintHeader() {
	fmt.Println("Resource Group Consumption:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("ResourceID,Resource Group,Service Name,Cost,Resource Type,Resource Location,Consumption Type,Meter")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
}


func (r ResourceGroupCost) Print() {
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
				pArray:= strings.Split(resourceType, "/")
				resourceType = pArray[len(pArray)-1]
			}

			if strings.Contains(resourceId, "/") {
				pArray:= strings.Split(resourceId, "/")
				resourceId = pArray[len(pArray)-1]
			}

			fmt.Printf("%s,%s,%s,$%s,%s,%s,%s,%s\n",resourceId, resourceGroupName, serviceName, costUSD,resourceType, resourceLocation, chargeType, meter )
		}


	}
}




