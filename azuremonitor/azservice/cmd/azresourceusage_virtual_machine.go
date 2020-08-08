package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ResourceUsageVirtualMachine struct {
	Tables []struct {
		Name    string `json:"name"`
		Columns []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"columns"`
		Rows [][]interface{} `json:"rows"`
	} `json:"tables"`
}

const (
	TB = 1000000000000
	GB = 1000000000
	MB = 1000000
	KB = 1000
)

var siFactors = map[string]float64{
	"":  1e0,
	"k": 1e3,
	"M": 1e6, // Sometimes, M (Roman numeral) for thousands and MM for millions
	"G": 1e9,
	"T": 1e12,
	"P": 1e15,
	"E": 1e18,
	"Z": 1e21,
	"Y": 1e24,
	"K": 1e3, // colloquial synonym for "k"
	"B": 1e9, // colloquial synonym for "G"
}

func parseNumber(s string) (float64, error) {
	fmt.Println("the number sent: ", s)
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f, nil
	}
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError {
		return 0, err
	}
	symbol := s[len(s)-size : len(s)]
	factor, ok := siFactors[symbol]
	if !ok {
		return 0, err
	}
	f, e := strconv.ParseFloat(s[:len(s)-len(symbol)], 64)
	if e != nil {
		return 0, err
	}
	return f * factor, nil
}

func (r *ResourceUsageVirtualMachine) getVirtualMachineByResourceId(id string) (*ResourceUsageVirtualMachine, error) {
	var at = &AccessToken{}
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	if len(id) <= 0 {
		return nil, fmt.Errorf("resource id name is required")
	}

	at, err = at.getAccessToken()
	if err != nil {
		return nil, err
	}

	// only 30-days increments
	startD := "2020-07-01"
	endD := "2020-07-30"
	url := fmt.Sprintf("https://management.azure.com//subscriptions/%s/resourcegroups/" +
		"defaultresourcegroup-eus/providers/microsoft.operationalinsights/workspaces/" +
		"defaultworkspace-%s-eus/query?api-version=2017-10-01",cl.AppConfig.AccessToken.SubscriptionID, cl.AppConfig.AccessToken.SubscriptionID)



	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	payload := strings.NewReader(fmt.Sprintf("{\"query\": \"let " +
		"startDateTime = datetime('%sT08:00:00.000Z');" +
		"let endDateTime = datetime('%sT16:00:00.000Z');" +
		"let trendBinSize = 8h;" +
		"let maxListSize = 1000;" +
		"let cpuMemory = materialize(InsightsMetrics| where TimeGenerated between (startDateTime .. endDateTime)| " +
		"where _ResourceId =~ '%s'| " +
		"where Origin == 'vm.azm.ms'| where (Namespace == 'Processor' and Name == 'UtilizationPercentage') or (Namespace == 'Memory' and Name == 'AvailableMB')| " +
		"project TimeGenerated, Name, Namespace, Val);let networkDisk = materialize(InsightsMetrics| where TimeGenerated between (startDateTime .. endDateTime)| " +
		"where _ResourceId =~ '%s'| " +
		"where Origin == 'vm.azm.ms'| " +
		"where (Namespace == 'Network' and Name in ('WriteBytesPerSecond', 'ReadBytesPerSecond'))    " +
		"or (Namespace == 'LogicalDisk' and Name in ('TransfersPerSecond', 'BytesPerSecond', 'TransferLatencyMs'))| " +
		"extend ComputerId = iff(isempty(_ResourceId), Computer, _ResourceId)| " +
		"summarize Val = sum(Val) by bin(TimeGenerated, 1m), ComputerId, Name, Namespace| project TimeGenerated, Name, Namespace, Val);" +
		"let rawDataCached = cpuMemory| union networkDisk| extend Val = iif(Name in ('WriteLatencyMs', 'ReadLatencyMs', 'TransferLatencyMs'), Val/1000.0, Val)|" +
		" project TimeGenerated,cName = case(Namespace == 'Processor' and Name == 'UtilizationPercentage', '% Processor Time'," +
		"Namespace == 'Memory' and Name == 'AvailableMB','Available MBytes',Namespace == 'LogicalDisk' and Name == 'TransfersPerSecond', 'Disk Transfers/sec'," +
		"Namespace == 'LogicalDisk' and Name == 'BytesPerSecond', 'Disk Bytes/sec',Namespace == 'LogicalDisk' " +
		"and Name == 'TransferLatencyMs', 'Avg. Disk sec/Transfer',Namespace == 'Network' " +
		"and Name == 'WriteBytesPerSecond', 'Bytes Sent/sec',Namespace == 'Network' " +
		"and Name == 'ReadBytesPerSecond', 'Bytes Received/sec',Name)," +
		"cValue = case(Val < 0, real(0),Val);rawDataCached| summarize min(cValue),avg(cValue),max(cValue)," +
		"percentiles(cValue, 5, 10, 50, 90, 95) by bin(TimeGenerated, trendBinSize), cName| " +
		"sort by TimeGenerated asc| summarize makelist(TimeGenerated, maxListSize)," +
		"makelist(min_cValue, maxListSize),makelist(avg_cValue, maxListSize),makelist(max_cValue, maxListSize),makelist(percentile_cValue_5, maxListSize)," +
		"makelist(percentile_cValue_10, maxListSize),makelist(percentile_cValue_50, maxListSize),makelist(percentile_cValue_90, maxListSize)," +
		"makelist(percentile_cValue_95, maxListSize) " +
		"by cName| join(rawDataCached    | summarize min(cValue), avg(cValue), max(cValue), " +
		"percentiles(cValue, 5, 10, 50, 90, 95) by cName)on cName\"," +
		"\"timespan\": \"%sT08:00:00.000Z/%sT16:00:00.000Z\"}",
		startD,
		endD,
		id,
		id,
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

func (r ResourceUsageVirtualMachine) Print() {
	dir := "cache"
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Println("creating cache directory")
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Printf("os.Mkdir('%s') failed with '%s'\n", dir)
		}
	}



	var strAvailableMemoryBytes  string
	//var availableMemory float64
	for i:= 0; i < len(r.Tables); i++ {
		for x:= 0; x < len(r.Tables[i].Rows); x++ {
			row := r.Tables[i].Rows[x]
			switch x {
			case 0:
				//Raw is in Kilo Bytes - need to convert to MegaBytes
				strTile := fmt.Sprintf("%v", row[0])
				strAvailableMemoryBytes = fmt.Sprintf("%v", row[12])
				n, err := parseNumber(strAvailableMemoryBytes)
				if err != nil {
					fmt.Printf("%q\t %g %v\n", strAvailableMemoryBytes, n, err)
				}


				availableMemoryGb := n / GB
				strDisplay := fmt.Sprintf("%v", availableMemoryGb)
				strDisplay = strDisplay[0:3]
				fmt.Printf("Available Memory: %s - %sGB [%gKB] \n", strTile,strDisplay, n) // round down

			case 1:

			}

			//cName := fmt.Sprintf("%v", row[0])
			//timeGenerated := fmt.Sprintf("%v", row[1])
			//minValue := fmt.Sprintf("%v", row[2])
			//avgValue := fmt.Sprintf("%v", row[3])
			//maxValue := fmt.Sprintf("%v", row[4])
			//percentileValue_five := fmt.Sprintf("%v", row[5])
			//percentileValue_ten := fmt.Sprintf("%v", row[6])
			//percentileValue_fifty := fmt.Sprintf("%v", row[7])
			//percentileValue_ninety := fmt.Sprintf("%v", row[8])
			//percentile_Value_ninety_five := fmt.Sprintf("%v", row[9])

			//summaryName := fmt.Sprintf("%v", row[10])
			//sMinValue := fmt.Sprintf("%v", row[11])
			//savgValue := fmt.Sprintf("%vMb", row[12])
			//smaxValue := fmt.Sprintf("%v", row[13])
			//sPercentileValue_five := fmt.Sprintf("%v", row[14])
			//sPercentileValue_ten := fmt.Sprintf("%v", row[15])
			//sPercentileValue_fifty := fmt.Sprintf("%v", row[16])
			//sPercentileValue_ninety := fmt.Sprintf("%v", row[17])
			//sPercentile_Value_ninety_five := fmt.Sprintf("%v", row[18])
		}
	}



}




