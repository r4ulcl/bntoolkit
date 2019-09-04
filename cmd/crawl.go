package cmd

import (
	"fmt"
	"log"

	"github.com/RaulCalvoLaorden/bntoolkit/dht"
	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/spf13/cobra"
)

var threads int

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl the BitTorrent Network to find hashes",
	Long: `Crawl the BitTorrent Network to find hashes and storage it in the DB.
For example:
	bntoolkit crawl -t 500`,
	Run: func(cmd *cobra.Command, args []string) {
		ScrapeCmd()
	},
}

//ScrapeCmd crawl the DHT
func ScrapeCmd() {
	fmt.Println("crawl called")

	db, err := utils.ConnectDb(cfgFile, debug, verbose)
	if err != nil {
		log.Fatal("Database error: ", err)
	}
	defer db.Close()
	for {
		dht.Crawler(db, debug, verbose, threads)
	}
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	crawlCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 500, "threads, over 500 you need to change the max files to connections (ulimit -n XXXX)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
