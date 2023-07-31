package cmd

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-webhooks/go"
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

	ctx := context.Background()

	// dashboard -- deprecated
	dashboard := &cobra.Command{
		Use:        "dashboard-access APP_ID",
		Short:      "Get a dashboard URL for the given app",
		Args:       validators.ExactArgs(1),
		Aliases:    []string{"dashboard"},
		Deprecated: "use app-portal instead",
		Run: func(cmd *cobra.Command, args []string) {
			appID := args[0]
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			svixClient := getSvixClientOrExit()
			da, err := svixClient.Authentication.DashboardAccess(ctx, appID)
			printer.CheckErr(err)

			printer.Print(da)
		},
	}
	ac.cmd.AddCommand(dashboard)

	// logout
	logout := &cobra.Command{
		Use:   "logout DASHBOARD_AUTH_TOKEN",
		Short: "Invalidates the given dashboard key",
		Args:  validators.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			authToken := args[0]
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			svixClient := svix.New(authToken, getSvixClientOptsOrExit())
			err := svixClient.Authentication.Logout(ctx)
			printer.CheckErr(err)
		},
	}
	ac.cmd.AddCommand(logout)

	// app-portal
	appPortal := &cobra.Command{
		Use:   "app-portal APP_ID [JSON_PAYLOAD]",
		Short: "Get app portal access for the given app ID",
		Args:  validators.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			appID := args[0]
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			var payloadBytes []byte

			if len(args) > 1 {
				payloadBytes = []byte(args[1])
			} else if readable, _ := utils.IsStdinReadable(); readable {
				bytes, err := utils.ReadStdin()
				printer.CheckErr(err)
				payloadBytes = bytes
			}

			var appPortalAccessIn svix.AppPortalAccessIn
			if len(payloadBytes) > 0 {
				err := json.Unmarshal(payloadBytes, &appPortalAccessIn)
				printer.CheckErr(err)
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Authentication.AppPortalAccess(ctx, appID, &appPortalAccessIn)

			printer.CheckErr(err)
			printer.Print(out)
		},
	}
	ac.cmd.AddCommand(appPortal)

	return ac
}
