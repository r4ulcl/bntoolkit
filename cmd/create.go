package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/RaulCalvoLaorden/bntoolkit/dht"
	"github.com/spf13/cobra"
)

var outfile string
var piecesize int64
var tracker string
var comment string
var private bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <file/folder>",
	Short: "Create a .torrent file",
	Long: `Create a .torrent file. You can specify the output file, the pieze size, the tracker and a comment
For example:
	bntoolkit create ubuntu.iso -o ubuntuTorrent -m test`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}

		file = strings.Join(args, " ")

		dht.CrateTorrent(debug, verbose, file, outfile, piecesize, tracker, comment)
		//dht.CrateTorrent2()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", "output", "Save the generated .torrent to this filename")
	//createCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "File to create")
	createCmd.PersistentFlags().Int64VarP(&piecesize, "piecesize", "p", 2*1024, "Save the generated .torrent to this filename")
	createCmd.PersistentFlags().StringVarP(&tracker, "tracker", "t", "", "Save the generated .torrent to this filename")
	createCmd.PersistentFlags().StringVarP(&comment, "comment", "m", "", "Save the generated .torrent to this filename")
	//createCmd.PersistentFlags().BoolVarP(&private, "private", "p", false, "Save the generated .torrent to this filename")
}
