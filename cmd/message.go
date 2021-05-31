package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	"github.com/svixhq/svix-cli/validators"
	svix "github.com/svixhq/svix-libs/go"
)

type messageCmd struct {
	cmd *cobra.Command
}

func newMessageCmd() *messageCmd {
	mc := &messageCmd{}
	mc.cmd = &cobra.Command{
		Use:   "message",
		Short: "List & create messages",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List current messages",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.Message.List(appID, getFilterOptions(cmd))
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
		Args:  validators.RangeArgs(2, 3),
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
				Payload:   payload,
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Message.Create(appID, msg)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	mc.cmd.AddCommand(create)

	get := &cobra.Command{
		Use:   "get APP_ID MSG_ID",
		Short: "get message by id",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			msgID := args[1]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Message.Get(appID, msgID)
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
