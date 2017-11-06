package command

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/kreuzwerker/awsu/yubikey"
	"github.com/spf13/cobra"
)

var registerFlags = struct {
	filename string
}{}

var registerCmd = &cobra.Command{

	Use:   "register",
	Short: "Initializes an device on AWS and Yubikey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		sess, err := newSession(rootFlags.workspace)

		if err != nil {
			return err
		}

		arn, err := yubikey.NewMFA(
			session.Must(
				session.NewSession(
					&aws.Config{
						Credentials: credentials.NewStaticCredentialsFromCreds(sess.Value),
					},
				),
			),
			registerFlags.filename,
			args[0],
		)

		if err != nil {
			return err
		}

		fmt.Println(*arn)

		return nil

	},
}

func init() {

	registerCmd.Flags().StringVarP(&registerFlags.filename, "filename", "f", "qr.png", "filename for the QR code PNG")

	rootCmd.AddCommand(registerCmd)

}
