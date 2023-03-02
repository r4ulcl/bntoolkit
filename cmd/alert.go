package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/r4ulcl/bntoolkit/utils"
	"github.com/spf13/cobra"
)

var projectName string

// addAlertCmd represents the addAlert command
var addAlertCmd = &cobra.Command{
	Use:   "addAlert <IP/Range>",
	Short: "Add an IP to the database alert table",
	Long: `Add an IP or range to the database alert table. When the daemon is executed an alert appears if that IP appears in the alerts table (IP or range).
For example:
		bntoolkit addAlert 1.1.0.0/16`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addAlert called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		ip := args[0]

		utils.InsertProject(db, debug, verbose, projectName)
		utils.InsertIP(db, debug, verbose, ip, projectName) //if return error exist

		err = utils.InsertAlert(db, debug, verbose, ip, "test", "username", projectName)
		if err != nil {
			log.Println(err)
		}
	},
}

var deleteAlertCmd = &cobra.Command{
	Use:   "deleteAlert <IP/Range>",
	Short: "Delete an IP or range from the database alert table. ",
	Long: `Delete an IP or range from the database alert table. 
For example:
	bntoolkit deleteAlert 1.1.0.0/16`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deleteAlert called")
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		ip := args[0]
		err = utils.DeleteAlert(db, debug, verbose, ip)
		if err != nil {
			log.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addAlertCmd)
	rootCmd.AddCommand(deleteAlertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addAlertCmd.PersistentFlags().String("foo", "", "A help for foo")
	addAlertCmd.PersistentFlags().StringVarP(&projectName, "projectName", "p", "default", "Monitoring projectName")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addAlertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
