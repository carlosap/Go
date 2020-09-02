package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/azure/costmanagement"
	"github.com/Go/azuremonitor/common/terminal"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	r, err := setResourceGroupUsageCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupUsageCommand() (*cobra.Command, error) {

	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.ResourceGroupUsage.DescriptionLine1,
		configuration.ResourceGroupUsage.DescriptionLine2,
		configuration.ResourceGroupUsage.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.ResourceGroupUsage.Command,
		Short: configuration.ResourceGroupUsage.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		costmanagement.StartDate = startDate
		costmanagement.EndDate = endDate
		costmanagement.IgnoreZeroCost = ignoreZeroCost
		costmanagement.SaveCsv = saveCsv
		usage := costmanagement.ResourceGroupUsage{}
		usage.RunAll()
		return nil
	}
	return cmd, nil
}
