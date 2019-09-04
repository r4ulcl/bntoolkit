package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/spf13/cobra"
)

var path string

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <file/hash/magnet>",
	Short: "Download a file from a hash, a magnet or a Torrent file",
	Long: `Download a file from a hash, a magnet or a Torrent file. 
For example:
	bntoolkit download e84213a794f3ccd890382a54a64ca68b7e925433
	bntoolkit download ubuntu.torrent`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called", len(args))

		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}

		if path == "" {
			path = "."
		}

		file := args[0]
		err := download(file, path)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func download(file string, path string) error {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = path
	cfg.NoUpload = true // Ensure that downloads are responsive.

	fmt.Println(file)
	c, err := torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	defer c.Close()

	var t *torrent.Torrent
	if strings.Contains(file, "magnet:?xt=urn:btih:") {
		t, _ = c.AddMagnet(file)
	} else if strings.Contains(file, ".torrent") {
		t, _ = c.AddTorrentFromFile(file)
	} else {
		t, _ = c.AddMagnet("magnet:?xt=urn:btih:" + file)
	}

	<-t.GotInfo()
	fmt.Println("Downloading", t.Name())
	t.DownloadAll()
	c.WaitAll()
	log.Print("Torrent downloaded")

	return nil
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().StringVarP(&download, "magnet", "m", "", "Magnet, hash or .torrent to download")

	downloadCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "path to download")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
