package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/pretty"
	"github.com/svixhq/svix-cli/utils"
	"github.com/svixhq/svix-cli/validators"
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

	create := &cobra.Command{
		Use:   "create APP_ID",
		Short: "Create a new endpoint",
		Args:  validators.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse positional args
			appID := args[0]

			ep := &svix.EndpointIn{}
			err := utils.TryMarshallPipe(ep)
			cobra.CheckErr(err)

			// get flags
			urlFlag, err := cmd.Flags().GetString("url")
			cobra.CheckErr(err)
			if urlFlag != "" {
				ep.Url = urlFlag
			}
			versionFlag, err := cmd.Flags().GetInt32("version")
			cobra.CheckErr(err)
			if versionFlag != 0 {
				ep.Version = versionFlag
			}
			filterTypesFlag, err := cmd.Flags().GetStringArray("filterTypes")
			cobra.CheckErr(err)
			if len(filterTypesFlag) > 0 {
				ep.FilterTypes = &filterTypesFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Create(appID, ep)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	create.Flags().String("url", "", "")
	create.Flags().Int32("version", 0, "")
	create.Flags().StringArray("filterTypes", []string{}, "")
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
		Use:   "update APP_ID ENDPOINT_ID",
		Short: "Update an application by id",
		Args:  validators.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			appID := args[0]
			endpointID := args[1]

			ep := &svix.EndpointIn{}
			err := utils.TryMarshallPipe(ep)
			cobra.CheckErr(err)

			// get flags
			urlFlag, err := cmd.Flags().GetString("url")
			cobra.CheckErr(err)
			if urlFlag != "" {
				ep.Url = urlFlag
			}
			versionFlag, err := cmd.Flags().GetInt32("version")
			cobra.CheckErr(err)
			if versionFlag != 0 {
				ep.Version = versionFlag
			}
			filterTypesFlag, err := cmd.Flags().GetStringArray("filterTypes")
			cobra.CheckErr(err)
			if len(filterTypesFlag) > 0 {
				ep.FilterTypes = &filterTypesFlag
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Endpoint.Update(appID, endpointID, ep)
			if err != nil {
				return err
			}
			pretty.Print(out, getPrintOptions(cmd))
			return nil
		},
	}
	update.Flags().String("url", "", "")
	update.Flags().Int32("version", 0, "")
	update.Flags().StringArray("filterTypes", []string{}, "")
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
