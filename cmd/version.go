package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/validators"
	"github.com/svixhq/svix-cli/version"
)

type versionCmd struct {
	cmd *cobra.Command
}

func newVersionCmd() *versionCmd {
	return &versionCmd{
		cmd: &cobra.Command{
			Use:   "version",
			Args:  validators.NoArgs(),
			Short: "Get the version of the Svix CLI",
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(version.String())
			},
		},
	}
}
