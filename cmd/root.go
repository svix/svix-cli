package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	svix "github.com/svixhq/svix-libs/go"
)

var cfgFile string

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
	// Set Config
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.svix.yaml)")
	// cfg := cfg.New() // TODO initialize
	key := os.Getenv("SVIX_KEY")
	if key == "" {
		fmt.Println("No SVIX_KEY found!")
		fmt.Println("Try setting your auth token via 'export SVIX_KEY=<AUTH_TOKEN> to get started!")
	}
	s := svix.New(os.Getenv("SVIX_KEY"), nil)

	// Global Flags
	// rootCmd.Flags().BoolP("verbose", "v", false, "Increases output, useful for debugging")

	// Register Commands
	rootCmd.AddCommand(newApplicationCmd(s).Cmd())
	rootCmd.AddCommand(newAuthenticationCmd(s).Cmd())
	rootCmd.AddCommand(newEventTypeCmd(s).Cmd())
	rootCmd.AddCommand(newEndpointCmd(s).Cmd())
	rootCmd.AddCommand(newMessageCmd(s).Cmd())
	rootCmd.AddCommand(newMessageAttemptCmd(s).Cmd())

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".svix" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".svix")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
