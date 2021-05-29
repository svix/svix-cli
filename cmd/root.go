package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	svix "github.com/svixhq/svix-libs/go"
)

var svixClient *svix.Svix

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "svix",
	Short: "A CLI to interact with the Svix API.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global Flags
	// rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Increases output, useful for debugging")
	rootCmd.PersistentFlags().Bool("json", false, "output results in json if possible")

	// Register Commands
	rootCmd.AddCommand(newApplicationCmd().Cmd())
	rootCmd.AddCommand(newAuthenticationCmd().Cmd())
	rootCmd.AddCommand(newEventTypeCmd().Cmd())
	rootCmd.AddCommand(newEndpointCmd().Cmd())
	rootCmd.AddCommand(newMessageCmd().Cmd())
	rootCmd.AddCommand(newMessageAttemptCmd().Cmd())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Find home directory.
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	// Search config in home directory with name ".svix" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".svix")

	// read in environment variables that match
	viper.SetEnvPrefix("svix")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	key := viper.GetString("key")
	if key == "" {
		fmt.Println("No SVIX_KEY found!")
		fmt.Println("Try running `svix init` to get started!")
		os.Exit(1)
	}
	svixClient = svix.New(key, nil)
}
