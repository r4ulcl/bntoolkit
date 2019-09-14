package cmd

import (
	"fmt"

	"github.com/RaulCalvoLaorden/bntoolkit/dht"
	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/spf13/cobra"
)

//var projectName string
var crawlDaemon bool

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Start the daemon to monitor the files in the monitor table and notify alerts.",
	Long: `Start the daemon to monitor the files in the monitor table, notify alerts and optionally crape DHT
For example:
	bntoolkit daemon -s`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("daemon called")

		//scrape
		if crawlDaemon {
			go ScrapeCmd()
			waitGroup.Add(1)
		}

		//alert monitor
		go utils.MonitorAlert(cfgFile, debug, verbose, projectName)
		waitGroup.Add(1)

		//Starting monitor
		go dht.DaemonPeers(cfgFile, debug, verbose, projectName)
		waitGroup.Add(1)

		waitGroup.Wait()
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	daemonCmd.PersistentFlags().StringVarP(&projectName, "projectName", "p", "default", "Monitoring project")
	daemonCmd.PersistentFlags().BoolVarP(&crawlDaemon, "crawl", "s", false, "Crawl DHT")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// daemonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
