package cmd

import (
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
	create := &cobra.Command{
		Use:   "create",
		Short: "Create a new event type",
		Args:  validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			et := &svix.EventTypeInOut{}
			err := utils.TryMarshallPipe(et)
			cobra.CheckErr(err)

			// get flags
			if cmd.Flags().Changed("name") {
				nameFlag, err := cmd.Flags().GetString("name")
				cobra.CheckErr(err)
				et.Name = nameFlag
			}
			if cmd.Flags().Changed("description") {
				descFlag, err := cmd.Flags().GetString("description")
				cobra.CheckErr(err)
				et.Description = descFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Create(et)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	create.Flags().String("name", "", "")
	create.Flags().String("description", "", "")
	etc.cmd.AddCommand(create)

	update := &cobra.Command{
		Use:   "update EVENT_TYPE_NAME",
		Short: "Update an event type by name",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// get poisitonal args
			eventName := args[0]

			et := &svix.EventTypeUpdate{}
			err := utils.TryMarshallPipe(et)
			cobra.CheckErr(err)

			// get flags
			if cmd.Flags().Changed("description") {
				descFlag, err := cmd.Flags().GetString("description")
				cobra.CheckErr(err)
				et.Description = descFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Update(eventName, et)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	update.Flags().String("description", "", "")
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
