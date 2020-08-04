package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"fmt"
	"os"
	"strings"
)


type ResourceGroups struct {
	Responses []struct {
		Content struct {
			TotalRecords int `json:"totalRecords"`
			Count        int `json:"count"`
			Data         struct {
				Rows [][]interface{} `json:"rows"`
			} `json:"data"`
		} `json:"content"`
	} `json:"responses"`
}

type ResourceGroupList []string

func init() {

	r, err := setResourceGroupCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupCommand() (*cobra.Command, error) {
	cl := Client{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	description := fmt.Sprintf("%s\n%s\n%s",
		cl.AppConfig.ResourceGroups.DescriptionLine1,
		cl.AppConfig.ResourceGroups.DescriptionLine2,
		cl.AppConfig.ResourceGroups.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   cl.AppConfig.ResourceGroups.Command,
		Short: cl.AppConfig.ResourceGroups.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		r := ResourceGroupList{}
		r, err := r.getResourceGroups()
		if err != nil {
			return err
		}

		clearTerminal()
		r.Print()
		return nil
	}
	return cmd, nil
}

func (r ResourceGroupList) getResourceGroups() (ResourceGroupList, error) {
	var at = &AccessToken{}
	cl := Client{}
	rg := ResourceGroups{}
	err := cl.New()
	if err != nil {
		return nil, err
	}

	at, err = at.getAccessToken()
	if err != nil {
		return nil, err
	}

	token := fmt.Sprintf("Bearer %s", at.AccessToken)
	payload := strings.NewReader(fmt.Sprintf("{\"requests\": [{\"content\": {\"subscriptions\": [\"%s\"]," +
		"\"query\": \"(resourcecontainers|where type in~ ('microsoft.resources/subscriptions/resourcegroups'))" +
		"|where type =~ 'microsoft.resources/subscriptions/resourcegroups'\\r\\n| " +
		"extend status = case(\\r\\n    (properties.provisioningState =~ 'accepted'), " +
		"'Creating',\\r\\n    (properties.provisioningState =~ 'deleted'), " +
		"'Deleted',\\r\\n    (properties.provisioningState =~ 'deleting'), " +
		"'Deleting',\\r\\n    (properties.provisioningState =~ 'failed'), " +
		"'Failed',\\r\\n    (properties.provisioningState =~ 'movingresources'), " +
		"'Moving Resources',\\r\\n    properties.provisioningState)\\r\\n| " +
		"project id, name, type, location, subscriptionId, resourceGroup, kind, tags, status\\r\\n|" +
		"extend subscriptionDisplayName=case(subscriptionId == '%s','SpartanAppSolutions',subscriptionId)|" +
		"extend locationDisplayName=case(location == 'eastus','East US',location == 'eastus2','East US 2'," +
		"location == 'southcentralus','South Central US',location == 'westus2','West US 2',location == 'australiaeast','Australia East'," +
		"location == 'southeastasia','Southeast Asia',location == 'northeurope','North Europe',location == 'uksouth','UK South'," +
		"location == 'westeurope','West Europe',location == 'centralus','Central US',location == 'northcentralus','North Central US'," +
		"location == 'westus','West US',location == 'southafricanorth','South Africa North',location == 'centralindia','Central India'," +
		"location == 'eastasia','East Asia',location == 'japaneast','Japan East',location == 'koreacentral','Korea Central'," +
		"location == 'canadacentral','Canada Central',location == 'francecentral','France Central'," +
		"location == 'germanywestcentral','Germany West Central',location == 'norwayeast','Norway East'," +
		"location == 'switzerlandnorth','Switzerland North',location == 'uaenorth','UAE North'," +
		"location == 'brazilsouth','Brazil South',location == 'centralusstage','Central US (Stage)'," +
		"location == 'eastusstage','East US (Stage)',location == 'eastus2stage','East US 2 (Stage)'," +
		"location == 'northcentralusstage','North Central US (Stage)'," +
		"location == 'southcentralusstage','South Central US (Stage)'," +
		"location == 'westusstage','West US (Stage)',location == 'westus2stage','West US 2 (Stage)'," +
		"location == 'asia','Asia',location == 'asiapacific','Asia Pacific',location == 'australia','Australia'," +
		"location == 'brazil','Brazil',location == 'canada','Canada',location == 'europe','Europe'," +
		"location == 'global','Global',location == 'india','India',location == 'japan','Japan'," +
		"location == 'uk','United Kingdom',location == 'unitedstates','United States'," +
		"location == 'eastasiastage','East Asia (Stage)',location == 'southeastasiastage','Southeast Asia (Stage)'," +
		"location == 'westcentralus','West Central US'," +
		"location == 'southafricawest','South Africa West'," +
		"location == 'australiacentral','Australia Central',location == 'australiacentral2','Australia Central 2'," +
		"location == 'australiasoutheast','Australia Southeast',location == 'japanwest','Japan West'," +
		"location == 'koreasouth','Korea South',location == 'southindia','South India'," +
		"location == 'westindia','West India',location == 'canadaeast','Canada East'," +
		"location == 'francesouth','France South',location == 'germanynorth','Germany North'," +
		"location == 'norwaywest','Norway West',location == 'switzerlandwest','Switzerland West'," +
		"location == 'ukwest','UK West',location == 'uaecentral','UAE Central',location)|where (type !~ ('microsoft.confluent/organizations'))|" +
		"where (type !~ ('microsoft.securitydetonation/chambers'))|where (type !~ ('microsoft.intelligentitdigitaltwin/digitaltwins'))|" +
		"where (type !~ ('microsoft.connectedcache/cachenodes'))|where (type !~ ('microsoft.serviceshub/connectors'))|" +
		"where not((type =~ ('microsoft.sql/servers/databases')) and ((kind in~ ('system','v2.0,system','v12.0,system','v12.0,user,datawarehouse,gen2,analytics'))))|" +
		"where not((type =~ ('microsoft.sql/servers')) and ((kind =~ ('v12.0,analytics'))))|" +
		"project name,subscriptionDisplayName,locationDisplayName,id,type,kind,location,subscriptionId,resourceGroup,tags|" +
		"sort by tolower(tostring(name)) asc\",\"options\": {\"$top\": 100,\"$skip\": 0,\"$skipToken\": \"\"}}," +
		"\"httpMethod\": \"POST\",\"name\": \"34cc625b-b20a-423a-9563-33faf337b033\"," +
		"\"requestHeaderDetails\": {\"commandName\": \"HubsExtension.BrowseResourceGroups.microsoft.resources/subscriptions/resourcegroups.InitialLoad\"}," +
		"\"url\": \"https://management.azure.com/providers/Microsoft.ResourceGraph/resources?api-version=2018-09-01-preview\"}]}",
		cl.AppConfig.AccessToken.SubscriptionID,
		cl.AppConfig.AccessToken.SubscriptionID,
	))

	client := &http.Client {}
	req, _ := http.NewRequest("POST", cl.AppConfig.ResourceGroups.URL, payload)
	req.Header.Add("Authorization", token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &rg)
	if err != nil {
		fmt.Println("recommendation list unmarshal body response: ", err)
	}

	for i := 0; i < len(rg.Responses); i++ {
		rp := rg.Responses[i]
		if rp.Content.TotalRecords > 0 {
			for x := 0; x < rp.Content.TotalRecords; x++ {
				row := rp.Content.Data.Rows[x]
				if len(row) > 0 {
					//casting interface to string
					str := fmt.Sprintf("%v", row[0])
					r = append(r, str)
				}
			}
		}
	}

	return r, nil
}


func (r ResourceGroupList) Print() {
	fmt.Println("Resource Groups:")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------")
	for i := 0; i < len(r); i++ {
		fmt.Println(r[i])
	}
}




