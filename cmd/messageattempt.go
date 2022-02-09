package cmd

import (
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-webhooks/go"
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
			opts, err := getMessageAttemptListOptions(cmd)
			printer.CheckErr(err)

			l, err := svixClient.MessageAttempt.List(appID, msgID, opts)
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
			opts, err := getMessageAttemptListOptions(cmd)
			printer.CheckErr(err)
			l, err := svixClient.MessageAttempt.ListAttemptedDestinations(appID, msgID, opts)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addMessageAttemptFilterFlags(listDestinations)
	mac.cmd.AddCommand(listDestinations)

	// list by endpoint
	// List Attempts For Endpoint
	listEndpoint := &cobra.Command{
		Use:     "list-attempts-for-endpoint APP_ID MSG_ID ENDPOINT_ID",
		Short:   "List attempts of the message filted by endpoint",
		Args:    validators.ExactArgs(3),
		Aliases: []string{"list-endpoint"},
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			svixClient := getSvixClientOrExit()

			opts, err := getMessageAttemptListOptions(cmd)
			printer.CheckErr(err)
			l, err := svixClient.MessageAttempt.ListAttemptsForEndpoint(appID, msgID, endpointID, opts)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addMessageAttemptFilterFlags(listEndpoint)
	mac.cmd.AddCommand(listEndpoint)

	// list all attempts for endpoint
	// List Attempts For Endpoint
	listAttemptedMessages := &cobra.Command{
		Use:   "list-attempted-messages APP_ID ENDPOINT_ID",
		Short: "List all attempts for a given endpoint",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			endpointID := args[1]

			svixClient := getSvixClientOrExit()

			opts, err := getMessageAttemptListOptions(cmd)
			printer.CheckErr(err)
			l, err := svixClient.MessageAttempt.ListAttemptedMessages(appID, endpointID, opts)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addMessageAttemptFilterFlags(listAttemptedMessages)
	mac.cmd.AddCommand(listAttemptedMessages)

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

func addMessageAttemptFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor id for list call")
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
	cmd.Flags().Int32P("status", "s", 0, "message status")
	cmd.Flags().StringArray("event-types", []string{}, "event types")
	cmd.Flags().StringP("before", "b", "", "before")
}

func getMessageAttemptListOptions(cmd *cobra.Command) (*svix.MessageAttemptListOptions, error) {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.MessageAttemptListOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if cmd.Flags().Changed("iterator") {
		opts.Iterator = &iteratorFlag
	}

	statusFlag, _ := cmd.Flags().GetInt32("status")
	if cmd.Flags().Changed("status") {
		status := svix.MessageStatus(statusFlag)
		opts.Status = &status
	}

	eventTypesFlag, _ := cmd.Flags().GetStringArray("event-types")
	if cmd.Flags().Changed("event-types") {
		opts.EventTypes = &eventTypesFlag
	}

	beforeFlag, _ := cmd.Flags().GetString("before")
	if cmd.Flags().Changed("before") {
		t, err := dateparse.ParseAny(beforeFlag)
		if err != nil {
			return nil, fmt.Errorf("invalid before flag: %s", err)
		}
		opts.Before = &t
	}

	return opts, nil
}
