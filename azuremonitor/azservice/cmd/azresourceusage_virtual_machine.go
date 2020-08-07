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

	fmt.Println(url)

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
	fmt.Println(string(body))

	err = json.Unmarshal(body,r)
	if err != nil {
		fmt.Println("recommendation list unmarshal body response: ", err)
	}

	return r, nil
}



func (r ResourceUsageVirtualMachine) Print() {
	fmt.Printf("usage information:::::%v\n", r)
}




