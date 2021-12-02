package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-libs/go"
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
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]

			svixClient := getSvixClientOrExit()
			l, err := svixClient.Endpoint.List(appID, getEndpointListOptions(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addEndpointFilterFlags(list)
	ec.cmd.AddCommand(list)

	// create
	urlFlagName := "data-url"
	versionFlagName := "data-version"
	filterTypesFlagName := "data-filterTypes"
	rateLimitFlagName := "data-rate-limit"
	disabledFlagName := "data-disabled"
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
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse positional args
			appID := args[0]

			var in []byte
			if len(args) > 1 {
				in = []byte(args[1])
			} else {
				var err error
				in, err = utils.ReadPipe()
				printer.CheckErr(err)
			}
			var ep svix.EndpointIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &ep)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(urlFlagName) {
				urlFlag, err := cmd.Flags().GetString(urlFlagName)
				printer.CheckErr(err)
				ep.Url = urlFlag
			}
			if cmd.Flags().Changed(versionFlagName) {
				versionFlag, err := cmd.Flags().GetInt32(versionFlagName)
				printer.CheckErr(err)
				ep.Version = versionFlag
			}
			if cmd.Flags().Changed(filterTypesFlagName) {
				filterTypesFlag, err := cmd.Flags().GetStringArray(filterTypesFlagName)
				printer.CheckErr(err)
				ep.FilterTypes = &filterTypesFlag
			}
			if cmd.Flags().Changed(rateLimitFlagName) {
				rateLimitFlag, err := cmd.Flags().GetInt32(rateLimitFlagName)
				printer.CheckErr(err)
				ep.RateLimit = &rateLimitFlag
			}
			if cmd.Flags().Changed(disabledFlagName) {
				disabledFlag, err := cmd.Flags().GetBool(disabledFlagName)
				printer.CheckErr(err)
				ep.Disabled = &disabledFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Create(appID, &ep)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	create.Flags().String(urlFlagName, "", "")
	create.Flags().Int32(versionFlagName, 0, "")
	create.Flags().StringArray(filterTypesFlagName, []string{}, "")
	create.Flags().Int32(rateLimitFlagName, 0, "Rate Limit of the endpoint (optional)")
	create.Flags().Bool(disabledFlagName, false, "")
	ec.cmd.AddCommand(create)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID ENDPOINT_ID",
		Short: "Get an endpoint by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			appID := args[0]
			endpointID := args[1]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Get(appID, endpointID)
			printer.CheckErr(err)

			printer.Print(out)
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
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse args
			appID := args[0]
			endpointID := args[1]

			var in []byte
			if len(args) > 2 {
				in = []byte(args[2])
			} else {
				var err error
				in, err = utils.ReadPipe()
				printer.CheckErr(err)
			}
			var ep svix.EndpointUpdate
			if len(in) > 0 {
				err := json.Unmarshal(in, &ep)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(urlFlagName) {
				urlFlag, err := cmd.Flags().GetString(urlFlagName)
				printer.CheckErr(err)
				ep.Url = urlFlag
			}
			if cmd.Flags().Changed(versionFlagName) {
				versionFlag, err := cmd.Flags().GetInt32(versionFlagName)
				printer.CheckErr(err)
				ep.Version = versionFlag
			}
			if cmd.Flags().Changed(filterTypesFlagName) {
				filterTypesFlag, err := cmd.Flags().GetStringArray(filterTypesFlagName)
				printer.CheckErr(err)
				ep.FilterTypes = &filterTypesFlag
			}
			if cmd.Flags().Changed(rateLimitFlagName) {
				rateLimitFlag, err := cmd.Flags().GetInt32(rateLimitFlagName)
				printer.CheckErr(err)
				ep.RateLimit = &rateLimitFlag
			}
			if cmd.Flags().Changed(disabledFlagName) {
				disabledFlag, err := cmd.Flags().GetBool(disabledFlagName)
				printer.CheckErr(err)
				ep.Disabled = &disabledFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Update(appID, endpointID, &ep)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	update.Flags().String(urlFlagName, "", "")
	update.Flags().Int32(versionFlagName, 0, "")
	update.Flags().StringArray(filterTypesFlagName, []string{}, "")
	update.Flags().Int32(rateLimitFlagName, 0, "Rate Limit of the endpoint (optional)")
	update.Flags().Bool(disabledFlagName, false, "")
	ec.cmd.AddCommand(update)

	delete := &cobra.Command{
		Use:   "delete APP_ID ENDPOINT_ID",
		Short: "Delete an endpoint by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse args
			appID := args[0]
			endpointID := args[1]

			utils.Confirm(fmt.Sprintf("Are you sure you want to delete the the endpoint with id: %s", endpointID))

			svixClient := getSvixClientOrExit()
			err := svixClient.Endpoint.Delete(appID, endpointID)
			printer.CheckErr(err)

			fmt.Printf("Endpoint \"%s\" Deleted!\n", endpointID)
		},
	}
	ec.cmd.AddCommand(delete)

	secret := &cobra.Command{
		Use:   "secret APP_ID ENDPOINT_ID",
		Short: "get an endpoint's secret by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse args
			appID := args[0]
			endpointID := args[1]

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.GetSecret(appID, endpointID)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	ec.cmd.AddCommand(secret)

	return ec
}

func addEndpointFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor id for list call")
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
}

func getEndpointListOptions(cmd *cobra.Command) *svix.EndpointListOptions {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.EndpointListOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if cmd.Flags().Changed("iterator") {
		opts.Iterator = &iteratorFlag
	}

	return opts
}
