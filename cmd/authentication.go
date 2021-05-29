package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
)

type authenticationCmd struct {
	cmd *cobra.Command
}

func newAuthenticationCmd() *authenticationCmd {
	ac := &authenticationCmd{}
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
			da, err := svixClient.Authentication.DashboardAccess(appID)
			if err != nil {
				return err
			}

			pretty.Print(da, getPrintOptions(cmd))
			return nil
		},
	}
	ac.cmd.AddCommand(dashboard)

	// // logout
	// logout := &cobra.Command{
	// 	Use:   "logout",
	// 	Short: "Get a dashboard URL for the given app",
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		err :=svixClient.Authentication.Logout()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	},
	// }
	// ac.cmd.AddCommand(logout)

	return ac
}

func (ac *authenticationCmd) Cmd() *cobra.Command {
	return ac.cmd
}
