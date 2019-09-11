package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/spf13/cobra"
)

var projectNameName string
var userName string

// addMonitorCmd represents the addMonitor command
var addMonitorCmd = &cobra.Command{
	Use:   "addMonitor <HASH>",
	Short: "Add a hash to the database monitor table",
	Long: `Add a hash to the database monitor table. 
For example:
	bntoolkit addMonitor e84213a794f3ccd890382a54a64ca68b7e925433`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addMonitor called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		hash := args[0]
		if len(hash) == 40 {
			utils.InsertProject(db, debug, verbose, projectName)

			utils.InsertHash(db, debug, verbose, hash, "cli")
			utils.InsertMonitor(db, debug, verbose, hash, userName, projectName)
		} else {
			fmt.Println("Error hash incorrect")
			cmd.Help()
			os.Exit(0)
		}
	},
}

var deleteMonitorCmd = &cobra.Command{
	Use:   "deleteMonitor <HASH>",
	Short: "Delete a hash from the database monitor table.",
	Long: `Delete a hash from the database monitor table.
For example:
	bntoolkit deleteMonitor e84213a794f3ccd890382a54a64ca68b7e925433`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteMonitor called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		hash := args[0]
		if len(hash) == 40 {
			utils.DeleteMonitor(db, debug, verbose, hash)
			utils.DeleteProject(db, debug, verbose, projectName)
		}
	},
}

func init() {
	rootCmd.AddCommand(addMonitorCmd)
	rootCmd.AddCommand(deleteMonitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	addMonitorCmd.PersistentFlags().StringVarP(&projectName, "projectName", "p", "default", "Monitoring projectName")
	addMonitorCmd.PersistentFlags().StringVarP(&userName, "userName", "u", "default", "Monitoring username")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addMonitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
