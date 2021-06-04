package cmd

import (
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/validators"
)

const docsURL = "https://docs.svix.com/"

type docsCmd struct {
	cmd *cobra.Command
}

func newDocsCmd() *versionCmd {
	return &versionCmd{
		cmd: &cobra.Command{
			Use:   "docs",
			Args:  validators.NoArgs(),
			Short: "Open the default browser with the Svix documentation",
			Run: func(cmd *cobra.Command, args []string) {
				open.Run(docsURL)
			},
		},
	}
}
