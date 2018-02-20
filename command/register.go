package command

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/yubikey"
	qr "github.com/mdp/qrterminal"
	"github.com/spf13/cobra"
)

var registerFlags = struct {
	issuer string
}{}

var registerCmd = &cobra.Command{

	Use:   "register",
	Short: "Initializes an device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		username := args[0]

		sess, err := newSession(rootFlags.noCache,
			rootFlags.profile,
			rootFlags.profiles)

		if err != nil {
			return err
		}

		arn, err := yubikey.NewMFA(
			session.Must(
				session.NewSession(
					&aws.Config{
						Credentials: sess.Credentials(),
					},
				),
			),
			func(secret string) {

				uri := fmt.Sprintf("otpauth://totp/%s@%s?secret=%s&issuer=%s",
					username,
					sess.Profile,
					secret,
					registerFlags.issuer,
				)

				qr.Generate(uri, qr.L, os.Stderr)

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

	registerCmd.Flags().StringVarP(&registerFlags.issuer, "issuer", "i", "Amazon", "issuer parameter in the QR key uri")

	rootCmd.AddCommand(registerCmd)

}
