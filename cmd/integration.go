package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-libs/go"
)

type integrationCmd struct {
	cmd *cobra.Command
}

func newIntegrationCmd() *integrationCmd {

	ac := &integrationCmd{}
	ac.cmd = &cobra.Command{
		Use:   "integration",
		Short: "List, create & modify integrations",
	}

	// list
	list := &cobra.Command{
		Use:   "list",
		Short: "List integrations by app id",
		Args:  validators.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			appID := args[0]

			l, err := svixClient.Integration.List(appID, getIntegrationListOptions(cmd))
			printer.CheckErr(err)

			printer.Print(l)
		},
	}
	addIntegrationFilterFlags(list)
	ac.cmd.AddCommand(list)

	return ac
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
