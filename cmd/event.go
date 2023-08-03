package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-webhooks/go"
)

type eventTypeCmd struct {
	cmd *cobra.Command
}

func newEventTypeCmd(ctx context.Context) *eventTypeCmd {
	etc := &eventTypeCmd{}
	etc.cmd = &cobra.Command{
		Use:   "event-type",
		Short: "List, create & modify event types",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List current event types",
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			svixClient := getSvixClientOrExit()
			l, err := svixClient.EventType.List(ctx, getEventListOptions(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addEventFilterFlags(list)
	etc.cmd.AddCommand(list)

	// create
	nameFlagName := "data-name"
	descriptionFlagName := "data-description"
	create := &cobra.Command{
		Use:   "create [JSON_PAYLOAD]",
		Short: "Create a new event type",
		Long: `Create a new event type

Example Schema:
{
  "description": "string",
  "name": "user.signup"
}
`,
		Args: validators.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			var in []byte
			if len(args) > 0 {
				in = []byte(args[0])
			} else {
				var err error
				in, err = utils.ReadStdin()
				printer.CheckErr(err)
			}
			var et svix.EventTypeIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &et)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(nameFlagName) {
				nameFlag, err := cmd.Flags().GetString(nameFlagName)
				printer.CheckErr(err)
				et.Name = nameFlag
			}
			if cmd.Flags().Changed(descriptionFlagName) {
				descFlag, err := cmd.Flags().GetString(descriptionFlagName)
				printer.CheckErr(err)
				et.Description = descFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Create(ctx, &et)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	create.Flags().String(nameFlagName, "", "")
	create.Flags().String(descriptionFlagName, "", "")
	etc.cmd.AddCommand(create)

	update := &cobra.Command{
		Use:   "update EVENT_TYPE_NAME [JSON_PAYLOAD]",
		Short: "Update an event type by name",
		Long: `Update an event type by name

Example Schema:
{
  "description": "string"
}
	`,
		Args: validators.RangeArgs(1, 2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// get poisitonal args
			eventName := args[0]

			var in []byte
			if len(args) > 1 {
				in = []byte(args[1])
			} else {
				var err error
				in, err = utils.ReadStdin()
				printer.CheckErr(err)
			}
			var et svix.EventTypeUpdate
			if len(in) > 0 {
				err := json.Unmarshal(in, &et)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(descriptionFlagName) {
				descFlag, err := cmd.Flags().GetString(descriptionFlagName)
				printer.CheckErr(err)
				et.Description = descFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Update(ctx, eventName, &et)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	update.Flags().String(descriptionFlagName, "", "")
	etc.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete EVENT_ID",
		Short: "Delete an event type by id",
		Args:  validators.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse args
			eventID := args[0]

			utils.Confirm(fmt.Sprintf("Are you sure you want to delete the the event with id: %s", eventID))

			svixClient := getSvixClientOrExit()
			err := svixClient.EventType.Delete(ctx, eventID)
			printer.CheckErr(err)

			fmt.Printf("Event Type \"%s\" Deleted!\n", eventID)
		},
	}
	etc.cmd.AddCommand(delete)

	return etc
}

func addEventFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor id for list call")
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
	cmd.Flags().Bool("with-content", false, "includes content like schemas")
	cmd.Flags().Bool("include-archived", false, "include archived event types")
}

func getEventListOptions(cmd *cobra.Command) *svix.EventTypeListOptions {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.EventTypeListOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if cmd.Flags().Changed("iterator") {
		opts.Iterator = &iteratorFlag
	}

	withContentFlag, _ := cmd.Flags().GetBool("with-content")
	if cmd.Flags().Changed("with-content") {
		opts.WithContent = &withContentFlag
	}

	includeArchivedFlag, _ := cmd.Flags().GetBool("include-archived")
	if cmd.Flags().Changed("include-archived") {
		opts.IncludeArchived = &includeArchivedFlag
	}

	return opts
}
