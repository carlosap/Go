package cmd

import (
	"fmt"
	"github.com/Go/azuremonitor/azure/advisor"
	"github.com/Go/azuremonitor/common/terminal"
	c "github.com/Go/azuremonitor/config"
	"github.com/spf13/cobra"
	"os"
)


func init() {

	r, err := setRecommendationListCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(r)
}

func setRecommendationListCommand() (*cobra.Command, error) {

	configuration, _ = c.GetCmdConfig()
	description := fmt.Sprintf("%s\n%s\n%s",
		configuration.RecommendationList.DescriptionLine1,
		configuration.RecommendationList.DescriptionLine2,
		configuration.RecommendationList.DescriptionLine3)

	fmt.Println(description)
	cmd := &cobra.Command{
		Use:   configuration.RecommendationList.Command,
		Short: configuration.RecommendationList.CommandComments,
		Long:  description}

	cmd.RunE = func(*cobra.Command, []string) error {
		terminal.Clear()
		rlist := advisor.RecommendationList{}
		rlist.ExecuteRequest(&rlist)
		rlist.Print()
		return nil
	}
	return cmd, nil
}
