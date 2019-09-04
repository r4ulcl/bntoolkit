package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "v0.9"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long: `Print the version number.
For example:
	bntoolkit version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bntoolkit (BitTorrent Network toolkit)", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	////insertCmd.Flags().StringVarP(&fileInsert, "file", "f", "", "File to insert")
	////insertCmd.Flags().StringVarP(&hashinsert, "hash", "H", "", "Hash to insert")
}
