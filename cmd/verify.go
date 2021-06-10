package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/svix/svix-cli/utils"
	"github.com/svix/svix-cli/validators"
	svix "github.com/svix/svix-libs/go"
)

type verifyCmd struct {
	cmd *cobra.Command
}

func newVerifyCmd() *verifyCmd {
	secretFlagName := "secret"
	signatureFlagName := "signature"
	msgIdFlagName := "msg-id"
	timestampFlagName := "timestamp"
	ac := &verifyCmd{}
	ac.cmd = &cobra.Command{
		Use:   "verify [JSON_PAYLOAD]",
		Short: "Verify the signature of a webhook message",
		Args:  validators.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// parse args
			var payload []byte
			if len(args) > 0 {
				payload = []byte(args[0])
			} else {
				var err error
				payload, err = utils.ReadPipe()
				cobra.CheckErr(err)
			}

			// ensure all flags are set
			if !cmd.Flags().Changed(secretFlagName) {
				return fmt.Errorf("Secret required for verification!")
			} else if !cmd.Flags().Changed(signatureFlagName) {
				return fmt.Errorf("Signature required for verification!")
			} else if !cmd.Flags().Changed(timestampFlagName) {
				return fmt.Errorf("Timestamp required for verification!")
			} else if !cmd.Flags().Changed(msgIdFlagName) {
				return fmt.Errorf("Message ID required for verifcation")
			}

			// get flags
			secret, err := cmd.Flags().GetString(secretFlagName)
			cobra.CheckErr(err)
			msgID, err := cmd.Flags().GetString(msgIdFlagName)
			cobra.CheckErr(err)
			timestamp, err := cmd.Flags().GetString(timestampFlagName)
			cobra.CheckErr(err)
			signature, err := cmd.Flags().GetString(signatureFlagName)
			cobra.CheckErr(err)

			headers := http.Header{}
			headers.Set("svix-id", msgID)
			headers.Set("svix-timestamp", timestamp)
			headers.Set("svix-signature", signature)

			wh, err := svix.NewWebhook(secret)
			cobra.CheckErr(err)
			err = wh.Verify(payload, headers)
			cobra.CheckErr(err)
			fmt.Println("Message Signature Is Valid!")
			return nil
		},
	}
	ac.cmd.Flags().String(secretFlagName, "", "")
	ac.cmd.Flags().String(msgIdFlagName, "", "")
	ac.cmd.Flags().String(timestampFlagName, "", "")
	ac.cmd.Flags().String(signatureFlagName, "", "")
	return ac
}
