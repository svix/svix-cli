package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/validators"
)

type authenticationCmd struct {
	cmd *cobra.Command
}

func newAuthenticationCmd() *authenticationCmd {
	ac := &authenticationCmd{}
	ac.cmd = &cobra.Command{
		Use:     "authentication",
		Short:   "Manage authentication tasks such as getting dashboard URLs",
		Aliases: []string{"auth"},
	}

	// dashboard
	dashboard := &cobra.Command{
		Use:   "dashboard APP_ID",
		Short: "Get a dashboard URL for the given app",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]

			svixClient := getSvixClientOrExit()
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
	// 		svixClient := getSvixClientOrExit()
	// 		err := svixClient.Authentication.Logout()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	},
	// }
	// ac.cmd.AddCommand(logout)

	return ac
}
