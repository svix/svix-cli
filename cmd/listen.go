package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svix/svix-cli/config"
	"github.com/svix/svix-cli/pretty"
	"github.com/svix/svix-cli/relay"
)

type listenCmd struct {
	cmd *cobra.Command
}

func newListenCmd() *listenCmd {
	noLoggingFlagName := "no-logging"
	lc := &listenCmd{}
	lc.cmd = &cobra.Command{
		Use:   `listen localURL (ex. http://localhost:8000/webhook/)`,
		Short: "Forward webhook requests a local url",
		Long: `listen creates an on-the-fly publicly accessible URL for use when testing webhooks.

The cli then acts as a proxy forwarding any requests to the given local URL.
This is useful for testing your webhook server locally without having to open a port or
change any nat configuration on your network.

Example:
	svix listen http://localhost:8000/webhook/

The above command will return you a unique URL and forward any POST requests it receives
to http://localhost:8000/webhook/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			printer := pretty.NewPrinter(getPrinterOptions(cmd))

			urlStr := args[0]
			url, err := url.Parse(urlStr)
			if err != nil {
				return fmt.Errorf("invalid local url %s", urlStr)
			}
			var token string
			if viper.IsSet("relay_token") {
				token = viper.GetString("relay_token")
			} else {
				token, err = relay.GenerateToken()
				printer.CheckErr(err)
				viper.Set("relay_token", token)
				err := config.Write(viper.AllSettings())
				printer.CheckErr(err)
			}
			noLogging, err := cmd.Flags().GetBool(noLoggingFlagName)
			printer.CheckErr(err)

			client := relay.NewClient(token, url, &relay.ClientOptions{
				DisableSecurity: viper.GetBool("relay_disable_security"),
				RelayDebugUrl:   viper.GetString("relay_debug_url"),
				Logging:         !noLogging,
			})
			client.Listen(context.Background())
			return nil
		},
	}
	lc.cmd.Flags().Bool(noLoggingFlagName, false, "Disables History Logging")
	return lc
}
