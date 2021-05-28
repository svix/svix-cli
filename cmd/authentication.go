package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type authenticationCmd struct {
	cmd *cobra.Command
	sc  *svix.Svix
}

func newAuthenticationCmd(s *svix.Svix) *authenticationCmd {
	ac := &authenticationCmd{
		sc: s,
	}
	ac.cmd = &cobra.Command{
		Use:     "authentication",
		Short:   "Manage authentication tasks such getting dashboard urls",
		Aliases: []string{"auth"},
	}

	// dashboard
	dashboard := &cobra.Command{
		Use:   "dashboard APP_ID",
		Short: "Get a dashboard URL for the given app",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			l, err := s.Authentication.DashboardAccess(appID)
			if err != nil {
				return err
			}

			pretty.PrintDashboardURL(appID, l.Url)
			return nil
		},
	}
	ac.cmd.AddCommand(dashboard)

	return ac
}

func (ac *authenticationCmd) Cmd() *cobra.Command {
	return ac.cmd
}
