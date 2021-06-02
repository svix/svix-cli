package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svixhq/svix-cli/config"
	"github.com/svixhq/svix-cli/validators"
)

type loginCmd struct {
	cmd *cobra.Command
}

func newLoginCmd() *loginCmd {
	lc := &loginCmd{}
	lc.cmd = &cobra.Command{
		Use:   "login",
		Short: "Interactively configure your Svix API credentials",
		Args:  validators.NoArgs(),
		Run:   lc.run,
	}
	return lc
}

func (lc *loginCmd) run(cmd *cobra.Command, args []string) {
	fmt.Printf("Welcome to the Svix CLI, enter your auth token to get started!\n\n")

	// get auth token
	keyPrompt := promptui.Prompt{
		Label:   "Svix Auth Token",
		Default: viper.GetString("auth_token"),
	}
	token, err := keyPrompt.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Initialization failed %v\n", err)
		os.Exit(1)
	}
	viper.Set("auth_token", token)

	if err := config.Write(viper.AllSettings()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, "Failed to configure the Svix CLI, please try again or try setting your auth token manually 'SVIX_AUTH_TOKEN' environment variable.")
		os.Exit(1)
	}

	fmt.Printf("All Set! Your config has been written to \"%s\"\n", viper.ConfigFileUsed())
	fmt.Println("Type `svix --help` to print the Svix CLI documentation!")
}
