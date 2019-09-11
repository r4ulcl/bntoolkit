package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/spf13/cobra"
)

//var sha1 string

// getTorrentCmd represents the download command
var getTorrentCmd = &cobra.Command{
	Use:   "getTorrent <hash/magnet>",
	Short: "Get torrent file from a hash or magnet",
	Long: `Get torrent file from a hash or magnet. 
For example:
	bntoolkit getTorrent e84213a794f3ccd890382a54a64ca68b7e925433 -p /tmp/`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getTorrent called", len(args))

		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}

		aux := args[0]

		fmt.Println(aux)

		var magnet string
		if strings.Contains(aux, "magnet:?xt=urn:btih:") {
			magnet = aux
		} else {
			magnet = "magnet:?xt=urn:btih:" + aux
		}
		err := getTorrent(magnet)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func getTorrent(magnet string) error {
	cfg := torrent.NewDefaultClientConfig()
	cl, err := torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	http.HandleFunc("/torrent", func(w http.ResponseWriter, r *http.Request) {
		cl.WriteStatus(w)
	})
	http.HandleFunc("/dht", func(w http.ResponseWriter, r *http.Request) {
		for _, ds := range cl.DhtServers() {
			ds.WriteStatus(w)
		}
	})

	t, err := cl.AddMagnet(magnet)
	if err != nil {
		log.Fatalf("error adding magnet to client: %s", err)
	}
	<-t.GotInfo()
	mi := t.Metainfo()
	t.Drop()
	fmt.Println("Name: " + t.Info().Name + "\n")
	files := t.Info().Files
	fmt.Println("FILES: ")
	for _, element := range files {
		fmt.Printf("\tName: %v , size: %d\n", element.Path, element.Length)
	}

	if path != "" {
		f, err := os.Create(t.Info().Name + ".torrent")
		if err != nil {
			return err
		}
		defer f.Close()
		err = bencode.NewEncoder(f).Encode(mi)
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(getTorrentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	//getTorrentCmd.PersistentFlags().StringVarP(&sha1, "sha1", "s", "", "sha1 hash to download torrent file")

	getTorrentCmd.PersistentFlags().StringVarP(&path, "path", "p", ".", "path to save torrent file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
