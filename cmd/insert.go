package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/RaulCalvoLaorden/bntoolkit/utils"

	"github.com/spf13/cobra"
)

var fileInsert string
var hashinsert string

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert <hash/file>",
	Short: "Insert a hash or a file of hashes in the DB",
	Long: `Insert a hash or a file of hashes in the DB.
For example:
	bntoolkit insert hashes.txt`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("insert called")
		db, err := utils.ConnectDb(cfgFile, debug, verbose)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		argument := args[0]
		if _, err := os.Stat(argument); !os.IsNotExist(err) {
			// path/to/whatever exists
			utils.InsertFile(db, debug, verbose, argument)
		} else {
			utils.InsertHash(db, debug, verbose, argument, "manual")
		}

	},
}

func init() {
	rootCmd.AddCommand(insertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// insertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	////insertCmd.Flags().StringVarP(&fileInsert, "file", "f", "", "File to insert")
	////insertCmd.Flags().StringVarP(&hashinsert, "hash", "H", "", "Hash to insert")
}
