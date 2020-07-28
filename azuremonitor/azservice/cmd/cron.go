package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/robfig/cron.v3"
	"os"
	"time"
)

const (
	// See http://golang.org/pkg/time/#Parse
	timeFormat = "2006-01-02 15:04 MST"
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
		Short: "Runs endless sequenses of API discoveries (ip, weather, forecast, etc..)",
		Long: `Scheduler that runs every hour to refresh latest API responses.
`,
	}

	cmd.RunE = func(*cobra.Command, []string) error {
		fmt.Println("....wait we are learning ip information")
		startTime = time.Now()

		c := cron.New()
		IpInfo, err := IpInfo.getIpInfo(false)
		if err != nil {
			fmt.Printf("Warning: ip info failed. Check internet connection %v\n", err)
		}

		IpInfo.Print()

		c.AddFunc("@every 0h0m15s", func() {
			ctr++
			delta := time.Now().Sub(startTime)
			clearTerminal()
			fmt.Printf("[%d] auto discovery running since: [%vhr]:[%vmin]:[%vsec]\n", ctr, delta.Hours(), delta.Minutes(), delta.Seconds())

			ip := sendIpInfo(IpInfo)
			if len(ip.IP) > 0 {
				sendWeather(ip)
				sendForecast(ip)
			}
		})

		c.Start()
		startServer()
		return nil
	}

	return cmd, nil
}

func sendIpInfo(IpInfo *IpapiResponse) *IpapiResponse {
	ip, err := IpInfo.getIpInfo(true)
	if err != nil {
		fmt.Printf("Warning: ip info failed. Check internet connection %v\n", err)
	}
	ip.Print()
	err = ip.clientIpInfoResponse()
	if err != nil {
		fmt.Printf("ip socket response:::%v\n", err)
	}
	return ip
}

func sendWeather(ip *IpapiResponse) {
	w := &Weather{}
	w, err := w.getWeather(ip)
	if err != nil {
		fmt.Printf("Warning: ip info failed. retrieve weather %v\n", err)
	}

	w.Print()
	err = w.clientWeatherResponse()
	if err != nil {
		fmt.Printf("Warnning: weather response:::%v\n", err)
	}
}

func sendForecast(ip *IpapiResponse) {
	f := &ForeCast{}
	f, err := f.getForeCast(ip)
	if err != nil {
		fmt.Printf("forecast failed. retrieve forecast %v\n", err)
	}

	f.Print()
	err = f.clientForecastResponse()
	if err != nil {
		fmt.Printf("forecast response:::%v\n", err)
	}
}
