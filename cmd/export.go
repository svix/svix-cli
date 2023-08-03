package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"io"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/inout"
	"github.com/svix/svix-cli/pretty"
)

type exportCmd struct {
	cmd *cobra.Command
}

func newExportCmd() *exportCmd {
	cmd := &cobra.Command{
		Use:   "export event-types",
		Short: "Export data from your Svix Organization",
	}

	exportEventTypes := &cobra.Command{
		Use:   "event-types [OUT_FILE]",
		Short: "Export event-types to a file",
		Long: `exports event-types from your Svix Organization to a .csv or .json file

If no OUT_FILE path is supplied it, output to stdout.

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

			eventTypes, err := inout.GetAllEventTypes(context.Background(), svixClient)
			printer.CheckErr(err)

			var outStream io.Writer = printer
			fileName := ""
			if len(args) > 0 {
				fileName = args[0]
				outFile, err := inout.CreateOrTruncateFile(fileName)
				printer.CheckErr(err)
				defer outFile.Close()
				outStream = outFile
			}
			fileType := getOrInferFileType(fileName)
			switch fileType {
			case "csv":
				err := inout.WriteEventTypesAsCsv(eventTypes, outStream)
				printer.CheckErr(err)
			default:
				enc := json.NewEncoder(outStream)
				err := enc.Encode(eventTypes)
				printer.CheckErr(err)
			}
		},
	}
	exportEventTypes.Flags().AddGoFlag(flag.Lookup(fileTypeFlagName))
	cmd.AddCommand(exportEventTypes)

	return &exportCmd{
		cmd: cmd,
	}
}
