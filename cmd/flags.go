package cmd

import (
	"flag"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svix/svix-cli/flags"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
	svix "github.com/svix/svix-libs/go"
)

var fileTypeFlagName = "type"
var fileTypeFlagValue string = "auto"
var fileTypeFlag = flags.NewEnum(&fileTypeFlagValue, "auto", "json", "csv")

func init() {
	flag.Var(fileTypeFlag, fileTypeFlagName, "auto|json|csv")
}

func getOrInferFileType(fileName string) string {
	fileType := fileTypeFlagValue
	if fileTypeFlagValue == "auto" {
		switch {
		case strings.HasSuffix(fileName, ".csv"):
			fileType = "csv"
		default:
			fileType = "json"
		}
	}
	return fileType
}

func getPrinterOptions(cmd *cobra.Command) *pretty.PrinterOptions {
	colorFlag := viper.GetString("color")
	color := false
	switch colorFlag {
	case "always":
		color = true
	case "never":
		color = false
	default:
		isTTY, _, err := utils.IsTTY(os.Stdout)
		if err == nil {
			// just defaults to false if an error occurs
			color = isTTY
		}
	}

	return &pretty.PrinterOptions{
		Color: color,
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
