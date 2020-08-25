package cmd

import (
	"encoding/json"
	"fmt"
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
	rootCmd.PersistentFlags().BoolVar(&saveDb, "db", false, "[=true]saves records to Postgres db")
	rootCmd.PersistentFlags().BoolVar(&saveCsv, "csv", false, "[=true]saves records into a csv output file")
	rootCmd.PersistentFlags().BoolVar(&ignoreZeroCost, "izcost", false, "[=true] ignores resources with zero cost")

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
		errors := requests.Execute()
		IfErrorsPrintThem(errors)

		if saveCsv {
			RemoveFile(csvRgcReportName)
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

func (r *ResourceGroupCost) getRequests(rsgroups []string) Requests {
	requests := Requests{}
	header := r.getHeader()
	for i := 0; i < len(rsgroups); i++ {
		rgName := rsgroups[i]
		request := Request{}
		request.Name = rgName
		request.Header = header
		request.Payload = r.getPayload()
		request.Url = r.getUrl(rgName)
		request.Method = Methods.POST
		request.IsCache = true
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
func (r *ResourceGroupCost) getHeader() http.Header {
	var at = &AccessToken{}
	at.ExecuteRequest(at)
	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	var header = http.Header{}
	header.Add("Authorization", token)
	header.Add("Accept", "application/json")
	header.Add("Content-Type", "application/json")
	return header
}

func (r ResourceGroupCost) PrintHeader() {
	fmt.Println("Consumption Report:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("Resource Group,ResourceID,Service Name,Resource Type,Resource Location,Consumption Type,Meter,Cost")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	if saveCsv {
		var matrix [][]string
		rec := []string{"Resource Group", "ResourceID", "Service Name", "Resource Type", "Resource Location", "Consumption Type", "Meter", "Cost"}
		matrix = append(matrix, rec)
		saveCSV(csvRgcReportName, matrix)
	}
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

			fmt.Printf("\t%s,%s,%s,%s,%s,%s,$%s\n", resourceId, serviceName, resourceType, resourceLocation, chargeType, meter, costUSD)
		}
	}
}

func (r ResourceGroupCost) writeCSV() {

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

				var rec []string
				rec = append(rec, r.ResourceGroupName)
				rec = append(rec, resourceId)
				rec = append(rec, serviceName)
				rec = append(rec, resourceType)
				rec = append(rec, resourceLocation)
				rec = append(rec, chargeType)
				rec = append(rec, meter)
				rec = append(rec, costUSD)
				matrix = append(matrix, rec)
			}
		}

		saveCSV(csvRgcReportName, matrix)
	}

}
