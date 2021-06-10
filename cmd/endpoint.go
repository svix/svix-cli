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

type endpointCmd struct {
	cmd *cobra.Command
}

func newEndpointCmd() *endpointCmd {
	ec := &endpointCmd{}
	ec.cmd = &cobra.Command{
		Use:   "endpoint",
		Short: "List, create & modify endpoints",
	}

	// list
	list := &cobra.Command{
		Use:   "list APP_ID",
		Short: "List current endpoints",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.Endpoint.List(appID, getFilterOptions(cmd))
			if err != nil {
				return err
			}

			pretty.Print(l, getPrintOptions(cmd))
			return nil
		},
	}
	addFilterFlags(list)
	ec.cmd.AddCommand(list)

	// create
	urlFlagName := "data-url"
	versionFlagName := "data-version"
	filterTypesFlagName := "data-filterTypes"
	create := &cobra.Command{
		Use:   "create APP_ID [JSON_PAYLOAD]",
		Short: "Create a new endpoint",
		Long: `Create a new endpoint

Example Schema:
{
	"url": "string",
	"version": 0,
	"description": "",
	"filterTypes": [
	  "string"
	]
  }
`,
		Args: validators.RangeArgs(1, 2),
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
			var ep svix.EndpointIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &ep)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(urlFlagName) {
				urlFlag, err := cmd.Flags().GetString(urlFlagName)
				cobra.CheckErr(err)
				ep.Url = urlFlag
			}
			if cmd.Flags().Changed(versionFlagName) {
				versionFlag, err := cmd.Flags().GetInt32(versionFlagName)
				cobra.CheckErr(err)
				ep.Version = versionFlag
			}
			if cmd.Flags().Changed(filterTypesFlagName) {
				filterTypesFlag, err := cmd.Flags().GetStringArray(filterTypesFlagName)
				cobra.CheckErr(err)
				ep.FilterTypes = &filterTypesFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Create(appID, &ep)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	create.Flags().String(urlFlagName, "", "")
	create.Flags().Int32(versionFlagName, 0, "")
	create.Flags().StringArray(filterTypesFlagName, []string{}, "")
	ec.cmd.AddCommand(create)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID ENDPOINT_ID",
		Short: "Get an endpoint by id",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appID := args[0]
			endpointID := args[1]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Get(appID, endpointID)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	ec.cmd.AddCommand(get)

	update := &cobra.Command{
		Use:   "update APP_ID ENDPOINT_ID [JSON_PAYLOAD]",
		Short: "Update an application by id",
		Long: `Update an application by id

Example Schema:
{
  "url": "string",
  "version": 0,
  "description": "",
  "filterTypes": [
    "string"
  ]
}
`,
		Args: validators.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			endpointID := args[1]

			var in []byte
			if len(args) > 2 {
				in = []byte(args[2])
			} else {
				var err error
				in, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}
			var ep svix.EndpointIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &ep)
				cobra.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(urlFlagName) {
				urlFlag, err := cmd.Flags().GetString(urlFlagName)
				cobra.CheckErr(err)
				ep.Url = urlFlag
			}
			if cmd.Flags().Changed(versionFlagName) {
				versionFlag, err := cmd.Flags().GetInt32(versionFlagName)
				cobra.CheckErr(err)
				ep.Version = versionFlag
			}
			if cmd.Flags().Changed(filterTypesFlagName) {
				filterTypesFlag, err := cmd.Flags().GetStringArray(filterTypesFlagName)
				cobra.CheckErr(err)
				ep.FilterTypes = &filterTypesFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Update(appID, endpointID, &ep)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	update.Flags().String(urlFlagName, "", "")
	update.Flags().Int32(versionFlagName, 0, "")
	update.Flags().StringArray(filterTypesFlagName, []string{}, "")
	ec.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete APP_ID ENDPOINT_ID",
		Short: "Delete an endpoint by id",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			endpointID := args[1]

			utils.Confirm(fmt.Sprintf("Are you sure you want to delete the the endpoint with id: %s", endpointID))

			svixClient := getSvixClientOrExit()
			err := svixClient.Endpoint.Delete(appID, endpointID)
			if err != nil {
				return err
			}

			fmt.Printf("Endpoint \"%s\" Deleted!\n", endpointID)
			return nil
		},
	}
	ec.cmd.AddCommand(delete)

	secret := &cobra.Command{
		Use:   "secret APP_ID ENDPOINT_ID",
		Short: "get an endpoint's secret by id",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			endpointID := args[1]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.GetSecret(appID, endpointID)
			if err != nil {
				return err
			}

			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	ec.cmd.AddCommand(secret)

	return ec
}
