package cmd

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/svix/svix-cli/config"
	"github.com/svix/svix-cli/flags"
	"github.com/svix/svix-cli/version"
	svix "github.com/svix/svix-webhooks/go"
)

var defaultApiUrl = "https://api.svix.com"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "svix",
	Short:   "A CLI to interact with the Svix API.",
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetVersionTemplate(version.String())

	// Root Flags
	rootCmd.Flags().BoolP("version", "v", false, "Get the version of the Svix CLI") // overrides default msg

	// Global Flags
	color := "auto"
	colorFlag := flags.NewEnum(&color, "auto", "always", "never")
	flag.Var(colorFlag, "color", "auto|always|never")
	rootCmd.PersistentFlags().AddGoFlag(flag.Lookup("color"))
	cobra.CheckErr(viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))) // allow color flag to be set in config

	// Register Commands
	rootCmd.AddCommand(newVersionCmd().cmd)
	rootCmd.AddCommand(newLoginCmd().cmd)
	rootCmd.AddCommand(newApplicationCmd().cmd)
	rootCmd.AddCommand(newAuthenticationCmd().cmd)
	rootCmd.AddCommand(newEventTypeCmd().cmd)
	rootCmd.AddCommand(newEndpointCmd().cmd)
	rootCmd.AddCommand(newMessageCmd().cmd)
	rootCmd.AddCommand(newMessageAttemptCmd().cmd)
	rootCmd.AddCommand(newVerifyCmd().cmd)
	rootCmd.AddCommand(newOpenCmd().cmd)
	rootCmd.AddCommand(newListenCmd().cmd)
	rootCmd.AddCommand(newImportCmd().cmd)
	rootCmd.AddCommand(newExportCmd().cmd)
	rootCmd.AddCommand(newIntegrationCmd().cmd)
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
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	_ = viper.ReadInConfig()
}

func getSvixClientOrExit() *svix.Svix {
	token := viper.GetString("auth_token")
	if token == "" {
		fmt.Fprintln(os.Stderr, "No SVIX_AUTH_TOKEN found!")
		fmt.Fprintln(os.Stderr, "Try running `svix login` to get started!")
		os.Exit(1)
	}

	opts := getSvixClientOptsOrExit()
	return svix.New(token, opts)
}

func getSvixClientOptsOrExit() *svix.SvixOptions {
	opts := &svix.SvixOptions{}
	rawServerUrl := viper.GetString("server_url")

	// fallback to debug_url for backwards compatibility
	if rawServerUrl == "" {
		rawServerUrl = viper.GetString("debug_url")
	}

	if rawServerUrl != "" {
		serverUrl, err := url.Parse(rawServerUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid server_url set: \"%s\"\n", rawServerUrl)
			os.Exit(1)
		}
		opts.ServerUrl = serverUrl
	}
	return opts
}
