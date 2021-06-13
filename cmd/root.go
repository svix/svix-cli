package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/svix/svix-cli/config"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/version"
	svix "github.com/svix/svix-libs/go"
)

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
	isTTY, _, err := utils.IsTTY(os.Stdout)
	cobra.CheckErr(err)
	rootCmd.PersistentFlags().Bool("color", isTTY, "colorize output json")              // on by default if TTY, off if not
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
	rootCmd.AddCommand(newCompletionCmd().cmd)
	rootCmd.AddCommand(newOpenCmd().cmd)
	rootCmd.AddCommand(newListenCmd().cmd)
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

	opts := &svix.SvixOptions{}
	rawDebugURL := viper.GetString("debug_url")
	if rawDebugURL != "" {
		debugURL, err := url.Parse(rawDebugURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid debug_url set: \"%s\"\n", rawDebugURL)
			os.Exit(1)
		}
		opts.DebugURL = debugURL
	}
	return svix.New(token, opts)
}
