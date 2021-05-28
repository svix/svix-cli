package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
