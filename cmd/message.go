package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/araddon/dateparse"
	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-webhooks/go"
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
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]

			svixClient := getSvixClientOrExit()

			opts, err := getMessageFilterFlags(cmd)
			printer.CheckErr(err)
			l, err := svixClient.Message.List(cmd.Context(), appID, opts)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addMessageFilterFlags(list)
	mc.cmd.AddCommand(list)

	// create
	eventTypeFlagName := "data-eventType"
	eventIdFlagName := "data-eventId"
	payloadFlagName := "data-payload"
	create := &cobra.Command{
		Use:   "create APP_ID [JSON_PAYLOAD]",
		Short: "Create a new message",
		Long: `Create a new message

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
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// get positional args
			appID := args[0]
			var in []byte
			if len(args) > 1 {
				in = []byte(args[1])
			} else {
				var err error
				in, err = utils.ReadStdin()
				printer.CheckErr(err)
			}
			var msg svix.MessageIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &msg)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(eventTypeFlagName) {
				eventTypeFlag, err := cmd.Flags().GetString(eventTypeFlagName)
				printer.CheckErr(err)
				msg.EventType = eventTypeFlag
			}
			if cmd.Flags().Changed(eventIdFlagName) {
				eventIdFlag, err := cmd.Flags().GetString(eventIdFlagName)
				printer.CheckErr(err)
				msg.EventId.Set(&eventIdFlag)
			}
			if cmd.Flags().Changed(payloadFlagName) {
				payloadFlag, err := cmd.Flags().GetString(payloadFlagName)
				printer.CheckErr(err)
				// unmarshal payload
				var payload map[string]interface{}
				err = json.Unmarshal([]byte(payloadFlag), &payload)
				if err != nil {
					printer.CheckErr(fmt.Errorf("invalid payload json supplied via flag"))
				}
				msg.Payload = payload
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Message.Create(cmd.Context(), appID, &msg)
			printer.CheckErr(err)

			printer.Print(out)
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
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			msgID := args[1]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Message.Get(cmd.Context(), appID, msgID)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	mc.cmd.AddCommand(get)

	return mc
}

func addMessageFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor id for list call")
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
	cmd.Flags().StringArray("event-types", []string{}, "event types")
	cmd.Flags().StringP("before", "b", "", "before")
}

func getMessageFilterFlags(cmd *cobra.Command) (*svix.MessageListOptions, error) {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.MessageListOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if cmd.Flags().Changed("iterator") {
		opts.Iterator = &iteratorFlag
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
