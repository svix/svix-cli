package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	svix "github.com/svixhq/svix-libs/go"
)

type applicationCmd struct {
	cmd *cobra.Command
}

func newApplicationCmd() *applicationCmd {

	ac := &applicationCmd{}
	ac.cmd = &cobra.Command{
		Use:   "application",
		Short: "List, create & modify applications",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List current applications",
		RunE: func(cmd *cobra.Command, args []string) error {

			svixClient := getSvixClientOrExit()
			l, err := svixClient.Application.List(getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
		},
	}
	addFilterFlags(list)
	ac.cmd.AddCommand(list)

	// create
	create := &cobra.Command{
		Use:   "create NAME [UID]",
		Short: "Create a new application",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			name := args[0]
			var uid *string
			if len(args) >= 2 {
				uid = &args[1]
			}

			app := &svix.ApplicationIn{
				Name: name,
				Uid:  uid,
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Application.Create(app)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	ac.cmd.AddCommand(create)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID",
		Short: "Get an application by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			appID := args[0]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Application.Get(appID)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	ac.cmd.AddCommand(get)

	update := &cobra.Command{
		Use:   "update APP_ID NAME [UID]",
		Short: "Update an application by id",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			name := args[1]
			var uid *string
			if len(args) >= 2 {
				uid = &args[2]
			}

			app := &svix.ApplicationIn{
				Name: name,
				Uid:  uid,
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Application.Update(appID, app)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	ac.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete APP_ID",
		Short: "Delete an application by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]

			svixClient := getSvixClientOrExit()
			err := svixClient.Application.Delete(appID)
			if err != nil {
				return err
			}

			fmt.Printf("Application \"%s\" Deleted!\n", appID)
			return nil
		},
	}
	ac.cmd.AddCommand(delete)

	return ac
}

func (ac *applicationCmd) Cmd() *cobra.Command {
	return ac.cmd
}
