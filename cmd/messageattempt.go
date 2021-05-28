package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/cfg"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type messageAttemptCmd struct {
	cmd *cobra.Command
	cfg *cfg.Config
	sc  *svix.Svix
}

func newMessageAttemptCmd(cfg *cfg.Config, s *svix.Svix) *messageAttemptCmd {
	mac := &messageAttemptCmd{
		sc:  s,
		cfg: cfg,
	}
	mac.cmd = &cobra.Command{
		Use:   "message-attempt",
		Short: "List & create messages",
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

			pretty.PrintListResponseMessageAttemptOut(l)
			return nil
		},
	}
	addMessageAttemptFilterFlags(list)
	mac.cmd.AddCommand(list)

	// list destinations
	listDestinations := &cobra.Command{
		Use:   "list-destinations APP_ID MSG_ID",
		Short: "List attempted messages",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			l, err := s.MessageAttempt.ListAttemptedDestinations(appID, msgID, getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.PrintListResponseMessageEndpointOut(l)
			return nil
		},
	}
	addFilterFlags(listDestinations)
	mac.cmd.AddCommand(listDestinations)

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
