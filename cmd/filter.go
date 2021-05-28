package cmd

import (
	"github.com/spf13/cobra"
	svix "github.com/svixhq/svix-libs/go"
)

func addFilterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("iterator", "i", "", "anchor for call TODO better desc")

	// TODO cobra has no ability to set true pointer flags, always requiring a default value
	// decide on a way to check if this flag has been explicitly set or not (or we can just leave a sane default max)
	cmd.Flags().Int32P("limit", "l", 50, "max items per request")
}

func getFilterOptions(cmd *cobra.Command) *svix.FetchOptions {
	limit, _ := cmd.Flags().GetInt32("limit")

	opts := &svix.FetchOptions{
		Limit: &limit,
	}

	iteratorFlag, _ := cmd.Flags().GetString("iterator")
	if iteratorFlag != "" {
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
	if statusFlag != "" {
		opts.Iterator = &statusFlag
	}
	return opts
}
