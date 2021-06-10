package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/validators"
)

var openableURLs = map[string]string{
	"docs": "https://docs.svix.com/",
	"api":  "https://api.svix.com/docs",
}

type openCmd struct {
	cmd *cobra.Command
}

func keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func newOpenCmd() *openCmd {
	keys := keys(openableURLs)
	oc := &openCmd{
		cmd: &cobra.Command{
			Use:       fmt.Sprintf("open [%s]", strings.Join(keys, "|")),
			ValidArgs: keys,
			Args:      validators.ExactValidArgs(1),
			Short:     "Quickly open Svix pages in your browser",
			Long: `Quickly open Svix pages in your browser:
docs - opens the Svix documentation
api  - opens the Svix API documentation
			`,
			Run: func(cmd *cobra.Command, args []string) {
				url := openableURLs[args[0]]
				err := open.Run(url)
				if err != nil {
					fmt.Fprintf(os.Stderr, `Failed to open %s in your default browser
To open it manually navigate to:
%s
`, args[0], pretty.MakeTerminalLink(url, url))
					os.Exit(1)
				}
			},
		},
	}
	return oc
}
