package cmd

import (
	"fmt"
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
		clearTerminal()
		developer, _ := cmd.Flags().GetString("developer")
		if developer != "" {
			fmt.Printf("Developer: %s\n", developer)
		}
		fmt.Println("Azmonitor ", version)



		//requests := Requests{
		//	{"accesstoken", "", url, "POST", strPayload, header, item, nil},
		//	{"google", "", "https://www.google.com", "GET", "", nil, nil, nil},
		//	{"msn", "", "https://www.msn.com", "GET", "", nil, nil, nil},
		//}
		//
	    //errors := requests.Execute()
		//if len(errors) > 0 {
		//	fmt.Fprintf(os.Stderr, "\n%d errors occurred:\n", len(errors))
		//	for _, err := range errors {
		//		fmt.Fprintf(os.Stderr, "--> %s\n", err)
		//	}
		//}
		//
		//for _, r := range requests {
		//	body := r.GetResponse()
		//	fmt.Printf("The body of %s - %s\n", r.Name, body)
		//}



		//fmt.Printf("The body of %s - %s\n", request.Name, body)
		//fmt.Printf("The value of %s - %v\n", request.Name,act.AccessToken )



	},
}
