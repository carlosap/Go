package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/robfig/cron.v3"
	"os"
	"time"
)


var ctr int = 0
var startTime time.Time

func init() {

	cron, err := setScheduleCronCommand()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rootCmd.AddCommand(cron)
}

func setScheduleCronCommand() (*cobra.Command, error) {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "",
		Long: `
`,
	}

	cmd.RunE = func(*cobra.Command, []string) error {
		fmt.Println("....wait we are learning ip information")
		startTime = time.Now()

		c := cron.New()
		c.AddFunc("@every 0h0m15s", func() {
			ctr++
			delta := time.Now().Sub(startTime)
			clearTerminal()
			fmt.Printf("[%d] auto discovery running since: [%vhr]:[%vmin]:[%vsec]\n", ctr, delta.Hours(), delta.Minutes(), delta.Seconds())

		})

		c.Start()
		//startServer()
		return nil
	}

	return cmd, nil
}
