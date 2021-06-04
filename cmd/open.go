package cmd

import (
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/validators"
)

const docsURL = "https://docs.svix.com/"
const apidocsURL = "https://api.svix.com/docs"

func newOpenCmd() *versionCmd {
	return &versionCmd{
		cmd: &cobra.Command{
			Use:       "open [docs|api]",
			ValidArgs: []string{"docs", "api"},
			Args:      validators.ExactValidArgs(1),
			Short:     "Open in browser",
			Long: `Opens information in default browser:
docs - opens the Svix documentation
api  - opens the Svix API documentation
			`,
			Run: func(cmd *cobra.Command, args []string) {
				url := ""
				switch args[0] {
				case "docs":
					url = docsURL
				case "api":
					url = apidocsURL
				}
				err := open.Run(url)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					fmt.Fprintln(os.Stderr, "Failed to open default application for", makeTerminalHyperlink(url, url))
					os.Exit(1)
				}
			},
		},
	}
}

// TODO: Enable this function in pretty package and switch to using that above
func makeTerminalHyperlink(name, url string) string {
	return fmt.Sprintf("\u001B]8;;%s\a%s\u001B]8;;\a", url, name)
}
