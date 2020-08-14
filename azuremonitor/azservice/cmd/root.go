package cmd

import (
	"fmt"

	"github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)

var cmdConfig config.CmdConfig

func init() {
	cmdConfig, err := config.GetCmdConfig("config.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("startign config file %v", cmdConfig)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "localization services",
	Short: "Azure Localization Services",
	Long:  ``,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
