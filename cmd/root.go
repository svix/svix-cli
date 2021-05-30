package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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

	// Setup config file
	configFolder, err := config.Folder()
	cobra.CheckErr(err)

	configFile := filepath.Join(configFolder, config.FileName)
	viper.SetConfigType("toml")
	viper.SetConfigFile(configFile)
	viper.SetConfigPermissions(config.FileMode)

	// read in environment variables that match
	viper.SetEnvPrefix("svix")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
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

	opts := &svix.SvixOptions{}
	rawBaseURL := viper.GetString("base-url")
	if rawBaseURL != "" {
		baseURL, err := url.Parse(rawBaseURL)
		if err != nil {
			log.Printf("Invalid base-url set: \"%s\"\n", rawBaseURL)
			os.Exit(1)
		}
		opts.BaseURL = baseURL
	}
	return svix.New(key, opts)
}
