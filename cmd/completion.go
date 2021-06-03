package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/validators"
)

type completionCmd struct {
	cmd *cobra.Command
}

func newCompletionCmd() *completionCmd {
	return &completionCmd{
		cmd: &cobra.Command{
			Use:   "completion [bash|zsh|fish|powershell]",
			Short: "Generate completion script",
			Long: `To load completions:

Bash:

  $ source <(svix completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ svix completion bash > /etc/bash_completion.d/svix
  # macOS:
  $ svix completion bash > /usr/local/etc/bash_completion.d/svix

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ svix completion zsh > "${fpath[1]}/_svix"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ svix completion fish | source

  # To load completions for each session, execute once:
  $ svix completion fish > ~/.config/fish/completions/svix.fish

PowerShell:

  PS> svix completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> svix completion powershell > svix.ps1
  # and source this file from your PowerShell profile.
`,
			DisableFlagsInUseLine: true,
			ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
			Args:                  validators.ExactValidArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				switch args[0] {
				case "bash":
					cobra.CheckErr(cmd.Root().GenBashCompletion(os.Stdout))
				case "zsh":
					cobra.CheckErr(cmd.Root().GenZshCompletion(os.Stdout))
				case "fish":
					cobra.CheckErr(cmd.Root().GenFishCompletion(os.Stdout, true))
				case "powershell":
					cobra.CheckErr(cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout))
				}
			},
		},
	}
}
