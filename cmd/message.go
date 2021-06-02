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
		Use:   "list APP_ID",
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
	eventTypeFlagName := "data-eventType"
	eventIdFlagName := "data-eventId"
	payloadFlagName := "data-payload"
	create := &cobra.Command{
		Use:   "create APP_ID [JSON_PAYLOAD]",
		Short: "Create a new messsage",
		Long: `Create a new messsage

Example Schema:
{
  "eventType": "user.signup",
  "eventId": "evt_pNZKtWg8Azow",
  "payload": {
    "username": "test_user",
    "email": "test@example.com"
  }
}
`,
		Args: validators.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get positional args
			appID := args[0]
			var in []byte
			if len(args) > 1 {
				in = []byte(args[1])
			} else {
				var err error
				in, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}
			var msg svix.MessageIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &msg)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(eventTypeFlagName) {
				eventTypeFlag, err := cmd.Flags().GetString(eventTypeFlagName)
				cobra.CheckErr(err)
				msg.EventType = eventTypeFlag
			}
			if cmd.Flags().Changed(eventIdFlagName) {
				eventIdFlag, err := cmd.Flags().GetString(eventIdFlagName)
				cobra.CheckErr(err)
				msg.EventId = &eventIdFlag
			}
			if cmd.Flags().Changed(payloadFlagName) {
				payloadFlag, err := cmd.Flags().GetString(payloadFlagName)
				cobra.CheckErr(err)
				// unmarshal payload
				var payload map[string]interface{}
				err = json.Unmarshal([]byte(payloadFlag), &payload)
				if err != nil {
					return fmt.Errorf("invalid payload json supplied via flag")
				}
				msg.Payload = payload
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Message.Create(appID, &msg)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	create.Flags().String(eventTypeFlagName, "", "")
	create.Flags().String(eventIdFlagName, "", "")
	create.Flags().String(payloadFlagName, "", "json message payload")
	mc.cmd.AddCommand(create)

	get := &cobra.Command{
		Use:   "get APP_ID MSG_ID",
		Short: "get message by id",
		Args:  validators.ExactArgs(2),
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
