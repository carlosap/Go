package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/common/filesystem"
	"github.com/spf13/cobra"
)

func init() {

	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean-cache",
	Short: "clears cache resources",
	Long:  `clear requests and local cache sources`,
	Run: func(cmd *cobra.Command, args []string) {
		dirName := "cache"
		dbName := "spartan"
		isSuccess := filesystem.RemoveFile(dbName)
		if !isSuccess {
			fmt.Println("failed to remove cache file")
		} else {
			fmt.Println("...done.")
		}

		isSuccess = filesystem.IsDirectoryExist(dirName)
		if isSuccess {
			filesystem.RemoveDirectory(dirName)
		}

	},
}
