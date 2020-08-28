package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/azure/costmanagement"
	"github.com/Go/azuremonitor/common/filesystem"
	"github.com/Go/azuremonitor/common/terminal"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	r, err := setResourceGroupCostCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupCostCommand() (*cobra.Command, error) {

	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.ResourceGroupCost.DescriptionLine1,
		configuration.ResourceGroupCost.DescriptionLine2,
		configuration.ResourceGroupCost.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.ResourceGroupCost.Command,
		Short: configuration.ResourceGroupCost.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		costmanagement.StartDate = startDate
		costmanagement.EndDate = endDate
		costmanagement.IgnoreZeroCost = ignoreZeroCost
		rgc := costmanagement.ResourceGroupCost{}
		rgc.ExecuteRequest(&rgc)
		rgc.Print()
		if saveCsv {
			filesystem.RemoveFile(csvRgcReportName)
			rgc.WriteCSV(csvRgcReportName)
			fmt.Printf("Done. report was generated - %s\n", csvRgcReportName)
		}

		return nil
	}
	return cmd, nil
}
