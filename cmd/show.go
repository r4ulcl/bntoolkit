package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"
	"github.com/spf13/cobra"
)

var sql string
var hash string
var source string

var name string

var user string

var ip string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the database data",
	Long:  `Show the database data`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")
		cmd.Help()
		os.Exit(0)
	},
}

var showHashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Show the database data of the table hash",
	Long:  `Show the database data of the table hash`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""
		if sql != "" {
			salida, err = utils.QueryHash(db, debug, verbose, sql, "", "")
		} else if hash != "" {
			if source != "" {
				salida, err = utils.QueryHash(db, debug, verbose, "", hash, source)
			} else {
				salida, err = utils.QueryHash(db, debug, verbose, "", hash, "")
			}
		} else if source != "" {
			salida, err = utils.QueryHash(db, debug, verbose, "", "", source)
		} else {
			salida, err = utils.QueryHash(db, debug, verbose, "", "", "")
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

var showpossiblesCmd = &cobra.Command{
	Use:   "possibles",
	Short: "Show the database data of the table possibles",
	Long:  `Show the database data of the table possibles`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""
		if sql != "" {
			salida, err = utils.QueryPossibles(db, debug, verbose, sql, "")
		} else if hash != "" {
			salida, err = utils.QueryPossibles(db, debug, verbose, "", hash)
		} else {
			salida, err = utils.QueryPossibles(db, debug, verbose, "", "")
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

var showProjectsCmd = &cobra.Command{
	Use:   "project",
	Short: "Show the database data of the table project",
	Long:  `Show the database data of the table project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("projects called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""
		if sql != "" {
			salida, err = utils.QueryProjects(db, debug, verbose, sql, "")
		} else if name != "" {
			salida, err = utils.QueryProjects(db, debug, verbose, "", name)
		} else {
			salida, err = utils.QueryProjects(db, debug, verbose, "", "")
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

var showMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Show the database data of the table monitor",
	Long:  `Show the database data of the table monitor`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""
		if sql != "" {
			salida, err = utils.QueryMonitor(db, debug, verbose, sql, "", "")
		} else if hash != "" {
			if user != "" {
				salida, err = utils.QueryMonitor(db, debug, verbose, "", hash, user)
			} else {
				salida, err = utils.QueryMonitor(db, debug, verbose, "", hash, "")
			}
		} else if user != "" {
			salida, err = utils.QueryMonitor(db, debug, verbose, "", "", user)
		} else {
			salida, err = utils.QueryMonitor(db, debug, verbose, "", "", "")
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

var showCountCmd = &cobra.Command{
	Use:   "count",
	Short: "Show the database count of hash",
	Long:  `Show the database count of hash`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""

		salida, err = utils.QueryCount(db, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

var showIPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Show the database data of the table ip",
	Long:  `Show the database data of the table ip`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""
		if sql != "" {
			salida, err = utils.QueryIP(db, debug, verbose, sql, "")
		} else if ip != "" {
			salida, err = utils.QueryIP(db, debug, verbose, "", ip)
		} else {
			salida, err = utils.QueryIP(db, debug, verbose, "", "")
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

var showAlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Show the database data of the table alert",
	Long:  `Show the database data of the table alert`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("results called")

		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		salida := ""
		if sql != "" {
			salida, err = utils.QueryAlert(db, debug, verbose, sql, "", "")
		} else if ip != "" {
			if user != "" {
				salida, err = utils.QueryAlert(db, debug, verbose, "", ip, user)
			} else {
				salida, err = utils.QueryAlert(db, debug, verbose, "", ip, "")
			}
		} else if user != "" {
			salida, err = utils.QueryAlert(db, debug, verbose, "", "", user)
		} else {
			salida, err = utils.QueryAlert(db, debug, verbose, "", "", "")
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(salida)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.AddCommand(showAlertCmd)
	showCmd.AddCommand(showHashCmd)
	showCmd.AddCommand(showIPCmd)
	showCmd.AddCommand(showMonitorCmd)
	showCmd.AddCommand(showpossiblesCmd)
	showCmd.AddCommand(showProjectsCmd)
	showCmd.AddCommand(showCountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	showCmd.PersistentFlags().StringVarP(&sql, "where", "w", "", "Where in SQL lenguaje")
	showHashCmd.PersistentFlags().StringVarP(&hash, "hash", "", "", "Hash to filter")
	showHashCmd.PersistentFlags().StringVarP(&source, "source", "", "", "Source to filter")

	showpossiblesCmd.PersistentFlags().StringVarP(&hash, "hash", "", "", "Hash to filter")

	showProjectsCmd.PersistentFlags().StringVarP(&name, "projectName", "p", "", "Project name")

	showIPCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "", "ip")

	showMonitorCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user")

	showAlertCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "", "ip")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resultsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
