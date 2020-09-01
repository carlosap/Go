package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/azure/subscription"
	"github.com/Go/azuremonitor/common/terminal"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)

func init() {

	r, err := setResourcesCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourcesCommand() (*cobra.Command, error) {

	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.Resources.DescriptionLine1,
		configuration.Resources.DescriptionLine2,
		configuration.Resources.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.Resources.Command,
		Short: configuration.Resources.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		resource := subscription.ResourceSubscription{}
		resource.ExecuteRequest(&resource)
		resource.Print()
		return nil
	}
	return cmd, nil
}
