package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svixhq/svix-cli/config"
)

type initCmd struct {
	cmd *cobra.Command
}

func newInitCmd() *initCmd {
	ic := &initCmd{}
	ic.cmd = &cobra.Command{
		Use:   "init",
		Short: "Interactively configure your Svix API credentials",
		Args:  cobra.ExactArgs(0),
		Run:   ic.run,
	}
	return ic
}

func (ic *initCmd) run(cmd *cobra.Command, args []string) {
	fmt.Printf("Welcome to the Svix CLI, enter your API key to get started!\n\n")

	// get api key
	keyPrompt := promptui.Prompt{
		Label:   "Svix API Key",
		Default: viper.GetString("key"),
		Validate: func(input string) error {
			if !strings.HasPrefix(input, "dhk_") {
				// TODO should this be enforced?
				return fmt.Errorf("API keys must start with dhk_")
			}
			if len(input) != 36 {
				// TODO is this true?
				return fmt.Errorf("API keys should be 36 characters long")
			}
			return nil
		},
	}
	apiKey, err := keyPrompt.Run()
	if err != nil {
		fmt.Printf("Initialization failed %v\n", err)
		return
	}
	viper.Set("key", apiKey)

	if err := config.Write(viper.AllSettings()); err != nil {
		fmt.Println(err)
		fmt.Println("Failed to configure the Svix CLI, please try again or try setting your api key manually 'SVIX_KEY' environment variable.")
		return
	}

	fmt.Printf("All Set! Your config has been written to \"%s\"\n", viper.ConfigFileUsed())
	fmt.Println("Type `svix --help` to print the Svix CLI documentation!")
}

func (ic *initCmd) Cmd() *cobra.Command {
	return ic.cmd
}
