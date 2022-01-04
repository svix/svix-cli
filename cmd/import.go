package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/inout"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
)

type importCmd struct {
	cmd *cobra.Command
}

func newImportCmd() *importCmd {
	forceFlagName := "force"

	cmd := &cobra.Command{
		Use:   "import event-types",
		Short: "Import data to your Svix Organization",
	}

	importEventTypes := &cobra.Command{
		Use:   "event-types [IN_FILE]",
		Short: "Import event-types from a file",
		Long: `imports event-types into your Svix Organization from json or csv

If no IN_FILE path is supplied it, we will read from stdin.

CSV Format:
name,description

Json Format:
[{
	name: "",
	description: ""
}]`,
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))
			svixClient := getSvixClientOrExit()

			force, err := cmd.Flags().GetBool(forceFlagName)
			printer.CheckErr(err)

			var reader io.Reader
			fileName := ""
			if len(args) > 0 {
				fileName = args[0]
				file, err := os.Open(fileName)
				printer.CheckErr(err)
				defer file.Close()
				reader = file
			} else {
				isReadable, err := utils.IsStdinReadable()
				printer.CheckErr(err)
				if !isReadable {
					printer.CheckErr(fmt.Errorf("stdin not readable"))
				}
				reader = os.Stdin
			}

			fileType := getOrInferFileType(fileName)
			switch fileType {
			case "csv":
				err := inout.ImportEventTypesCsv(svixClient, reader, force)
				printer.CheckErr(err)
			default:
				err := inout.ImportEventTypesJson(svixClient, reader, force)
				printer.CheckErr(err)
			}
		},
	}
	importEventTypes.Flags().AddGoFlag(flag.Lookup(fileTypeFlagName))
	importEventTypes.Flags().Bool(forceFlagName, false, "Update event type if already exists (defaults to skipping)")
	cmd.AddCommand(importEventTypes)

	return &importCmd{
		cmd: cmd,
	}
}
