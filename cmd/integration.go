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

type integrationCmd struct {
	cmd *cobra.Command
}

func newIntegrationCmd(ctx context.Context) *integrationCmd {

	ic := &integrationCmd{}
	ic.cmd = &cobra.Command{
		Use:   "integration",
		Short: "List, create & modify integrations",
	}

	// list
	list := &cobra.Command{
		Use:   "list APP_ID",
		Short: "List integrations by app id",
		Args:  validators.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			appID := args[0]

			l, err := svixClient.Integration.List(ctx, appID, getIntegrationListOptions(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addIntegrationFilterFlags(list)
	ic.cmd.AddCommand(list)

	// create
	nameFlagName := "data-name"
	create := &cobra.Command{
		Use:   "create APP_ID [JSON_PAYLOAD]",
		Short: "Create a new integration",
		Long: `Create a new integration

Example Schema:
{
    "name": "Example Integration"
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
				in, err = utils.ReadStdin()
				printer.CheckErr(err)
			}
			var integration svix.IntegrationIn
			if len(in) > 0 {
				err := json.Unmarshal(in, &integration)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(nameFlagName) {
				nameFlag, err := cmd.Flags().GetString(nameFlagName)
				printer.CheckErr(err)
				integration.Name = nameFlag
			}

			// validate args
			if integration.Name == "" {
				printer.CheckErr(fmt.Errorf("name required"))
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Integration.Create(ctx, appID, &integration)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	create.Flags().String(nameFlagName, "", "")
	ic.cmd.AddCommand(create)

	// get
	get := &cobra.Command{
		Use:   "get APP_ID INTEGRATION_ID",
		Short: "Get an integration by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			appID := args[0]
			integrationID := args[1]

			l, err := svixClient.Integration.Get(ctx, appID, integrationID)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	ic.cmd.AddCommand(get)

	// update
	update := &cobra.Command{
		Use:   "update APP_ID INTEGRATION_ID [JSON_PAYLOAD]",
		Short: "update an integration by id",
		Long: `update an integration by id

Example Schema:
{
    "name": "Example Integration"
}
`,
		Args: validators.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			// parse positional args
			appID := args[0]
			integrationID := args[1]

			var in []byte
			if len(args) > 2 {
				in = []byte(args[2])
			} else {
				var err error
				in, err = utils.ReadStdin()
				printer.CheckErr(err)
			}
			var integration svix.IntegrationUpdate
			if len(in) > 0 {
				err := json.Unmarshal(in, &integration)
				printer.CheckErr(err)
			}

			// get flags
			if cmd.Flags().Changed(nameFlagName) {
				nameFlag, err := cmd.Flags().GetString(nameFlagName)
				printer.CheckErr(err)
				integration.Name = nameFlag
			}

			// validate args
			if integration.Name == "" {
				printer.CheckErr(fmt.Errorf("name required"))
			}

			svixClient := getSvixClientOrExit()
			out, err := svixClient.Integration.Update(ctx, appID, integrationID, &integration)
			printer.CheckErr(err)

			printer.Print(out)
		},
	}
	update.Flags().String(nameFlagName, "", "")
	ic.cmd.AddCommand(update)

	// delete
	delete := &cobra.Command{
		Use:   "delete APP_ID INTEGRATION_ID",
		Short: "Delete an integration by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			appID := args[0]
			integrationID := args[1]

			err := svixClient.Integration.Delete(ctx, appID, integrationID)
			printer.CheckErr(err)
		},
	}
	ic.cmd.AddCommand(delete)

	// get-key
	getKey := &cobra.Command{
		Use:   "get-key APP_ID INTEGRATION_ID",
		Short: "Get an integration's key by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			appID := args[0]
			integrationID := args[1]

			l, err := svixClient.Integration.GetKey(ctx, appID, integrationID)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	ic.cmd.AddCommand(getKey)

	// rotate-key
	rotateKey := &cobra.Command{
		Use:   "rotate-key APP_ID INTEGRATION_ID",
		Short: "Rotate an integration's key by id",
		Args:  validators.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			appID := args[0]
			integrationID := args[1]

			l, err := svixClient.Integration.RotateKey(ctx, appID, integrationID)
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	ic.cmd.AddCommand(rotateKey)

	return ic
}

func addIntegrationFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor id for list call")
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
}

func getIntegrationListOptions(cmd *cobra.Command) *svix.IntegrationListOptions {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.IntegrationListOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if cmd.Flags().Changed("iterator") {
		opts.Iterator = &iteratorFlag
	}

	return opts
}
