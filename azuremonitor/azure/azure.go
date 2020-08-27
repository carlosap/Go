package azure

type Resource struct {
	ResourceGroup string `json:"resourcegroup"`
	ResourceID    string `json:"resourceid"`
	Service       string `json:"service"`
	ServiceType   string `json:"serviceType"`
	Location      string `json:"location"`
	ChargeType    string `json:"chargetype"`
	Meter         string `json:"meter"`
	Cost          string `json:"cost"`
}

type Resources []Resource



var LocationNames = "location == 'eastus','East US'," +
	"location == 'eastus2','East US 2',"+
	"location == 'southcentralus','South Central US'," +
	"location == 'westus2','West US 2'," +
	"location == 'australiaeast','Australia East',"+
	"location == 'southeastasia','Southeast Asia'," +
	"location == 'northeurope','North Europe'," +
	"location == 'uksouth','UK South',"+
	"location == 'westeurope','West Europe'," +
	"location == 'centralus','Central US'," +
	"location == 'northcentralus','North Central US',"+
	"location == 'westus','West US'," +
	"location == 'southafricanorth','South Africa North'," +
	"location == 'centralindia','Central India',"+
	"location == 'eastasia','East Asia'," +
	"location == 'japaneast','Japan East'," +
	"location == 'koreacentral','Korea Central',"+
	"location == 'canadacentral','Canada Central'," +
	"location == 'francecentral','France Central',"+
	"location == 'germanywestcentral','Germany West Central'," +
	"location == 'norwayeast','Norway East',"+
	"location == 'switzerlandnorth','Switzerland North'," +
	"location == 'uaenorth','UAE North',"+
	"location == 'brazilsouth','Brazil South'," +
	"location == 'centralusstage','Central US (Stage)',"+
	"location == 'eastusstage','East US (Stage)'," +
	"location == 'eastus2stage','East US 2 (Stage)',"+
	"location == 'northcentralusstage','North Central US (Stage)',"+
	"location == 'southcentralusstage','South Central US (Stage)',"+
	"location == 'westusstage','West US (Stage)'," +
	"location == 'westus2stage','West US 2 (Stage)',"+
	"location == 'asia','Asia'," +
	"location == 'asiapacific','Asia Pacific'," +
	"location == 'australia','Australia',"+
	"location == 'brazil','Brazil'," +
	"location == 'canada','Canada'," +
	"location == 'europe','Europe',"+
	"location == 'global','Global'," +
	"location == 'india','India'," +
	"location == 'japan','Japan',"+
	"location == 'uk','United Kingdom'," +
	"location == 'unitedstates','United States',"+
	"location == 'eastasiastage','East Asia (Stage)'," +
	"location == 'southeastasiastage','Southeast Asia (Stage)',"+
	"location == 'westcentralus','West Central US',"+
	"location == 'southafricawest','South Africa West',"+
	"location == 'australiacentral','Australia Central'," +
	"location == 'australiacentral2','Australia Central 2',"+
	"location == 'australiasoutheast','Australia Southeast'," +
	"location == 'japanwest','Japan West',"+
	"location == 'koreasouth','Korea South'," +
	"location == 'southindia','South India',"+
	"location == 'westindia','West India'," +
	"location == 'canadaeast','Canada East',"+
	"location == 'francesouth','France South'," +
	"location == 'germanynorth','Germany North',"+
	"location == 'norwaywest','Norway West'," +
	"location == 'switzerlandwest','Switzerland West',"+
	"location == 'ukwest','UK West'," +
	"location == 'uaecentral','UAE Central',"
