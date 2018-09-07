package command

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/generator/yubikey"
	qr "github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

var registerFlags = struct {
	issuer string
	qr     bool
}{}

var registerCmd = &cobra.Command{

	Use:   "register :username",
	Short: "Initializes an device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		username := args[0]

		creds, err := newSession(rootFlags)

		if err != nil {
			return err
		}

		arn, err := yubikey.NewMFA(creds.UpdateSession(session.Must(session.NewSession())),
			func(secret string) {

				if registerFlags.qr {

					uri := fmt.Sprintf("otpauth://totp/%s@%s?secret=%s&issuer=%s",
						username,
						creds.Profile,
						secret,
						registerFlags.issuer,
					)

					qr.Generate(uri, qr.L, os.Stderr)

				}

			},
			username,
		)

		if err != nil {
			return err
		}

		fmt.Println(*arn)

		return nil

	},
}

func init() {

	flag(registerCmd.Flags(),
		&registerFlags.issuer,
		"Amazon",
		"issuer",
		"i",
		"AWSU_QR_ISSUER",
		"issuer parameter in the QR key uri",
	)

	flag(registerCmd.Flags(),
		&registerFlags.qr,
		true,
		"qr",
		"q",
		"AWSU_QR",
		"generate a QR barcode as backup for soft tokens",
	)

	rootCmd.AddCommand(registerCmd)

}
