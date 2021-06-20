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
		Run: func(cmd *cobra.Command, args []string) {
			appID := args[0]
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			svixClient := getSvixClientOrExit()
			da, err := svixClient.Authentication.DashboardAccess(appID)
			printer.CheckErr(err)

			printer.Print(da)
		},
	}
	ac.cmd.AddCommand(dashboard)

	// // logout
	// logout := &cobra.Command{
	// 	Use:   "logout",
	// 	Short: "Get a dashboard URL for the given app",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		svixClient := getSvixClientOrExit()
	// 		err := svixClient.Authentication.Logout()
	// 		printer.CheckErr(err)
	// 	},
	// }
	// ac.cmd.AddCommand(logout)

	return ac
}
