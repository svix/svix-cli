package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svix/svix-cli/pretty"
	svix "github.com/svix/svix-libs/go"
)

func getPrinterOptions(cmd *cobra.Command) *pretty.PrinterOptions {
	colorFlag := viper.GetBool("color")
	if !colorFlag {
		return nil
	}
	return &pretty.PrinterOptions{
		Color: true,
	}
}

func addFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor id for list call")
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
}

func getFilterOptions(cmd *cobra.Command) *svix.FetchOptions {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.FetchOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if cmd.Flags().Changed("iterator") {
		opts.Iterator = &iteratorFlag
	}
	return opts
}

func addMessageAttemptFilterFlags(cmd *cobra.Command) {
	addFilterFlags(cmd)
	cmd.Flags().StringP("status", "s", "", "message status")
}

func getFilterOptionsMessageAttempt(cmd *cobra.Command) *svix.FetchOptionsMessageAttempt {
	baseOpts := getFilterOptions(cmd)
	opts := &svix.FetchOptionsMessageAttempt{
		FetchOptions: *baseOpts,
	}

	statusFlag, _ := cmd.Flags().GetString("status")
	if cmd.Flags().Changed("status") {
		opts.Iterator = &statusFlag
	}
	return opts
}
