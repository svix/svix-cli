package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svixhq/svix-cli/config"
)

type loginCmd struct {
	cmd *cobra.Command
}

func newLoginCmd() *loginCmd {
	lc := &loginCmd{}
	lc.cmd = &cobra.Command{
		Use:   "init",
		Short: "Interactively configure your Svix API credentials",
		Args:  cobra.ExactArgs(0),
		Run:   lc.run,
	}
	return lc
}

func (lc *loginCmd) run(cmd *cobra.Command, args []string) {
	fmt.Printf("Welcome to the Svix CLI, enter your API key to get started!\n\n")

	// get api key
	keyPrompt := promptui.Prompt{
		Label:   "Svix API Key",
		Default: viper.GetString("key"),
	}
	apiKey, err := keyPrompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed %v\n", err)
		os.Exit(1)
	}
	viper.Set("key", apiKey)

	if err := config.Write(viper.AllSettings()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "Failed to configure the Svix CLI, please try again or try setting your api key manually 'SVIX_KEY' environment variable.")
		os.Exit(1)
	}

	fmt.Printf("All Set! Your config has been written to \"%s\"\n", viper.ConfigFileUsed())
	fmt.Println("Type `svix --help` to print the Svix CLI documentation!")
}

func (lc *loginCmd) Cmd() *cobra.Command {
	return lc.cmd
}
