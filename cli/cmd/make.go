package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	migrationTimeFormat = "20060102150405"
	migrationExt        = ".sql"
)

// Make migration command
var makeMigrationCmd = &cobra.Command{
	Use: "make:migration [name]",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Not enough arguments (missing: 'name')")
		}
		return nil
	},
	Short: "Create a new migration file",
	Long:  "Create a new migration file",
	Run: func(cmd *cobra.Command, args []string) {
		base := fmt.Sprintf("%s%s_%v.", viper.GetString("migration.folder")+"/", time.Now().Format(migrationTimeFormat), args[0])

		os.MkdirAll(viper.GetString("migration.folder"), os.ModePerm)

		createFile(base + "down" + migrationExt)
		createFile(base + "up" + migrationExt)
	},
}

func init() {
	RootCmd.AddCommand(makeMigrationCmd)
}

func createFile(filename string) {
	if _, err := os.Create(filename); err != nil {
		fmt.Println(errors.Wrap(err, "[make:migration]"))
	}
}
