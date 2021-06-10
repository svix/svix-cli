package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/validators"
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
	nameFlagName := "data-name"
	uidFlagName := "data-uid"
	create := &cobra.Command{
		Use:   "create [JSON_PAYLOAD]",
		Short: "Create a new application",
		Long: `Creates a new application

Example Schema:
{
  "uid": "string",
  "name": "string"
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
			var app svix.ApplicationIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &app)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(nameFlagName) {
				nameFlag, err := cmd.Flags().GetString(nameFlagName)
				cobra.CheckErr(err)
				app.Name = nameFlag
			}
			if cmd.Flags().Changed(uidFlagName) {
				uidFlag, err := cmd.Flags().GetString(uidFlagName)
				cobra.CheckErr(err)
				app.Uid = &uidFlag
			}

			// validate args
			if app.Name == "" {
				return fmt.Errorf("name required")
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Application.Create(&app)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	create.Flags().String(nameFlagName, "", "Name of the Application")
	create.Flags().String(uidFlagName, "", "UID of the application (optional)")
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
		Use:   "update APP_ID [JSON_PAYLOAD]",
		Short: "Update an application by id",
		Long: `Update an application by id

Example Schema:
{
  "uid": "string",
  "name": "string"
}
`,
		Args: validators.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse positional args
			appID := args[0]

			var in []byte
			if len(args) > 1 {
				in = []byte(args[1])
			} else {
				var err error
				in, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}
			var app svix.ApplicationIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &app)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(nameFlagName) {
				nameFlag, err := cmd.Flags().GetString(nameFlagName)
				cobra.CheckErr(err)
				app.Name = nameFlag
			}
			if cmd.Flags().Changed(uidFlagName) {
				uidFlag, err := cmd.Flags().GetString(uidFlagName)
				cobra.CheckErr(err)
				app.Uid = &uidFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Application.Update(appID, &app)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	update.Flags().String(nameFlagName, "", "Name of the Application")
	update.Flags().String(uidFlagName, "", "UID of the application (optional)")
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
