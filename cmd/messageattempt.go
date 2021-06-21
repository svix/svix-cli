package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/validators"
)

type messageAttemptCmd struct {
	cmd *cobra.Command
}

func newMessageAttemptCmd() *messageAttemptCmd {
	mac := &messageAttemptCmd{}
	mac.cmd = &cobra.Command{
		Use:   "message-attempt",
		Short: "List, lookup & resend message attempts",
	}

	// list TODO add remaining list endpoints to this single command
	list := &cobra.Command{
		Use:   "list APP_ID MSG_ID",
		Short: "List attempted messages by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			msgID := args[1]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.MessageAttempt.List(appID, msgID, getFilterOptionsMessageAttempt(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addMessageAttemptFilterFlags(list)
	mac.cmd.AddCommand(list)

	// list destinations
	listDestinations := &cobra.Command{
		Use:   "list-destinations APP_ID MSG_ID",
		Short: "List attempted destinations",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			msgID := args[1]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.MessageAttempt.ListAttemptedDestinations(appID, msgID, getFilterOptions(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addFilterFlags(listDestinations)
	mac.cmd.AddCommand(listDestinations)

	// list by endpoint
	// List Attempts For Endpoint
	listEndpoint := &cobra.Command{
		Use:   "list-endpoint APP_ID MSG_ID ENDPOINT_ID",
		Short: "List attempts for endpoint",
		Args:  validators.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.MessageAttempt.ListAttemptsForEndpoint(appID, msgID, endpointID, *getFilterOptionsMessageAttempt(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addMessageAttemptFilterFlags(listEndpoint)
	mac.cmd.AddCommand(listEndpoint)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID MSG_ID ATTEMPT_ID",
		Short: "Get attempt by id",
		Args:  validators.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse args
			appID := args[0]
			msgID := args[1]
			attemptID := args[2]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.MessageAttempt.Get(appID, msgID, attemptID)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	mac.cmd.AddCommand(get)

	// resend
	resend := &cobra.Command{
		Use:   "resend APP_ID MSG_ID ENDPOINT_ID",
		Short: "resends a webhook message by id",
		Args:  validators.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse args
			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			svixClient := getSvixClientOrExit()
			err := svixClient.MessageAttempt.Resend(appID, msgID, endpointID)
			printer.CheckErr(err)
		},
	}
	mac.cmd.AddCommand(resend)

	return mac
}
