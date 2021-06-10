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
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.MessageAttempt.List(appID, msgID, getFilterOptionsMessageAttempt(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
		},
	}
	addMessageAttemptFilterFlags(list)
	mac.cmd.AddCommand(list)

	// list destinations
	listDestinations := &cobra.Command{
		Use:   "list-destinations APP_ID MSG_ID",
		Short: "List attempted destinations",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.MessageAttempt.ListAttemptedDestinations(appID, msgID, getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.MessageAttempt.ListAttemptsForEndpoint(appID, msgID, endpointID, *getFilterOptionsMessageAttempt(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
		},
	}
	addMessageAttemptFilterFlags(listEndpoint)
	mac.cmd.AddCommand(listEndpoint)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID MSG_ID ATTEMPT_ID",
		Short: "Get attempt by id",
		Args:  validators.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			msgID := args[1]
			attemptID := args[2]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.MessageAttempt.Get(appID, msgID, attemptID)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	mac.cmd.AddCommand(get)

	// resend
	resend := &cobra.Command{
		Use:   "resend APP_ID MSG_ID ENDPOINT_ID",
		Short: "resends a webhook message by id",
		Args:  validators.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			svixClient := getSvixClientOrExit()
			err := svixClient.MessageAttempt.Resend(appID, msgID, endpointID)
			if err != nil {
				return err
			}
			return nil
		},
	}
	mac.cmd.AddCommand(resend)

	return mac
}
