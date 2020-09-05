package azure

type Resource struct {
	SubscriptionID         string   `json:"subscription_id"`
	ResourceGroupName      string   `json:"resource_group_name"`
	ResourceID             string   `json:"resource_id"`
	ResourceType           string   `json:"resource_type"`
	ResourceLocation       string   `json:"resource_location"`
	ChargeType             string   `json:"charge_type"`
	ServiceName            string   `json:"service_name"`
	Meter                  string   `json:"meter"`
	MeterCategory          string   `json:"meter_category"`
	MeterSubCategory       string   `json:"meter_subcategory"`
	ServiceFamily          string   `json:"service_family"`
	UnitOfMeasure          string   `json:"unit_of_measure"`
	CostAllocationRuleName string   `json:"cost_allocation_rule_name"`
	Product                string   `json:"product"`
	Frequency              string   `json:"frequency"`
	PricingModel           string   `json:"pricing_model"`
	Tags                   []string `json:"tags"`
	Currency               string   `json:"currency"`
	PreTaxCostUSD          float64  `json:"pre_tax_cost_usd"`
	UsageQuantity          float64  `json:"usage_quantity"`
}

type Resources []Resource

var QueryUrl = "https://management.azure.com/batch?api-version=2015-11-01"

var ActualCostManagementPayload = "{\"type\":\"ActualCost\",\"dataSet\":{\"granularity\":\"None\",\"aggregation\":{\"pretaxCost\":{\"name\":\"PreTaxCostUSD\",\"function\":\"Sum\"},\"usageQuantity\":{\"name\":\"UsageQuantity\",\"function\":\"Sum\"}},\"grouping\":[{\"type\":\"Dimension\",\"name\":\"ResourceId\"},{\"type\":\"Dimension\",\"name\":\"ResourceType\"},{\"type\":\"Dimension\",\"name\":\"ResourceLocation\"},{\"type\":\"Dimension\",\"name\":\"ChargeType\"},{\"type\":\"Dimension\",\"name\":\"ResourceGroupName\"},{\"type\":\"Dimension\",\"name\":\"ServiceName\"},{\"type\":\"Dimension\",\"name\":\"Meter\"},{\"type\":\"Dimension\",\"name\":\"MeterCategory\"},{\"type\":\"Dimension\",\"name\":\"MeterSubCategory\"},{\"type\":\"Dimension\",\"name\":\"ServiceFamily\"},{\"type\":\"Dimension\",\"name\":\"UnitOfMeasure\"},{\"type\":\"Dimension\",\"name\":\"CostAllocationRuleName\"},{\"type\":\"Dimension\",\"name\":\"Product\"},{\"type\":\"Dimension\",\"name\":\"Frequency\"},{\"type\":\"Dimension\",\"name\":\"PricingModel\"}],\"include\":[\"Tags\"]},\"timeframe\":\"Custom\",\"timePeriod\":" +
	"{\"from\":\"{{startdate}}T00:00:00+00:00\"," +
	"\"to\":\"{{enddate}}T23:59:59+00:00\"}}"

var StorageStorageAccountPayload = "{\"requests\":[{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:00:00.000Z/{{enddate}}T20:00:00.000Z&interval=FULL&metricnames=Egress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:00:00.000Z/{{enddate}}T20:00:00.000Z&interval=FULL&metricnames=Ingress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:00:00.000Z/{{enddate}}T20:00:00.000Z&interval=FULL&metricnames=Transactions&aggregation=total&metricNamespace=microsoft.storage%2Fstorageaccounts&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/blobServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=Egress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fblobservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/blobServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=Ingress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fblobservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/blobServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=Transactions&aggregation=total&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fblobservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/blobServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=BlobCount&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fblobservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/fileServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=FileCount&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ffileservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/fileServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=Egress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ffileservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/fileServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=Ingress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ffileservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/fileServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=Transactions&aggregation=total&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ffileservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/queueServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=QueueCount&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fqueueservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/queueServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=Egress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fqueueservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/queueServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=Ingress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fqueueservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/queueServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T20:15:00.000Z/{{enddate}}T20:15:00.000Z&interval=FULL&metricnames=Transactions&aggregation=total&metricNamespace=microsoft.storage%2Fstorageaccounts%2Fqueueservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/tableServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=TableCount&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ftableservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/tableServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=TableEntityCount&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ftableservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/tableServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T08:25:00.000Z/{{enddate}}T20:25:00.000Z&interval=FULL&metricnames=Egress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ftableservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/tableServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T08:25:00.000Z/{{enddate}}T20:25:00.000Z&interval=FULL&metricnames=Ingress&aggregation=average&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ftableservices&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\":\"GET\",\"relativeUrl\":\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Storage/storageAccounts/" +
	"{{resourceid}}/tableServices/default/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T08:25:00.000Z/{{enddate}}T20:25:00.000Z&interval=FULL&metricnames=Transactions&aggregation=total&metricNamespace=microsoft.storage%2Fstorageaccounts%2Ftableservices&validatedimensions=false&api-version=2019-07-01\"}]}"

