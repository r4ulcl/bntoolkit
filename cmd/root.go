package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string //config file
var verbose bool
var debug bool

var rootCmd = &cobra.Command{
	Use:   "bntoolkit <command>",
	Short: "BNT (BitTorrent Network toolkit)",
	Long:  `Set of utilities to monitor, download, create and find files in the BitTorent network`,
	Args:  cobra.MinimumNArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//generate doc
	/*err := doc.GenMarkdownTree(rootCmd, "./doc")
	if err != nil {
		log.Fatal(err)
	}*/

	gopath := os.Getenv("GOPATH")

	//Persistent Flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "cfgFile", "c", gopath+"/src/github.com/RaulCalvoLaorden/bntoolkit/configFile.toml", "Config file to DB (host, port, user, password)")
}
