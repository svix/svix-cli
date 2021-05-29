package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type messageCmd struct {
	cmd *cobra.Command
	sc  *svix.Svix
}

func newMessageCmd(s *svix.Svix) *messageCmd {
	mc := &messageCmd{
		sc: s,
	}
	mc.cmd = &cobra.Command{
		Use:   "message",
		Short: "List & create messages",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List current messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			l, err := s.Message.List(appID, getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
		},
	}
	addFilterFlags(list)
	mc.cmd.AddCommand(list)

	// create
	create := &cobra.Command{
		Use:   "create APP_ID EVENT_TYPE [EVENT_ID] JSON_PAYLOAD",
		Short: "Create a new messsage",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			eventType := args[1]
			var eventID string
			var payloadStr string
			if len(args) > 3 {
				eventID = args[2]
				payloadStr = args[3]
			} else {
				payloadStr = args[2]
			}

			// unmarshal payload
			var payload map[string]interface{}
			err := json.Unmarshal([]byte(payloadStr), &payload)
			if err != nil {
				return fmt.Errorf("invalid payload json")
			}

			msg := &svix.MessageIn{
				EventType: eventType,
				EventId:   &eventID,
				Data:      payload,
			}

			out, err := s.Message.Create(appID, msg)
			if err != nil {
				return err
			}
			fmt.Println("Message Created!")
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	mc.cmd.AddCommand(create)

	get := &cobra.Command{
		Use:   "get APP_ID MSG_ID",
		Short: "get message by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			out, err := s.Message.Get(appID, msgID)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	mc.cmd.AddCommand(get)

	return mc
}

func (mc *messageCmd) Cmd() *cobra.Command {
	return mc.cmd
}
