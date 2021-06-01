package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	"github.com/svixhq/svix-cli/utils"
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
		Use:   "create APP_ID",
		Short: "Create a new messsage",
		Args:  validators.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get positional args
			appID := args[0]

			msg := &svix.MessageIn{}
			err := utils.TryMarshallPipe(msg)
			cobra.CheckErr(err)

			// get flags
			eventTypeFlag, err := cmd.Flags().GetString("eventType")
			cobra.CheckErr(err)
			if eventTypeFlag != "" {
				msg.EventType = eventTypeFlag
			}
			eventIdFlag, err := cmd.Flags().GetString("eventId")
			cobra.CheckErr(err)
			if eventIdFlag != "" {
				msg.EventId = &eventIdFlag
			}
			payloadFlag, err := cmd.Flags().GetString("payload")
			cobra.CheckErr(err)
			if payloadFlag != "" {
				// unmarshal payload
				var payload map[string]interface{}
				err := json.Unmarshal([]byte(payloadFlag), &payload)
				if err != nil {
					return fmt.Errorf("invalid payload json supplied via flag")
				}
				msg.Payload = payload
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
	create.Flags().String("eventType", "", "")
	create.Flags().String("eventId", "", "")
	create.Flags().String("payload", "", "json message payload")
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
