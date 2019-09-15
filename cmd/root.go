package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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
		//generate doc
		err := doc.GenMarkdownTree(rootCmd, "./doc")
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

func init() {
	gopath := os.Getenv("GOPATH")

	//Persistent Flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug output")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "cfgFile", "c", gopath+"/src/github.com/RaulCalvoLaorden/bntoolkit/configFile.toml", "Config file to DB (host, port, user, password)")
}
