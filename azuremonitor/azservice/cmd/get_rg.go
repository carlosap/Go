package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/azure/batch"
	"github.com/Go/azuremonitor/common/terminal"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)

func init() {

	r, err := setResourceGroupCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setResourceGroupCommand() (*cobra.Command, error) {
	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.ResourceGroups.DescriptionLine1,
		configuration.ResourceGroups.DescriptionLine2,
		configuration.ResourceGroups.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.ResourceGroups.Command,
		Short: configuration.ResourceGroups.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		rgl := batch.ResourceGroupList{}
		rgl.ExecuteRequest(&rgl)
		rgl.Print()
		return nil
	}
	return cmd, nil
}
