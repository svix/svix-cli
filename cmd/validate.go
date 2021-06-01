package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/svixhq/svix-cli/utils"
	"github.com/svixhq/svix-cli/validators"
	svix "github.com/svixhq/svix-libs/go"
)

type verifyCmd struct {
	cmd *cobra.Command
}

func newVerifyCmd() *verifyCmd {
	ac := &verifyCmd{}
	ac.cmd = &cobra.Command{
		Use:   "verify SECRET MSG_ID TIMESTAMP SIGNATURE [JSON_PAYLOAD]",
		Short: "Verify the signature of a webhook message",
		Args:  validators.RangeArgs(4, 5),
		Run: func(cmd *cobra.Command, args []string) {
			// parse args
			secret := args[0]
			msgID := args[1]
			timestamp := args[2]
			signature := args[3]

			var payload []byte
			if len(args) > 4 {
				payload = []byte(args[4])
			} else {
				var err error
				payload, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}

			headers := http.Header{}
			headers.Set("svix-id", msgID)
			headers.Set("svix-timestamp", timestamp)
			headers.Set("svix-signature", signature)

			wh, err := svix.NewWebhook(secret)
			cobra.CheckErr(err)
			err = wh.Verify(payload, headers)
			cobra.CheckErr(err)
			fmt.Println("Message Signature Is Valid!")
		},
	}

	return ac
}
