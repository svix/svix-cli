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

type eventTypeCmd struct {
	cmd *cobra.Command
}

func newEventTypeCmd() *eventTypeCmd {
	etc := &eventTypeCmd{}
	etc.cmd = &cobra.Command{
		Use:   "event-type",
		Short: "List, create & modify event types",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List current event types",
		RunE: func(cmd *cobra.Command, args []string) error {

			svixClient := getSvixClientOrExit()
			l, err := svixClient.EventType.List(getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
		},
	}
	addFilterFlags(list)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var in []byte
			if len(args) > 0 {
				in = []byte(args[0])
			} else {
				var err error
				in, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}
			var et svix.EventTypeIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &et)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(nameFlagName) {
				nameFlag, err := cmd.Flags().GetString(nameFlagName)
				cobra.CheckErr(err)
				et.Name = nameFlag
			}
			if cmd.Flags().Changed(descriptionFlagName) {
				descFlag, err := cmd.Flags().GetString(descriptionFlagName)
				cobra.CheckErr(err)
				et.Description = descFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Create(&et)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
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
		Args: validators.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get poisitonal args
			eventName := args[0]

			var in []byte
			if len(args) > 1 {
				in = []byte(args[1])
			} else {
				var err error
				in, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}
			var et svix.EventTypeUpdate
			if len(in) > 0 {
				err := json.Unmarshal(in, &et)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(descriptionFlagName) {
				descFlag, err := cmd.Flags().GetString(descriptionFlagName)
				cobra.CheckErr(err)
				et.Description = descFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Update(eventName, &et)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	update.Flags().String(descriptionFlagName, "", "")
	etc.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete EVENT_ID",
		Short: "Delete an event type by id",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			eventID := args[0]

			utils.Confirm(fmt.Sprintf("Are you sure you want to delete the the event with id: %s", eventID))

			svixClient := getSvixClientOrExit()
			err := svixClient.EventType.Delete(eventID)
			if err != nil {
				return err
			}

			fmt.Printf("Event Type \"%s\" Deleted!\n", eventID)
			return nil
		},
	}
	etc.cmd.AddCommand(delete)

	return etc
}
