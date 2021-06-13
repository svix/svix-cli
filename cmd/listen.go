package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/svix/svix-cli/relay"
)

type listenCmd struct {
	cmd *cobra.Command
}

func newListenCmd() *listenCmd {
	lc := &listenCmd{}
	lc.cmd = &cobra.Command{
		Use:   "listen localURL",
		Short: "Forward webhook requests a local url",
		Long: `listen creates an on-the-fly publicly accessible URL for use when testing webhooks.

The cli then acts as a proxy forwarding any requests to the given localURL.
This is useful for testing your webhook server locally without having to open a port or
change any nat configuration on your network.

Example:
	svix listen http://localhost:8000/webhook/

The above command will return you a unique URL and forward any POST requests it receives
to http://localhost:8000/webhook/`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			urlStr := args[0]
			url, err := url.Parse(urlStr)
			if err != nil {
				return fmt.Errorf("invalid local url %s", urlStr)
			}
			client := relay.NewClient(url, &relay.ClientOptions{
				DisableSecurity: viper.GetBool("relay_disable_security"),
				RelayDebugUrl:   viper.GetString("relay_debug_url"),
			})
			client.Listen(context.Background())
			return nil
		},
	}
	return lc
}
