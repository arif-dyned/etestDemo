package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	HostServerName = "dyned.sirl.me"
)
// MigrationSourceURL returns migration source URL
func MigrationSourceURL() string {
	return fmt.Sprintf(
		"%s://%s",
		viper.Get("migration.driver"),
		viper.Get("migration.folder"),
	)
}

// MigrationDatabaseURL returns migration database URL
func MigrationDatabaseURL() string {
	return fmt.Sprintf(
		"%s://%s:%s@tcp(%s:%s)/%s",
		viper.Get("database.driver"),
		viper.Get("database.username"),
		viper.Get("database.password"),
		viper.Get("database.host"),
		viper.Get("database.port"),
		viper.Get("database.database"),
	)
}
