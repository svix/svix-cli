package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/cfg"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type authenticationCmd struct {
	cmd *cobra.Command
	cfg *cfg.Config
	sc  *svix.Svix
}

func newAuthenticationCmd(cfg *cfg.Config, s *svix.Svix) *authenticationCmd {
	ac := &authenticationCmd{
		sc:  s,
		cfg: cfg,
	}
	ac.cmd = &cobra.Command{
		Use:     "authentication",
		Short:   "Get Dashboard urls",
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

	//TODO do we need to use logout from the CLI?

	return ac
}

func (ac *authenticationCmd) Cmd() *cobra.Command {
	return ac.cmd
}