var LogicAppUsagePayload = "{\"requests\": [{\"httpMethod\": \"GET\",\"relativeUrl\": \"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Logic/workflows/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T16:00:00.000Z/{{enddate}}T16:00:00.000Z&interval=FULL&metricnames=TotalBillableExecutions" +
	"&aggregation=average&metricNamespace=microsoft.logic%2Fworkflows&validatedimensions=false&api-version=2019-07-01\"}, " +
	"{\"httpMethod\": \"GET\",\"relativeUrl\": \"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Logic/workflows/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T16:00:00.000Z/{{enddate}}T16:00:00.000Z&interval=FULL&metricnames=BillableActionExecutions" +
	"&aggregation=total&metricNamespace=microsoft.logic%2Fworkflows&validatedimensions=false&api-version=2019-07-01\"}, " +
	"{\"httpMethod\": \"GET\",\"relativeUrl\": \"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Logic/workflows/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T16:00:00.000Z/{{enddate}}T16:00:00.000Z&interval=FULL&" +
	"metricnames=BillingUsageNativeOperation&aggregation=total&metricNamespace=microsoft.logic%2Fworkflows&validatedimensions=false&api-version=2019-07-01\"},  " +
	"{\"httpMethod\": \"GET\",\"relativeUrl\": \"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Logic/workflows/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T16:00:00.000Z/{{enddate}}T16:00:00.000Z&interval=FULL&metricnames=BillingUsageStandardConnector" +
	"&aggregation=average&metricNamespace=microsoft.logic%2Fworkflows&validatedimensions=false&api-version=2019-07-01\"}, " +
	"{\"httpMethod\": \"GET\",\"relativeUrl\": \"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Logic/workflows/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T16:00:00.000Z/{{enddate}}T16:00:00.000Z&interval=FULL&metricnames=BillingUsageStorageConsumption&" +
	"aggregation=average&metricNamespace=microsoft.logic%2Fworkflows&validatedimensions=false&api-version=2019-07-01\"}]}"

var StorageDiskUsagePayload = "{\"requests\": [{\"httpMethod\": \"GET\",\"url\": " +
	"\"https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?" +
	"timespan={{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=OS%20Disk%20Read%20Bytes%2Fsec&" +
	"aggregation=average&validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\",\"url\": " +
	"\"https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?" +
	"timespan={{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=OS%20Disk%20Write%20Bytes%2Fsec&" +
	"aggregation=average&validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\",\"url\": " +
	"\"https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?" +
	"timespan={{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&" +
	"metricnames=OS%20Disk%20Read%20Operations%2FSec&" +
	"aggregation=average&validatedimensions=false&api-version=2019-07-01\"},{\"httpMethod\": \"GET\",\"url\": " +
	"\"https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?" +
	"timespan={{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=OS%20Disk%20Write%20Operations%2FSec&" +
	"aggregation=average&validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\",\"url\": " +
	"\"https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?" +
	"timespan={{startdate}}T22:00:00.000Z/{{enddate}}T22:00:00.000Z&interval=FULL&metricnames=OS%20Disk%20Queue%20Depth&aggregation=average&" +
	"validatedimensions=false&api-version=2019-07-01\"}]}"

var VmUsagePayload = "{\"requests\": [{\"httpMethod\": \"GET\",\"url\": " +
	"\"https://management.azure.com/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}?api-version=2020-06-01&$expand=instanceView\"}," +
	"{\"httpMethod\": \"GET\",\"relativeUrl\": " +
	"\"/subscriptions/{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T19:45:00.000Z/{{enddate}}T19:45:00.000Z&" +
	"interval=FULL&metricnames=Percentage CPU&aggregation=average&" +
	"metricNamespace=microsoft.compute%2Fvirtualmachines&" +
	"validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\",\"relativeUrl\": " +
	"\"/subscriptions/{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T19:45:00.000Z/{{enddate}}T19:45:00.000Z&" +
	"interval=FULL&metricnames=Disk Read Bytes&" +
	"aggregation=total&metricNamespace=microsoft.compute%2Fvirtualmachines&" +
	"validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\",\"relativeUrl\": " +
	"\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T19:45:00.000Z/{{enddate}}T19:45:00.000Z&" +
	"interval=FULL&metricnames=Disk Write Bytes&aggregation=total" +
	"&metricNamespace=microsoft.compute%2Fvirtualmachines&validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\",\"relativeUrl\": " +
	"\"/subscriptions/" +
	"{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T19:45:00.000Z/{{enddate}}T19:45:00.000Z" +
	"&interval=FULL&metricnames=Network In Total&aggregation=total&" +
	"metricNamespace=microsoft.compute%2Fvirtualmachines&validatedimensions=false&api-version=2019-07-01\"}, {\"httpMethod\": \"GET\"," +
	"\"relativeUrl\": \"/subscriptions/{{subscriptionid}}/resourceGroups/" +
	"{{resourcegroup}}/providers/Microsoft.Compute/virtualMachines/" +
	"{{resourceid}}/providers/microsoft.Insights/metrics?timespan=" +
	"{{startdate}}T19:45:00.000Z/{{enddate}}T19:45:00.000Z&" +
	"interval=FULL&metricnames=Network Out Total" +
	"&aggregation=total&metricNamespace=microsoft.compute%2Fvirtualmachines&validatedimensions=false&api-version=2019-07-01\"}]}"

