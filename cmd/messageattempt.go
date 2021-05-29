package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type messageAttemptCmd struct {
	cmd *cobra.Command
	sc  *svix.Svix
}

func newMessageAttemptCmd(s *svix.Svix) *messageAttemptCmd {
	mac := &messageAttemptCmd{
		sc: s,
	}
	mac.cmd = &cobra.Command{
		Use:   "message-attempt",
		Short: "List, lookup & resend message attempts",
	}

	// list TODO add remaining list endpoints to this single command
	list := &cobra.Command{
		Use:   "list APP_ID MSG_ID",
		Short: "List attempted messages by id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			l, err := s.MessageAttempt.List(appID, msgID, getFilterOptionsMessageAttempt(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l)
			return nil
		},
	}
	addMessageAttemptFilterFlags(list)
	mac.cmd.AddCommand(list)

	// list destinations
	listDestinations := &cobra.Command{
		Use:   "list-destinations APP_ID MSG_ID",
		Short: "List attempted destinations",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			l, err := s.MessageAttempt.ListAttemptedDestinations(appID, msgID, getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l)
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
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			l, err := s.MessageAttempt.ListAttemptsForEndpoint(appID, msgID, endpointID, *getFilterOptionsMessageAttempt(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l)
			return nil
		},
	}
	addMessageAttemptFilterFlags(listEndpoint)
	mac.cmd.AddCommand(listEndpoint)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID MSG_ID ATTEMPT_ID",
		Short: "Get attempt by id",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			msgID := args[1]
			attemptID := args[2]

			out, err := s.MessageAttempt.Get(appID, msgID, attemptID)
			if err != nil {
				return err
			}

			pretty.Print(out)
			return nil
		},
	}
	mac.cmd.AddCommand(get)

	// resend
	resend := &cobra.Command{
		Use:   "resend APP_ID MSG_ID ENDPOINT_ID",
		Short: "resends a webhook message by id",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			msgID := args[1]
			endpointID := args[2]

			err := s.MessageAttempt.Resend(appID, msgID, endpointID)
			if err != nil {
				return err
			}
			fmt.Println("Message Resent!")
			return nil
		},
	}
	mac.cmd.AddCommand(resend)

	return mac
}

func (mac *messageAttemptCmd) Cmd() *cobra.Command {
	return mac.cmd
}
