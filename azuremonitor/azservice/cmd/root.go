package cmd

import "C"
import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "localization services",
	Short: "Azure Localization Services",
	Long:  ``,
}

func init() {

	now := time.Now()
	month := now.AddDate(0, 0, -29)
	rootCmd.PersistentFlags().StringVar(&startDate, "from", month.Format(layoutISO), "start date of report (i.e. YYYY-MM-DD)")
	rootCmd.PersistentFlags().StringVar(&endDate, "to", now.Format(layoutISO), "end date of report (i.e. YYYY-MM-DD)")
	rootCmd.PersistentFlags().BoolVar(&saveDb, "db", false, "[=true]saves records to Postgres db")
	rootCmd.PersistentFlags().BoolVar(&saveCsv, "csv", false, "[=true]saves records into a csv output file")
	rootCmd.PersistentFlags().BoolVar(&ignoreZeroCost, "izcost", false, "[=true] ignores resources with zero cost")

}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
