package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

func (r *ResourceUsageVirtualMachine) getVirtualMachineByResourceId(id string, startD string,endD string) (*ResourceUsageVirtualMachine, error) {
	c := &Cache{}


	//Validate
	if id == "" || startD == "" || endD == "" {
		return nil, fmt.Errorf("resource id name is required")
	}

	//Cache lookup
	cKey := fmt.Sprintf("GetVirtualMachineByResourceId_%s_%s_%s",id, startD, endD)
	cHashVal := c.Get(cKey)
	if len(cHashVal) <= 0 {
		//Execute Request
		r, err := r.executeRequest(id, startD, endD, cKey)
		if err != nil {
			return r, err
		}

	} else {
		//Load From Cache
		err := LoadFromCache(cKey, r)
		if err != nil {
			fmt.Println("******WARNNING!!!!!!!!!MISSING FILE:::RESTORING WITH NEW REQUEST:::", err)
			r, err := r.executeRequest(id, startD, endD, cKey)
			if err != nil {
				return r, err
			}
		}
		//fmt.Println(r)
	}

	return r, nil
}

func (r *ResourceUsageVirtualMachine) executeRequest(id string, startD string, endD string, cKey string) (*ResourceUsageVirtualMachine, error) {

	var at = &AccessToken{}
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	at, err = at.getAccessToken()
	if err != nil {
		return nil, err
	}

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
		return r, fmt.Errorf("recommendation list unmarshal body response: ", err)
	}

	//cached it
	err = saveCache(cKey, r)
	if err != nil {
		return r, fmt.Errorf("error: failed to save to cache folder - %s: %v", cKey, err)
	}

	return r, nil
}

func (r ResourceUsageVirtualMachine) Print() {

	//var availableMemory float64
	for i := 0; i < len(r.Tables); i++ {
		for x := 0; x < len(r.Tables[i].Rows); x++ {
			row := r.Tables[i].Rows[x]
			strTile := fmt.Sprintf("%v", row[0])
			switch strTile {
			case "Available MBytes":
				getVmAvailableMemory(row)
			case "% Processor Time":

			case "Avg. Disk sec/Transfer":

			case "Disk Bytes/sec":

			case "Disk Transfers/sec":

			case "Bytes Sent/sec":

			case "Bytes Received/sec":
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
}

// interface raw is in Kilo Bytes - need to convert to MegaBytes
func getVmAvailableMemory(row []interface{}) (float64, float64) {
	m := fmt.Sprintf("%v", row[12])
	kbValue, err := stringToFloat(m)
	if err != nil {
		fmt.Printf("%q\t %g %v\n", m, kbValue, err)
	}

	gbValue := kbValue / GB
	strDisplay := fmt.Sprintf("%v", gbValue)
	fmt.Printf("Available Memory: %sGB [%gKB] \n", strDisplay[0:3], kbValue)
	return gbValue, kbValue
}




