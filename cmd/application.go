package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	"github.com/svixhq/svix-cli/utils"
	"github.com/svixhq/svix-cli/validators"
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
		Use:   "create",
		Short: "Create a new application",
		Args:  validators.NoArgs(),
		RunE: func(cmd *cobra.Command, args []string) error {
			app := &svix.ApplicationIn{}
			err := utils.TryMarshallPipe(app)
			cobra.CheckErr(err)

			// get flags
			if cmd.Flags().Changed("name") {
				nameFlag, err := cmd.Flags().GetString("name")
				cobra.CheckErr(err)
				app.Name = nameFlag
			}
			if cmd.Flags().Changed("uid") {
				uidFlag, err := cmd.Flags().GetString("uid")
				cobra.CheckErr(err)
				app.Uid = &uidFlag
			}

			// validate args
			if app.Name == "" {
				return fmt.Errorf("name required")
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
	create.Flags().String("name", "", "Name of the Application")
	create.Flags().String("uid", "", "UID of the application (optional)")
	ac.cmd.AddCommand(create)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID",
		Short: "Get an application by id",
		Args:  validators.ExactArgs(1),
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
		Use:   "update APP_ID",
		Short: "Update an application by id",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse positional args
			appID := args[0]

			app := &svix.ApplicationIn{}
			err := utils.TryMarshallPipe(app)
			cobra.CheckErr(err)

			// get flags
			if cmd.Flags().Changed("name") {
				nameFlag, err := cmd.Flags().GetString("name")
				cobra.CheckErr(err)
				app.Name = nameFlag
			}
			if cmd.Flags().Changed("uid") {
				uidFlag, err := cmd.Flags().GetString("uid")
				cobra.CheckErr(err)
				app.Uid = &uidFlag
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
	update.Flags().String("name", "", "Name of the Application")
	update.Flags().String("uid", "", "UID of the application (optional)")
	ac.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete APP_ID",
		Short: "Delete an application by id",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]

			svixClient := getSvixClientOrExit()

			utils.Confirm(fmt.Sprintf("Are you sure you want to delete the app with id: %s", appID))

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
