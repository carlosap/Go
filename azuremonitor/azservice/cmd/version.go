package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of elysium localization",
	Long:  `All software has versions`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf("args %v\n", args)
		}

		//clearTerminal()
		fmt.Println("Elysium Localization v0.2")
	},
}
