package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/validators"
	"github.com/svix/svix-cli/version"
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
