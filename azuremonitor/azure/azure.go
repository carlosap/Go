package azure

type Resource struct {
	SubscriptionID string `json:"subscription_id"`
	ResourceGroup string `json:"resource_group"`
	ResourceID    string `json:"resource_id"`
	Service       string `json:"service"`
	ServiceType   string `json:"service_type"`
	Location      string `json:"location"`
	ChargeType    string `json:"charge_type"`
	Meter         string `json:"meter"`
	Cost          string `json:"cost"`
}

type Resources []Resource



//TODO:::this may require some location
var QueryUrl = "https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourcegroups/defaultresourcegroup-eus/providers/microsoft.operationalinsights/workspaces/" +
	"defaultworkspace-{{subscriptionid}}-eus/query?api-version=2017-10-01"

var VmUsagePayload = "{\"query\": " +
	"\"let startDateTime = datetime('{{startdate}}T08:00:00.000Z');" +
	"let endDateTime = datetime('{{enddate}}T16:00:00.000Z');" +
	"let trendBinSize = 8h;let maxListSize = 1000;" +
	"let cpuMemory = materialize(InsightsMetrics| where TimeGenerated between (startDateTime .. endDateTime)| " +
	"where _ResourceId =~ '/subscriptions/{{subscriptionid}}/resourcegroups/{{resourcegroup}}/providers/microsoft.compute/" +
	"virtualmachines/{{resourceid}}'| " +
	"where Origin == 'vm.azm.ms'| where (Namespace == 'Processor' and Name == 'UtilizationPercentage') or (Namespace == 'Memory' and Name == 'AvailableMB')| " +
	"project TimeGenerated, Name, Namespace, Val);" +
	"let networkDisk = materialize(InsightsMetrics| " +
	"where TimeGenerated between (startDateTime .. endDateTime)| " +
	"where _ResourceId =~ '/subscriptions/" +
	"{{subscriptionid}}/resourcegroups/" +
	"{{resourcegroup}}/providers/microsoft.compute/" +
	"virtualmachines/" +
	"{{resourceid}}'| " +
	"where Origin == 'vm.azm.ms'| " +
	"where (Namespace == 'Network' and Name in ('WriteBytesPerSecond', 'ReadBytesPerSecond'))    " +
	"or (Namespace == 'LogicalDisk' and Name in ('TransfersPerSecond', 'BytesPerSecond', 'TransferLatencyMs'))| " +
	"extend ComputerId = iff(isempty(_ResourceId), Computer, _ResourceId)| " +
	"summarize Val = sum(Val) by bin(TimeGenerated, 1m), " +
	"ComputerId, Name, Namespace| project TimeGenerated, Name, Namespace, Val);" +
	"let rawDataCached = cpuMemory| union networkDisk| " +
	"extend Val = iif(Name in ('WriteLatencyMs', 'ReadLatencyMs', 'TransferLatencyMs'), Val/1000.0, Val)| " +
	"project TimeGenerated,cName = case(Namespace == 'Processor' and Name == 'UtilizationPercentage', '% Processor Time'," +
	"Namespace == 'Memory' and Name == 'AvailableMB', 'Available MBytes',        " +
	"Namespace == 'LogicalDisk' and Name == 'TransfersPerSecond', 'Disk Transfers/sec',        " +
	"Namespace == 'LogicalDisk' and Name == 'BytesPerSecond', 'Disk Bytes/sec',        " +
	"Namespace == 'LogicalDisk' and Name == 'TransferLatencyMs', 'Avg. Disk sec/Transfer',        " +
	"Namespace == 'Network' and Name == 'WriteBytesPerSecond', 'Bytes Sent/sec',        " +
	"Namespace == 'Network' and Name == 'ReadBytesPerSecond', 'Bytes Received/sec',        " +
	"Name),cValue = case(Val < 0, real(0),Val);rawDataCached| summarize min(cValue),    " +
	"avg(cValue),max(cValue),percentiles(cValue, 5, 10, 50, 90, 95) by bin(TimeGenerated, trendBinSize), " +
	"cName| sort by TimeGenerated asc| summarize makelist(TimeGenerated, maxListSize),    makelist(min_cValue, maxListSize)," +
	"makelist(avg_cValue, maxListSize),makelist(max_cValue, maxListSize),makelist(percentile_cValue_5, maxListSize),    " +
	"makelist(percentile_cValue_10, maxListSize),makelist(percentile_cValue_50, maxListSize)," +
	"makelist(percentile_cValue_90, maxListSize),makelist(percentile_cValue_95, maxListSize) by cName| " +
	"join(rawDataCached    | summarize min(cValue), avg(cValue), max(cValue), " +
	"percentiles(cValue, 5, 10, 50, 90, 95) by cName)on cName\"," +
	"\"timespan\": \"{{startdate}}T08:00:00.000Z/{{enddate}}T16:00:00.000Z\"}"

