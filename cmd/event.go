package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
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

	// TODO add function once svix-libs v0.14.0 is released
	// // list
	// list := &cobra.Command{
	// 	Use:   "list",
	// 	Short: "List current event types",
	// 	RunE: func(cmd *cobra.Command, args []string) error {

	// 		svixClient := getSvixClientOrExit()
	// 		l, err :=svixClient.EventType.List(getFilterOptions(cmd))
	// 		if err != nil {
	// 			return err
	// 		}

	// 		pretty.Print(l, getPrintOptions(cmd))
	// 		return nil
	// 	},
	// }
	// addFilterFlags(list)
	// etc.cmd.AddCommand(list)

	// create
	create := &cobra.Command{
		Use:   "create NAME DESCRIPTION",
		Short: "Create a new event type",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			et := &svix.EventTypeInOut{
				Name:        args[0],
				Description: args[1],
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
	etc.cmd.AddCommand(create)

	update := &cobra.Command{
		Use:   "update EVENT_ID DESCRIPTION",
		Short: "Update an event type by id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID := args[0]
			et := &svix.EventTypeUpdate{
				Description: args[1],
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.EventType.Update(eventID, et)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	etc.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete EVENT_ID",
		Short: "Delete an event type by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			eventID := args[0]

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

func (etc *eventTypeCmd) Cmd() *cobra.Command {
	return etc.cmd
}
