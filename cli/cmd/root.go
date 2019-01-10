package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appVersion = "v0.1.0"
	appCluster = "master"
	cfgFile    string
)

// RootCmd represents the base command when called without any subcommands
// The base command does not have handler
var RootCmd = &cobra.Command{
	Use:   "etest",
	Short: "DynEd E-Test Application (" + appVersion + ")",
	Long:  "DynEd E-Test Application (" + appVersion + ")",
}

// Execute adds all child commands to the RootCmd command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if cfgFile == "" {
			cfgFile = "./config.json"
		}

		// TODO: handle when config file is not found
		viper.SetConfigFile("./config.json")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(errors.Wrap(err, "[config]"))
		}
	})

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.json)")
}
