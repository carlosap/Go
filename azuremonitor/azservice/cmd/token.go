package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/azure/oauth2"
	"github.com/Go/azuremonitor/common/terminal"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	at, err := setAccessTokenCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(at)
}

func setAccessTokenCommand() (*cobra.Command, error) {

	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.AccessToken.DescriptionLine1,
		configuration.AccessToken.DescriptionLine2,
		configuration.AccessToken.DescriptionLine3)

	cmd := &cobra.Command{
		Use:   configuration.AccessToken.Command,
		Short: configuration.AccessToken.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		at := oauth2.AccessToken{}
		at.ExecuteRequest(&at)
		at.Print()
		return nil
	}
	return cmd, nil
}