//var VmUsagePayload = "{\"query\": " +
//	"\"let startDateTime = datetime('{{startdate}}T08:00:00.000Z');" +
//	"let endDateTime = datetime('{{enddate}}T16:00:00.000Z');" +
//	"let trendBinSize = 8h;let maxListSize = 1000;" +
//	"let cpuMemory = materialize(InsightsMetrics| where TimeGenerated between (startDateTime .. endDateTime)| " +
//	"where _ResourceId =~ '/subscriptions/{{subscriptionid}}/resourcegroups/{{resourcegroup}}/providers/microsoft.compute/" +
//	"virtualmachines/{{resourceid}}'| " +
//	"where Origin == 'vm.azm.ms'| where (Namespace == 'Processor' and Name == 'UtilizationPercentage') or (Namespace == 'Memory' and Name == 'AvailableMB')| " +
//	"project TimeGenerated, Name, Namespace, Val);" +
//	"let networkDisk = materialize(InsightsMetrics| " +
//	"where TimeGenerated between (startDateTime .. endDateTime)| " +
//	"where _ResourceId =~ '/subscriptions/" +
//	"{{subscriptionid}}/resourcegroups/" +
//	"{{resourcegroup}}/providers/microsoft.compute/" +
//	"virtualmachines/" +
//	"{{resourceid}}'| " +
//	"where Origin == 'vm.azm.ms'| " +
//	"where (Namespace == 'Network' and Name in ('WriteBytesPerSecond', 'ReadBytesPerSecond'))    " +
//	"or (Namespace == 'LogicalDisk' and Name in ('TransfersPerSecond', 'BytesPerSecond', 'TransferLatencyMs'))| " +
//	"extend ComputerId = iff(isempty(_ResourceId), Computer, _ResourceId)| " +
//	"summarize Val = sum(Val) by bin(TimeGenerated, 1m), " +
//	"ComputerId, Name, Namespace| project TimeGenerated, Name, Namespace, Val);" +
//	"let rawDataCached = cpuMemory| union networkDisk| " +
//	"extend Val = iif(Name in ('WriteLatencyMs', 'ReadLatencyMs', 'TransferLatencyMs'), Val/1000.0, Val)| " +
//	"project TimeGenerated,cName = case(Namespace == 'Processor' and Name == 'UtilizationPercentage', '% Processor Time'," +
//	"Namespace == 'Memory' and Name == 'AvailableMB', 'Available MBytes'," +
//	"Namespace == 'LogicalDisk' and Name == 'TransfersPerSecond', 'Disk Transfers/sec'," +
//	"Namespace == 'LogicalDisk' and Name == 'BytesPerSecond', 'Disk Bytes/sec'," +
//	"Namespace == 'LogicalDisk' and Name == 'TransferLatencyMs', 'Avg. Disk sec/Transfer'," +
//	"Namespace == 'Network' and Name == 'WriteBytesPerSecond', 'Bytes Sent/sec'," +
//	"Namespace == 'Network' and Name == 'ReadBytesPerSecond', 'Bytes Received/sec'," +
//	"Name),cValue = case(Val < 0, real(0),Val);rawDataCached| summarize min(cValue)," +
//	"avg(cValue),max(cValue),percentiles(cValue, 5, 10, 50, 90, 95) by bin(TimeGenerated, trendBinSize), " +
//	"cName| sort by TimeGenerated asc| summarize makelist(TimeGenerated, maxListSize),    makelist(min_cValue, maxListSize)," +
//	"makelist(avg_cValue, maxListSize),makelist(max_cValue, maxListSize),makelist(percentile_cValue_5, maxListSize),    " +
//	"makelist(percentile_cValue_10, maxListSize),makelist(percentile_cValue_50, maxListSize)," +
//	"makelist(percentile_cValue_90, maxListSize),makelist(percentile_cValue_95, maxListSize) by cName| " +
//	"join(rawDataCached    | summarize min(cValue), avg(cValue), max(cValue), " +
//	"percentiles(cValue, 5, 10, 50, 90, 95) by cName)on cName\"," +
//	"\"timespan\": \"{{startdate}}T08:00:00.000Z/{{enddate}}T16:00:00.000Z\"}"

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

//===================================prefix region and workspace=============================
//If you are wondering why your Azure subscription has a resource group
//called DefaultResourceGroup-XXX (the XXX is related to your region) and
//within that same resource group you have a DefaultWorkspace-<SubscriptionID>-XXX,
//there is a logical explanation, and it is associated with Azure Security Center.

//====================================Action Items===========================================
// - Need to determine geo per vm request to optian usage. there may be an easier way on doing that.