var LocationNames = "location == 'eastus','East US'," +
	"location == 'eastus2','East US 2'," +
	"location == 'southcentralus','South Central US'," +
	"location == 'westus2','West US 2'," +
	"location == 'australiaeast','Australia East'," +
	"location == 'southeastasia','Southeast Asia'," +
	"location == 'northeurope','North Europe'," +
	"location == 'uksouth','UK South'," +
	"location == 'westeurope','West Europe'," +
	"location == 'centralus','Central US'," +
	"location == 'northcentralus','North Central US'," +
	"location == 'westus','West US'," +
	"location == 'southafricanorth','South Africa North'," +
	"location == 'centralindia','Central India'," +
	"location == 'eastasia','East Asia'," +
	"location == 'japaneast','Japan East'," +
	"location == 'koreacentral','Korea Central'," +
	"location == 'canadacentral','Canada Central'," +
	"location == 'francecentral','France Central'," +
	"location == 'germanywestcentral','Germany West Central'," +
	"location == 'norwayeast','Norway East'," +
	"location == 'switzerlandnorth','Switzerland North'," +
	"location == 'uaenorth','UAE North'," +
	"location == 'brazilsouth','Brazil South'," +
	"location == 'centralusstage','Central US (Stage)'," +
	"location == 'eastusstage','East US (Stage)'," +
	"location == 'eastus2stage','East US 2 (Stage)'," +
	"location == 'northcentralusstage','North Central US (Stage)'," +
	"location == 'southcentralusstage','South Central US (Stage)'," +
	"location == 'westusstage','West US (Stage)'," +
	"location == 'westus2stage','West US 2 (Stage)'," +
	"location == 'asia','Asia'," +
	"location == 'asiapacific','Asia Pacific'," +
	"location == 'australia','Australia'," +
	"location == 'brazil','Brazil'," +
	"location == 'canada','Canada'," +
	"location == 'europe','Europe'," +
	"location == 'global','Global'," +
	"location == 'india','India'," +
	"location == 'japan','Japan'," +
	"location == 'uk','United Kingdom'," +
	"location == 'unitedstates','United States'," +
	"location == 'eastasiastage','East Asia (Stage)'," +
	"location == 'southeastasiastage','Southeast Asia (Stage)'," +
	"location == 'westcentralus','West Central US'," +
	"location == 'southafricawest','South Africa West'," +
	"location == 'australiacentral','Australia Central'," +
	"location == 'australiacentral2','Australia Central 2'," +
	"location == 'australiasoutheast','Australia Southeast'," +
	"location == 'japanwest','Japan West'," +
	"location == 'koreasouth','Korea South'," +
	"location == 'southindia','South India'," +
	"location == 'westindia','West India'," +
	"location == 'canadaeast','Canada East'," +
	"location == 'francesouth','France South'," +
	"location == 'germanynorth','Germany North'," +
	"location == 'norwaywest','Norway West'," +
	"location == 'switzerlandwest','Switzerland West'," +
	"location == 'ukwest','UK West'," +
	"location == 'uaecentral','UAE Central',"
