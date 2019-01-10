package cmd

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	// Migration database driver dialect
	_ "github.com/golang-migrate/migrate/database/mysql"
	// Migration source
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/DynEd/etest/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO: handling environment when set to production: showing confirmation?
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run the database migrations",
	Long:  "Run the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		m, err := migrate.New(config.MigrationSourceURL(), config.MigrationDatabaseURL())
		if err != nil {
			fmt.Println(errors.Wrap(err, "[migrate]"))
			return
		}
		defer m.Close()

		version, _, _ := m.Version()
		if version == 0 {
			fmt.Printf("[migrate]: Creating schema_migrations table\n")
		} else {
			fmt.Printf("[migrate]: Current database version is %d\n", version)
		}

		m.Log = &logger{verbose: true}

		if err := m.Up(); err != nil {
			fmt.Println(errors.Wrap(err, "[migrate]"))
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
