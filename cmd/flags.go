package cmd

import (
	"flag"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svix/svix-cli/flags"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/utils"
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
		if runtime.GOOS != "windows" {
			isTTY, _, err := utils.IsTTY(os.Stdout)
			if err == nil {
				// just defaults to false if an error occurs
				color = isTTY
			}
		}
	}

	return &pretty.PrinterOptions{
		Color: color,
	}
}
