package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/svixhq/svix-cli/config"
	svix "github.com/svixhq/svix-libs/go"
)

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
	rootCmd.AddCommand(newInitCmd().Cmd())
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
	path, err := config.Path()
	cobra.CheckErr(err)

	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)

	// read in environment variables that match
	viper.SetEnvPrefix("svix")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
}

func getSvixClientOrExit() *svix.Svix {
	key := viper.GetString("key")
	if key == "" {
		fmt.Println("No SVIX_KEY found!")
		fmt.Println("Try running `svix init` to get started!")
		os.Exit(1)
	}
	return svix.New(key, nil)
}
