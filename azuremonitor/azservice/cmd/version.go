package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/common/terminal"
	"github.com/spf13/cobra"
)

func init() {
	//cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&developer, "developer", "Carlos Perez", "Developer name.")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of elysium localization",
	Long:  `All software has versions`,
	Run: func(cmd *cobra.Command, args []string) {
		terminal.Clear()
		developer, _ := cmd.Flags().GetString("developer")
		if developer != "" {
			fmt.Printf("Developer: %s\n", developer)
		}
		fmt.Println("Azmonitor ", version)

	},
}
