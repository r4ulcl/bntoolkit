package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RaulCalvoLaorden/bntoolkit/dht"
	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/spf13/cobra"
)

var file string
var noadd bool
var typeF int
var timeout int

var waitGroup sync.WaitGroup

//FileInfo : Path plus lenght
type FileInfo struct {
	Path   []string
	Length int
}

var hashesChan chan string
var inicio = 10

var (
	builtinAnnounceList = [][]string{
		{"udp://tracker.openbittorrent.com:80"},
		{"udp://tracker.publicbt.com:80"},
		{"udp://tracker.istole.it:6969"},
	}
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find <path to file>",
	Short: "Find the file in Bittorrent network",
	Long: `Find the file in Bittorrent network using the DHT, a trackers list and the local database. In this command the hashes can be: Possibles, Valid or Downloaded. The first are the ones that could exist because they are valid, the second are the ones that have been found in BitTorrent and the third is that it has peers and can be downloaded.
For example: 
	bntoolkit find ubuntu-18.04.1-desktop-amd64.iso.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("find called", len(args))

		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}

		file = strings.Join(args, " ")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		utils.InsertProject(db, debug, verbose, projectName) //if return error exists, so OK

		switch typeF {
		case 0:
			dht.WorkersTorrents(cfgFile, debug, verbose, file, projectName) //create hashes and check in db
			dht.ScrapeTrackers(db, debug, verbose, projectName)             //scrape trackers
			fmt.Println("Possibles (possibles hashes to the file)")
			err := utils.SelectPossiblesWhere("possible", db, debug, verbose, projectName) //get only possibles
			if err != nil {
				log.Println(err)
			}
		case 1:
			dht.WorkersTorrents(cfgFile, debug, verbose, file, projectName) //create hashes and check in db
			dht.ScrapeTrackers(db, debug, verbose, projectName)             //scrape trackers
			dht.SearchDHT(db, debug, verbose, projectName)                  //search DHT
			fmt.Println("Possibles (possibles hashes to the file)")
			err := utils.SelectPossiblesWhere("possible", db, debug, verbose, projectName) //get only possibles
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Valid (found in DHT)")
			err = utils.SelectPossiblesWhere("valid", db, debug, verbose, projectName) //get only valids
			if err != nil {
				log.Println(err)
			}
		case 2:
			dht.WorkersTorrents(cfgFile, debug, verbose, file, projectName) //create hashes and check in db
			dht.ScrapeTrackers(db, debug, verbose, projectName)             //scrape trackers
			//dht.ScrapeDHT(db, debug, verbose)                  //search DHT

			fmt.Println("Possibles (possibles hashes to the file)")
			err := utils.SelectPossiblesWhere("possible", db, debug, verbose, projectName) //get only possibles
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Valid (found in DHT)")
			err = utils.SelectPossiblesWhere("valid", db, debug, verbose, projectName) //get only valids
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Download valid")
			err = utils.DownloadValid(db, timeout, debug, verbose, projectName) //download valid
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Downloaded")
			err = utils.SelectPossiblesWhere("download", db, debug, verbose, projectName) //get only download
			if err != nil {
				log.Println(err)
			}

		case 3:
			dht.WorkersTorrents(cfgFile, debug, verbose, file, projectName) //create hashes and check in db
			dht.ScrapeTrackers(db, debug, verbose, projectName)             //scrape trackers
			//dht.ScrapeDHT(db, debug, verbose)                  //search DHT

			fmt.Println("Possibles (possibles hashes to the file)")
			err := utils.SelectPossiblesWhere("possible", db, debug, verbose, projectName) //get only possibles
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Valid (found in DHT)")
			err = utils.SelectPossiblesWhere("valid", db, debug, verbose, projectName) //get only valids
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Download possibles")
			err = utils.DownloadPossibles(db, timeout, debug, verbose, projectName) //download possibles
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Downloaded")
			err = utils.SelectPossiblesWhere("download", db, debug, verbose, projectName) //get only download
			if err != nil {
				log.Println(err)
			}

		default:
			fmt.Println("Invalid type option")
			cmd.Help()
			os.Exit(0)
		}

	},
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.PersistentFlags().StringVarP(&projectName, "projectName", "p", "default", "ProjectName to database")
	findCmd.PersistentFlags().StringVarP(&tracker, "tracker", "", "", "<not implemented> tracker")
	findCmd.PersistentFlags().BoolVarP(&noadd, "no-add", "n", false, "<not implemented> no add to database")
	findCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 5, "timeout download in minutes")
	findCmd.PersistentFlags().IntVarP(&typeF, "mode", "m", 2, `Opciones
	0) Trackers 
	1) Trackers + DHT, Bool valid
	2) Trackers + DHT + Download valid, Bool valid and download
	3) Try to download all possibles
	`)
}
